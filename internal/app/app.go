package app

import (
	"fmt"
	"strconv"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatusLoadedMsg is sent when status is loaded
type StatusLoadedMsg struct {
	status *cli.StatusInfo
}

// helpLoadedMsg is sent when help content is loaded
type helpLoadedMsg struct {
	command string
	content string
}

// helpErrorMsg is sent when help loading fails
type helpErrorMsg struct {
	command string
	err     error
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return m.loadStatusCmd()
}

// Update handles all messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Clear consumed status messages
	if m.statusMsgConsumed {
		m.statusMessage = ""
		m.statusMsgConsumed = false
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Forward to active view if it exists
		if m.activeView != nil {
			var cmd tea.Cmd
			m.activeView, cmd = m.activeView.Update(msg)
			return m, cmd
		}
		return m, nil

	case StatusLoadedMsg:
		m.statusInfo = msg.status
		return m, nil

	case helpLoadedMsg:
		m.viewer.SetContent(msg.command, msg.content)
		return m, nil

	case helpErrorMsg:
		m.viewer.SetError(msg.err)
		return m, nil

	case ui.ShowHelpMsg:
		m.viewKind = ViewHelp
		m.activeView = nil // viewer rendered inline
		return m, m.loadHelpCmd(msg.Command)

	case ui.ShowMenuMsg:
		m.viewKind = ViewCommands
		m.activeView = nil
		return m, nil

	case ui.BackMsg:
		return m.handleBack()

	// --- Detail/editor navigation messages ---

	case ui.OpenDetailMsg:
		return m.openDetail(msg.Data)

	case ui.EditItemMsg:
		return m.openEditor(msg.Data)

	case ui.SaveEditorMsg:
		return m.handleSave(msg)

	case ui.DeleteItemMsg:
		m.viewKind = ViewConfirmDelete
		m.confirmDialog = ui.NewConfirmDialogModel(msg.Label, msg.Data)
		return m, nil

	case ui.ConfirmDeleteYesMsg:
		return m.executeDelete(msg.Data)

	case ui.ConfirmDeleteNoMsg:
		// Return to detail view
		m.viewKind = ViewDetail
		m.activeView = nil
		return m, nil

	case ui.ScheduleItemMsg:
		m.statusMessage = fmt.Sprintf("Schedule not yet implemented for %T", msg.Data)
		return m, nil

	case ui.StatusMessage:
		m.statusMessage = msg.Text
		return m, nil

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case tea.KeyMsg:
		// Editor and confirm dialog handle all their own keys
		if m.viewKind == ViewEditor || m.viewKind == ViewConfirmDelete {
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			if m.activeView != nil {
				var cmd tea.Cmd
				m.activeView, cmd = m.activeView.Update(msg)
				return m, cmd
			}
			if m.viewKind == ViewConfirmDelete {
				var cmd tea.Cmd
				dialogModel, cmd := m.confirmDialog.Update(msg)
				m.confirmDialog = dialogModel.(ui.ConfirmDialogModel)
				return m, cmd
			}
		}

		// Global keys
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.viewKind == ViewDetail {
				m.viewKind = m.prevKind
				m.activeView = m.prevView
				m.focus = false
				return m, nil
			}
			if !m.focus {
				m.focus = true
				return m, nil
			}
		case "tab":
			if m.focus {
				m.focus = false
			}
			return m, nil
		case "f10":
			m.viewKind = ViewCommands
			m.activeView = nil
			m.focus = true
			return m, nil
		}

		// Menu navigation when focused
		if m.focus && m.viewKind != ViewDetail {
			switch msg.String() {
			case "left", "h":
				if m.menuCursor > 0 {
					m.menuCursor--
					m.updateSelectedHint()
				}
				return m, nil
			case "right", "l":
				if m.menuCursor < len(m.menuItems)-1 {
					m.menuCursor++
					m.updateSelectedHint()
				}
				return m, nil
			case "enter", "space":
				return m.handleMenuSelection()
			}
		}

		// Forward to content view when not focused on menu
		if !m.focus && m.activeView != nil {
			var cmd tea.Cmd
			m.activeView, cmd = m.activeView.Update(msg)
			return m, cmd
		}

		// Special content views without activeView
		if !m.focus {
			switch m.viewKind {
			case ViewDetail:
				var cmd tea.Cmd
				detailModel, cmd := m.detailView.Update(msg)
				m.detailView = detailModel.(ui.DetailViewModel)
				return m, cmd
			case ViewHelp:
				var cmd tea.Cmd
				viewerModel, cmd := m.viewer.Update(msg)
				m.viewer = viewerModel.(ui.ViewerModel)
				return m, cmd
			case ViewCommands:
				var cmd tea.Cmd
				menuModel, cmd := m.menu.Update(msg)
				m.menu = menuModel.(ui.MenuModel)
				return m, cmd
			}
		}
	}

	// Forward non-key messages to active view (for async data loading)
	if m.activeView != nil {
		var cmd tea.Cmd
		m.activeView, cmd = m.activeView.Update(msg)
		return m, cmd
	}

	// Forward to special views
	switch m.viewKind {
	case ViewDetail:
		var cmd tea.Cmd
		detailModel, cmd := m.detailView.Update(msg)
		m.detailView = detailModel.(ui.DetailViewModel)
		return m, cmd
	case ViewHelp:
		var cmd tea.Cmd
		viewerModel, cmd := m.viewer.Update(msg)
		m.viewer = viewerModel.(ui.ViewerModel)
		return m, cmd
	case ViewCommands:
		var cmd tea.Cmd
		menuModel, cmd := m.menu.Update(msg)
		m.menu = menuModel.(ui.MenuModel)
		return m, cmd
	}

	return m, nil
}

