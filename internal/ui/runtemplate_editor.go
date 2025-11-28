package ui

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// RunTemplateEditorModel represents the editor for a run template
type RunTemplateEditorModel struct {
	template  cli.RunTemplate
	nameInput textinput.Model
	cursor    int
	width     int
	height    int
}

// SaveRunTemplateMsg is sent when the run template should be saved
type SaveRunTemplateMsg struct {
	Template cli.RunTemplate
}

// NewRunTemplateEditorModel creates a new run template editor model
func NewRunTemplateEditorModel(template cli.RunTemplate) RunTemplateEditorModel {
	nameInput := textinput.New()
	nameInput.SetValue(template.Name)
	nameInput.Focus()

	return RunTemplateEditorModel{
		template:  template,
		nameInput: nameInput,
	}
}

// Init initializes the run template editor model
func (m RunTemplateEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the run template editor model
func (m RunTemplateEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.template.Name = m.nameInput.Value()
			return m, func() tea.Msg {
				return SaveRunTemplateMsg{Template: m.template}
			}
		}
	}

	m.nameInput, cmd = m.nameInput.Update(msg)
	return m, cmd
}

// View renders the run template editor model
func (m RunTemplateEditorModel) View() string {
	return fmt.Sprintf(
		"Editing Run Template: %d\n\n%s\n\n%s",
		m.template.ID,
		m.nameInput.View(),
		"(press 'enter' to save, 'esc' to cancel)",
	)
}