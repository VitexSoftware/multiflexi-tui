package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
)

// MenuModel represents the command list screen
type MenuModel struct {
	list     list.Model
	commands []cli.Command
	width    int
	height   int
	err      error
}

// menuItem represents an item in the command list
type menuItem struct {
	command cli.Command
}

func (i menuItem) FilterValue() string {
	return i.command.Name
}

func (i menuItem) Title() string {
	return GetSelectedItemStyle().Render(i.command.Name)
}

func (i menuItem) Description() string {
	return GetItemDescriptionStyle().Render(i.command.Description)
}

// menuDelegate defines how list items are rendered
type menuDelegate struct{}

func (d menuDelegate) Height() int { return 2 }
func (d menuDelegate) Spacing() int { return 1 }
func (d menuDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d menuDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(menuItem)
	if !ok {
		return
	}

	var style lipgloss.Style
	if index == m.Index() {
		style = GetSelectedItemStyle()
	} else {
		style = GetUnselectedItemStyle()
	}

	title := style.Render(i.command.Name)
	description := GetItemDescriptionStyle().Render(i.command.Description)
	
	fmt.Fprint(w, title+"\n"+description)
}

// NewMenuModel creates a new menu model
func NewMenuModel(commands []cli.Command) MenuModel {
	items := make([]list.Item, len(commands))
	for i, cmd := range commands {
		items[i] = menuItem{command: cmd}
	}

	l := list.New(items, menuDelegate{}, 80, 20)
	l.Title = "MultiFlexi Commands"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = GetTitleStyle()
	l.Styles.PaginationStyle = GetFooterStyle()
	l.Styles.HelpStyle = GetFooterStyle()

	return MenuModel{
		list:     l,
		commands: commands,
	}
}

// Init initializes the menu model
func (m MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the menu model
func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width - 4)
		m.list.SetHeight(msg.Height - 8)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				item := selectedItem.(menuItem)
				return m, func() tea.Msg {
					return ShowHelpMsg{Command: item.command.Name}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the menu
func (m MenuModel) View() string {
	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error: %v", m.err))
	}

	content := m.list.View()
	footer := GetFooterStyle().Render("↑/↓: navigate • enter: select • q: quit")
	
	return content + "\n" + footer
}

// ShowHelpMsg is a message to show help for a command
type ShowHelpMsg struct {
	Command string
}