// View renders the full UI
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var content string
	switch m.viewKind {
	case ViewHome:
		content = m.renderSystemStatus()
	case ViewDetail:
		content = m.detailView.View()
	case ViewHelp:
		content = m.viewer.View()
	case ViewCommands:
		content = m.menu.View()
	case ViewConfirmDelete:
		content = m.confirmDialog.View()
	default:
		if m.activeView != nil {
			content = m.activeView.View()
		} else {
			content = "No view"
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderMenuBar(),
		content,
		m.renderHelpFooter(),
	)
}

// --- Back navigation ---

func (m Model) handleBack() (tea.Model, tea.Cmd) {
	switch m.viewKind {
	case ViewDetail, ViewEditor, ViewConfirmDelete:
		m.viewKind = m.prevKind
		m.activeView = m.prevView
		m.focus = false
		return m, nil
	default:
		if !m.focus {
			m.focus = true
		}
		return m, nil
	}
}

// --- Detail opening ---

func (m Model) openDetail(data interface{}) (tea.Model, tea.Cmd) {
	m.prevView = m.activeView
	m.prevKind = m.viewKind
	m.viewKind = ViewDetail

	switch d := data.(type) {
	case cli.RunTemplate:
		m.detailView.SetRunTemplate(d)
	case cli.Application:
		m.detailView.SetApplication(d)
	case cli.Job:
		m.detailView.SetJob(d)
	case cli.Company:
		m.detailView.SetCompany(d)
	}
	m.activeView = nil
	return m, nil
}

// --- Editor opening ---

