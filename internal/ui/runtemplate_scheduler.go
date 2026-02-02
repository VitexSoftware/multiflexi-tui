package ui

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// RunTemplateSchedulerModel represents the scheduler for a run template
type RunTemplateSchedulerModel struct {
	template      cli.RunTemplate
	intervalInput textinput.Model
	cursor        int
	width         int
	height        int
}

// NewRunTemplateSchedulerModel creates a new run template scheduler model
func NewRunTemplateSchedulerModel(template cli.RunTemplate) RunTemplateSchedulerModel {
	intervalInput := textinput.New()
	intervalInput.SetValue(template.Interv)
	intervalInput.Focus()

	return RunTemplateSchedulerModel{
		template:      template,
		intervalInput: intervalInput,
	}
}

// Init initializes the run template scheduler model
func (m RunTemplateSchedulerModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the run template scheduler model
func (m RunTemplateSchedulerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.template.Interv = m.intervalInput.Value()
			return m, func() tea.Msg {
				return SaveRunTemplateMsg{Template: m.template}
			}
		}
	}

	m.intervalInput, cmd = m.intervalInput.Update(msg)
	return m, cmd
}

// View renders the run template scheduler model
func (m RunTemplateSchedulerModel) View() string {
	return fmt.Sprintf(
		"Scheduling Run Template: %s\n\nInterval: %s\n\n%s",
		m.template.Name,
		m.intervalInput.View(),
		"(press 'enter' to save, 'esc' to cancel)",
	)
}
