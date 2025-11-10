package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewerModel represents the help text viewer screen
type ViewerModel struct {
	viewport    viewport.Model
	content     string
	commandName string
	width       int
	height      int
	ready       bool
	err         error
}

// NewViewerModel creates a new viewer model
func NewViewerModel(commandName, content string) ViewerModel {
	return ViewerModel{
		content:     content,
		commandName: commandName,
	}
}

// Init initializes the viewer model
func (m ViewerModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the viewer model
func (m ViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		headerHeight := 4 // Title + borders + padding
		footerHeight := 3 // Footer + margin
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width-4, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "q", "esc":
			return m, func() tea.Msg {
				return BackToMenuMsg{}
			}
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the viewer
func (m ViewerModel) View() string {
	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error: %v", m.err))
	}

	if !m.ready {
		return "Loading..."
	}

	title := GetTitleStyle().Render(fmt.Sprintf("Help: %s", m.commandName))
	
	// Create a styled viewport container
	viewportContent := GetViewerStyle().Render(m.viewport.View())
	
	footer := GetFooterStyle().Render("↑/↓: scroll • q: back to menu • ctrl+c: quit")
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		viewportContent,
		footer,
	)
}

// BackToMenuMsg is a message to return to the menu
type BackToMenuMsg struct{}

// SetContent updates the viewer content
func (m *ViewerModel) SetContent(commandName, content string) {
	m.commandName = commandName
	m.content = content
	if m.ready {
		m.viewport.SetContent(content)
		m.viewport.GotoTop()
	}
}

// SetError sets an error to display
func (m *ViewerModel) SetError(err error) {
	m.err = err
}