func (m Model) openEditor(data interface{}) (tea.Model, tea.Cmd) {
	m.prevView = m.activeView
	m.prevKind = m.viewKind
	m.viewKind = ViewEditor

	var editor ui.EditorModel
	switch d := data.(type) {
	case cli.Job:
		editor = ui.NewEditorModel(fmt.Sprintf("Editing Job: %d", d.ID), d, []ui.EditorField{
			{Label: "Command", Placeholder: "Command", Value: d.Command},
			{Label: "Executor", Placeholder: "Executor", Value: d.Executor},
			{Label: "Schedule Type", Placeholder: "Schedule Type", Value: d.ScheduleType},
		})
	case cli.Company:
		editor = ui.NewEditorModel(fmt.Sprintf("Editing Company: %s", d.Name), d, []ui.EditorField{
			{Label: "Name", Placeholder: "Company name", Value: d.Name},
			{Label: "Email", Placeholder: "Email", Value: d.Email},
			{Label: "IC", Placeholder: "IC", Value: d.IC},
			{Label: "Slug", Placeholder: "Slug", Value: d.Slug},
		})
	case cli.Application:
		editor = ui.NewEditorModel(fmt.Sprintf("Editing Application: %s", d.Name), d, []ui.EditorField{
			{Label: "Name", Placeholder: "Application name", Value: d.Name},
		})
	case cli.RunTemplate:
		editor = ui.NewEditorModel(fmt.Sprintf("Editing RunTemplate: %s", d.Name), d, []ui.EditorField{
			{Label: "Name", Placeholder: "Template name", Value: d.Name},
		})
	default:
		m.statusMessage = fmt.Sprintf("Edit not supported for %T", data)
		m.viewKind = m.prevKind
		m.activeView = m.prevView
		return m, nil
	}

	m.activeView = editor
	return m, editor.Init()
}

// --- Save handling ---

func (m Model) handleSave(msg ui.SaveEditorMsg) (tea.Model, tea.Cmd) {
	var err error
	var label string

	switch d := msg.Data.(type) {
	case cli.Job:
		d.Command = msg.Fields["Command"]
		d.Executor = msg.Fields["Executor"]
		d.ScheduleType = msg.Fields["Schedule Type"]
		err = cli.UpdateJob(d)
		label = fmt.Sprintf("Job %d", d.ID)
	case cli.Company:
		d.Name = msg.Fields["Name"]
		d.Email = msg.Fields["Email"]
		d.IC = msg.Fields["IC"]
		d.Slug = msg.Fields["Slug"]
		err = cli.UpdateCompany(d)
		label = fmt.Sprintf("Company %s", d.Name)
	case cli.Application:
		d.Name = msg.Fields["Name"]
		err = cli.UpdateApplication(d)
		label = fmt.Sprintf("Application %s", d.Name)
	case cli.RunTemplate:
		d.Name = msg.Fields["Name"]
		err = cli.UpdateRunTemplate(d)
		label = fmt.Sprintf("RunTemplate %s", d.Name)
	}

	if err != nil {
		m.statusMessage = fmt.Sprintf("Error saving %s: %v", label, err)
	} else {
		m.statusMessage = fmt.Sprintf("Saved %s", label)
	}

	m.viewKind = m.prevKind
	m.activeView = m.prevView
	return m, nil
}

// --- Delete handling ---

func (m Model) executeDelete(data interface{}) (tea.Model, tea.Cmd) {
	var err error
	var label string

	switch d := data.(type) {
	case cli.Job:
		err = cli.DeleteJob(d.ID)
		label = fmt.Sprintf("Job %d", d.ID)
	case cli.Application:
		err = cli.DeleteApplication(d.ID)
		label = fmt.Sprintf("Application %s", d.Name)
	case cli.Company:
		err = cli.DeleteCompany(d.ID)
		label = fmt.Sprintf("Company %s", d.Name)
	case cli.RunTemplate:
		err = cli.DeleteRunTemplate(d.ID)
		label = fmt.Sprintf("RunTemplate %s", d.Name)
	default:
		m.statusMessage = "Delete not supported for this item type"
		m.viewKind = m.prevKind
		m.activeView = m.prevView
		return m, nil
	}

	if err != nil {
		m.statusMessage = fmt.Sprintf("Error deleting %s: %v", label, err)
	} else {
		m.statusMessage = fmt.Sprintf("Deleted %s", label)
	}

	m.viewKind = m.prevKind
	m.activeView = m.prevView
	m.focus = false
	return m, nil
}

// --- Async loaders ---

func (m Model) loadStatusCmd() tea.Cmd {
	return func() tea.Msg {
		status, err := cli.GetStatusInfo()
		if err != nil {
			errMsg := err.Error()
			if len(errMsg) > 20 {
				errMsg = errMsg[:20]
			}
			return StatusLoadedMsg{status: &cli.StatusInfo{
				VersionCli: "Error",
				User:       errMsg,
			}}
		}
		return StatusLoadedMsg{status: status}
	}
}

