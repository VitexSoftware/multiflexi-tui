package app

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// App is the top-level bubbletea model.
type App struct {
	Client cli.Client
	nav    Navigator
	items  []MenuItem

	// Current view (nil = home/status view)
	activeView tea.Model

	// Menu bar
	menuCursor     int
	activeMenuItem int
	menuFocus      bool // true = menu focused

	// Layout
	width, height int

	// Status
	statusInfo    *cli.StatusInfo
	statusMessage string
}

// New creates a new App with the given client and menu items.
func New(client cli.Client, items []MenuItem) *App {
	return &App{
		Client:    client,
		items:     items,
		menuFocus: true,
	}
}

// statusLoadedMsg carries the loaded status.
type statusLoadedMsg struct{ status *cli.StatusInfo }

func (a *App) Init() tea.Cmd {
	return func() tea.Msg {
		status, err := a.Client.GetStatus()
		if err != nil {
			return statusLoadedMsg{status: &cli.StatusInfo{VersionCli: "Error", User: err.Error()}}
		}
		return statusLoadedMsg{status: status}
	}
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		if a.activeView != nil {
			var cmd tea.Cmd
			a.activeView, cmd = a.activeView.Update(msg)
			return a, cmd
		}
		return a, nil

	case statusLoadedMsg:
		a.statusInfo = msg.status
		return a, nil

	case ui.NavigateToMsg:
		// Push current view onto stack, switch to new view
		a.nav.Push(ViewState{View: a.activeView, MenuIdx: a.activeMenuItem})
		a.activeView = msg.View
		a.menuFocus = false
		return a, msg.View.Init()

	case ui.NavigateBackMsg:
		return a.goBack()

	case ui.StatusMsg:
		a.statusMessage = msg.Text
		return a, nil

	case ui.ConfirmMsg:
		confirm := ui.NewConfirmDialog(msg.Label, msg.Action)
		a.nav.Push(ViewState{View: a.activeView, MenuIdx: a.activeMenuItem})
		a.activeView = confirm
		a.menuFocus = false
		return a, nil

	case ui.ConfirmYesMsg:
		prev, _ := a.nav.Pop()
		a.activeView = prev.View
		if msg.Action != nil {
			return a, func() tea.Msg { return msg.Action() }
		}
		return a, nil

	case ui.ConfirmNoMsg:
		prev, _ := a.nav.Pop()
		a.activeView = prev.View
		return a, nil

	case tea.MouseMsg:
		return a.handleMouse(msg)

	case tea.KeyMsg:
		return a.handleKey(msg)
	}

	// Forward non-key messages to active view
	if a.activeView != nil {
		var cmd tea.Cmd
		a.activeView, cmd = a.activeView.Update(msg)
		return a, cmd
	}
	return a, nil
}

func (a *App) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global quit
	if key == "ctrl+c" {
		return a, tea.Quit
	}

	// If active view is a confirm dialog, let it handle keys
	if _, ok := a.activeView.(*ui.ConfirmDialog); ok {
		var cmd tea.Cmd
		a.activeView, cmd = a.activeView.Update(msg)
		return a, cmd
	}

	switch key {
	case "q":
		if a.menuFocus && a.activeView == nil {
			return a, tea.Quit
		}
		if a.menuFocus {
			return a, tea.Quit
		}
	case "esc":
		if !a.menuFocus {
			if a.nav.Depth() > 0 {
				return a.goBack()
			}
			a.menuFocus = true
			return a, nil
		}
	case "tab":
		a.menuFocus = !a.menuFocus
		return a, nil
	}

	// Menu navigation when focused
	if a.menuFocus {
		switch key {
		case "left", "h":
			if a.menuCursor > 0 {
				a.menuCursor--
			}
			return a, nil
		case "right", "l":
			if a.menuCursor < len(a.items)-1 {
				a.menuCursor++
			}
			return a, nil
		case "enter", " ":
			return a.selectMenuItem()
		case "q":
			return a, tea.Quit
		}
		return a, nil
	}

	// Forward to active view
	if a.activeView != nil {
		var cmd tea.Cmd
		a.activeView, cmd = a.activeView.Update(msg)
		return a, cmd
	}
	return a, nil
}

func (a *App) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.MouseLeft:
		if msg.Y == 0 {
			// Click on menu bar — calculate which item
			xPos := len("MultiFlexi TUI") + 4 + 1
			for i, item := range a.items {
				itemWidth := len(item.Label) + 3
				if msg.X >= xPos && msg.X < xPos+itemWidth {
					a.menuCursor = i
					return a.selectMenuItem()
				}
				xPos += itemWidth
			}
		} else if msg.Y >= 3 && a.menuFocus {
			a.menuFocus = false
		}
	case tea.MouseWheelUp:
		if !a.menuFocus && a.activeView != nil {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			var cmd tea.Cmd
			a.activeView, cmd = a.activeView.Update(keyMsg)
			return a, cmd
		}
	case tea.MouseWheelDown:
		if !a.menuFocus && a.activeView != nil {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			var cmd tea.Cmd
			a.activeView, cmd = a.activeView.Update(keyMsg)
			return a, cmd
		}
	}
	return a, nil
}

