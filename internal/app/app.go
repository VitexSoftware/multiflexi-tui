package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

// ViewState represents the current view being displayed
type ViewState int

const (
	MenuView ViewState = iota
	HelpView
)

// Model represents the main application model
type Model struct {
	state    ViewState
	menu     ui.MenuModel
	viewer   ui.ViewerModel
	width    int
	height   int
}

// NewModel creates a new application model
func NewModel() (*Model, error) {
	// Load commands from multiflexi-cli
	commands, err := cli.GetCommands()
	if err != nil {
		return nil, fmt.Errorf("failed to load commands: %w", err)
	}

	menu := ui.NewMenuModel(commands)
	viewer := ui.NewViewerModel("", "")

	return &Model{
		state:  MenuView,
		menu:   menu,
		viewer: viewer,
	}, nil
}

// Init initializes the application model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages for the application
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Forward size message to current view
		switch m.state {
		case MenuView:
			var cmd tea.Cmd
			menuModel, cmd := m.menu.Update(msg)
			m.menu = menuModel.(ui.MenuModel)
			return m, cmd
		case HelpView:
			var cmd tea.Cmd
			viewerModel, cmd := m.viewer.Update(msg)
			m.viewer = viewerModel.(ui.ViewerModel)
			return m, cmd
		}

	case ui.ShowHelpMsg:
		// Switch to help view and load help content
		m.state = HelpView
		return m, m.loadHelpCmd(msg.Command)

	case ui.BackToMenuMsg:
		// Switch back to menu view
		m.state = MenuView
		return m, nil

	case helpLoadedMsg:
		// Update viewer with help content
		m.viewer.SetContent(msg.command, msg.content)
		return m, nil

	case helpErrorMsg:
		// Display error in viewer
		m.viewer.SetError(msg.err)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	// Forward message to current view
	switch m.state {
	case MenuView:
		var cmd tea.Cmd
		menuModel, cmd := m.menu.Update(msg)
		m.menu = menuModel.(ui.MenuModel)
		return m, cmd
	case HelpView:
		var cmd tea.Cmd
		viewerModel, cmd := m.viewer.Update(msg)
		m.viewer = viewerModel.(ui.ViewerModel)
		return m, cmd
	}

	return m, nil
}

// View renders the current view
func (m Model) View() string {
	switch m.state {
	case MenuView:
		return m.menu.View()
	case HelpView:
		return m.viewer.View()
	default:
		return "Unknown view state"
	}
}

// helpLoadedMsg is sent when help content is loaded successfully
type helpLoadedMsg struct {
	command string
	content string
}

// helpErrorMsg is sent when there's an error loading help content
type helpErrorMsg struct {
	command string
	err     error
}

// loadHelpCmd returns a command that loads help for a command
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
	model, err := NewModel()
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	
	if err := p.Start(); err != nil {
		return fmt.Errorf("failed to start program: %w", err)
	}

	return nil
}