func (m Model) loadHelpCmd(command string) tea.Cmd {
	return func() tea.Msg {
		content, err := cli.GetCommandHelp(command)
		if err != nil {
			return helpErrorMsg{command: command, err: err}
		}
		return helpLoadedMsg{command: command, content: content}
	}
}

// Run starts the Bubbletea program
func Run() error {
	model := NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to start program: %w", err)
	}
	return nil
}

// --- Encryption/Prune/Queue simple models (inline, no separate file needed) ---

type encryptionView struct {
	status       string
	err          error
	initializing bool
}

func newEncryptionView(status string) encryptionView {
	return encryptionView{status: status}
}

func (v encryptionView) Init() tea.Cmd { return nil }

func (v encryptionView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case encryptionInitMsg:
		v.initializing = false
		v.status = "active"
		return v, nil
	case encryptionErrMsg:
		v.initializing = false
		v.err = msg.err
		return v, nil
	case tea.KeyMsg:
		if msg.String() == "enter" && !v.initializing && v.status != "active" {
			v.initializing = true
			return v, func() tea.Msg {
				if err := cli.InitEncryption(); err != nil {
					return encryptionErrMsg{err: err}
				}
				return encryptionInitMsg{}
			}
		}
	}
	return v, nil
}

func (v encryptionView) View() string {
	s := ui.GetTitleStyle().Render("Encryption Status") + "\n\n"
	if v.err != nil {
		s += ui.GetErrorStyle().Render(fmt.Sprintf("Error: %v", v.err)) + "\n\n"
	}
	s += fmt.Sprintf("Current Status: %s\n\n", v.status)
	if v.initializing {
		s += "Initializing encryption..."
	} else if v.status != "active" {
		s += ui.GetButtonStyle().Render("Initialize Encryption")
	}
	return s
}

type encryptionInitMsg struct{}
type encryptionErrMsg struct{ err error }

// --- Prune view ---

type pruneView struct {
	logs    bool
	jobs    bool
	keepStr string
	cursor  int
	pruning bool
	err     error
}

func newPruneView() pruneView {
	return pruneView{keepStr: "1000"}
}

func (v pruneView) Init() tea.Cmd { return nil }

func (v pruneView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pruneDoneMsg:
		v.pruning = false
		return v, nil
	case pruneErrMsg:
		v.pruning = false
		v.err = msg.err
		return v, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.cursor > 0 {
				v.cursor--
			}
		case "down", "j":
			if v.cursor < 2 {
				v.cursor++
			}
		case "enter":
			switch v.cursor {
			case 0:
				v.logs = !v.logs
			case 1:
				v.jobs = !v.jobs
			case 2:
				if !v.pruning {
					v.pruning = true
					keep, _ := strconv.Atoi(v.keepStr)
					logs, jobs := v.logs, v.jobs
					return v, func() tea.Msg {
						if err := cli.Prune(logs, jobs, keep); err != nil {
							return pruneErrMsg{err: err}
						}
						return pruneDoneMsg{}
					}
				}
			}
		}
	}
	return v, nil
}

func (v pruneView) View() string {
	s := ui.GetTitleStyle().Render("Prune Logs and Jobs") + "\n\n"
	if v.err != nil {
		s += ui.GetErrorStyle().Render(fmt.Sprintf("Error: %v", v.err)) + "\n\n"
	}
	check := func(b bool) string {
		if b {
			return "[x]"
		}
		return "[ ]"
	}
	options := []string{
		fmt.Sprintf("%s Prune logs", check(v.logs)),
		fmt.Sprintf("%s Prune jobs", check(v.jobs)),
		fmt.Sprintf("Keep: %s", v.keepStr),
	}
	for i, opt := range options {
		if i == v.cursor {
			s += ui.GetSelectedItemStyle().Render(opt) + "\n"
		} else {
			s += opt + "\n"
		}
	}
	s += "\n"
	if v.pruning {
		s += "Pruning..."
	} else {
		s += ui.GetButtonStyle().Render("Prune")
	}
	return s
}

type pruneDoneMsg struct{}
type pruneErrMsg struct{ err error }
