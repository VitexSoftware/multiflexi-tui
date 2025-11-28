package ui

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ApplicationEditorModel represents the editor for an application
type ApplicationEditorModel struct {
	app       cli.Application
	nameInput textinput.Model
	cursor    int
	width     int
	height    int
}

// SaveApplicationMsg is sent when the application should be saved
type SaveApplicationMsg struct {
	App cli.Application
}

// NewApplicationEditorModel creates a new application editor model
func NewApplicationEditorModel(app cli.Application) ApplicationEditorModel {
	nameInput := textinput.New()
	nameInput.SetValue(app.Name)
	nameInput.Focus()

	return ApplicationEditorModel{
		app:       app,
		nameInput: nameInput,
	}
}

// Init initializes the application editor model
func (m ApplicationEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the application editor model
func (m ApplicationEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return BackMsg{} }
		case "enter":
			m.app.Name = m.nameInput.Value()
			return m, func() tea.Msg {
				return SaveApplicationMsg{App: m.app}
			}
		}
	}

	m.nameInput, cmd = m.nameInput.Update(msg)
	return m, cmd
}

// View renders the application editor model
func (m ApplicationEditorModel) View() string {
	return fmt.Sprintf(
		"Editing Application: %d\n\n%s\n\n%s",
		m.app.ID,
		m.nameInput.View(),
		"(press 'enter' to save, 'esc' to cancel)",
	)
}