func (a *App) goBack() (tea.Model, tea.Cmd) {
	prev, ok := a.nav.Pop()
	if !ok {
		a.activeView = nil
		a.menuFocus = true
		return a, nil
	}
	a.activeView = prev.View
	if prev.View == nil {
		a.menuFocus = true
	}
	return a, nil
}

func (a *App) selectMenuItem() (tea.Model, tea.Cmd) {
	if a.menuCursor < 0 || a.menuCursor >= len(a.items) {
		return a, nil
	}
	a.activeMenuItem = a.menuCursor
	item := a.items[a.menuCursor]
	if item.Action != nil {
		// Clear nav stack when selecting from menu
		a.nav.Clear()
		view, cmd := item.Action(a)
		a.activeView = view
		a.menuFocus = false
		if view != nil && cmd == nil {
			cmd = view.Init()
		}
		return a, cmd
	}
	return a, nil
}

// View renders the full UI.
func (a *App) View() string {
	if a.width == 0 {
		return "Initializing..."
	}

	var content string
	if a.activeView != nil {
		content = a.activeView.View()
	} else {
		content = a.renderStatus()
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		a.renderMenuBar(),
		content,
		a.renderFooter(),
	)
}

func (a *App) renderMenuBar() string {
	w := a.width
	if w == 0 {
		w = 80
	}
	title := "MultiFlexi TUI"
	var parts []string
	for i, item := range a.items {
		var rendered string
		if i == a.menuCursor && a.menuFocus {
			rendered = ui.SelectedStyle().Render(" " + item.Label + " ")
		} else if i == a.activeMenuItem {
			rendered = ui.ActiveMenuStyle().Render(" " + item.Label + " ")
		} else {
			rendered = ui.UnselectedStyle().Render(" " + item.Label + " ")
		}
		parts = append(parts, rendered)
	}
	menuLine := ui.TitleStyle().Render(" "+title+" ") + " " + strings.Join(parts, " ")

	hint := "Navigation: ←/→ to move, Enter to select"
	if a.menuCursor >= 0 && a.menuCursor < len(a.items) {
		hint = a.items[a.menuCursor].Hint
	}
	hintLine := ui.DescriptionStyle().Render(" " + hint + " ")
	sep := strings.Repeat("═", w)

	return menuLine + "\n" + hintLine + "\n" + sep + "\n"
}

func (a *App) renderFooter() string {
	w := a.width
	if w == 0 {
		w = 80
	}
	sep := strings.Repeat("═", w)
	var helpLine string
	if a.menuFocus {
		helpLine = ui.FooterStyle().Render(" ←/→: navigate menu • enter: select • tab: content • q: quit ")
	} else {
		helpLine = ui.FooterStyle().Render(" ↑/↓: rows • ←/→: pages • enter: detail • e: edit • n: new • esc: back • tab: menu ")
	}
	statusLine := ""
	if a.statusMessage != "" {
		statusLine = ui.FooterStyle().Render(" " + a.statusMessage + " ")
	}
	return sep + "\n" + statusLine + "\n" + helpLine
}

func (a *App) renderStatus() string {
	var b strings.Builder
	b.WriteString(ui.TitleStyle().Render("🖥️  MultiFlexi System Dashboard"))
	b.WriteString("\n\n")

	if a.statusInfo == nil {
		b.WriteString(ui.DescriptionStyle().Render("Loading system status..."))
		b.WriteString("\n")
		return b.String()
	}
	s := a.statusInfo
	rows := []struct{ icon, label, value string }{
		{"🔧", "CLI Version", s.VersionCli},
		{"🗄️", "DB Migration", s.DbMigration},
		{"👤", "User", s.User},
		{"🐘", "PHP", s.PHP},
		{"💻", "OS", s.OS},
		{"🧠", "Memory", fmt.Sprintf("%d KB", s.Memory)},
		{"🏢", "Companies", fmt.Sprintf("%d", s.Companies)},
		{"📱", "Applications", fmt.Sprintf("%d", s.Apps)},
		{"📄", "RunTemplates", fmt.Sprintf("%d", s.RunTemplates)},
		{"🏷️", "Topics", fmt.Sprintf("%d", s.Topics)},
		{"🔑", "Credentials", fmt.Sprintf("%d", s.Credentials)},
		{"🎭", "Credential Types", fmt.Sprintf("%d", s.CredentialTypes)},
		{"💼", "Jobs", s.Jobs},
		{"⚙️", "Executor", s.Executor},
		{"📅", "Scheduler", s.Scheduler},
		{"🔐", "Encryption", s.Encryption},
		{"📊", "Zabbix", s.Zabbix},
		{"📈", "Telemetry", s.Telemetry},
		{"🕒", "Timestamp", s.Timestamp},
	}
	for _, r := range rows {
		style := ui.DescriptionStyle()
		if r.value == "active" {
			style = ui.ActiveStatusStyle()
		} else if r.value == "inactive" || r.value == "disabled" {
			style = ui.DisabledStatusStyle()
		}
		b.WriteString(fmt.Sprintf("%s %-18s %s\n", r.icon, r.label+":", style.Render(r.value)))
	}
	return b.String()
}

// Run starts the TUI application.
func Run(client cli.Client, items []MenuItem) error {
	app := New(client, items)
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
