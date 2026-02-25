package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// JobEditorModel represents the editor for a job
type JobEditorModel struct {
	job              cli.Job
	commandInput     textinput.Model
	executorInput    textinput.Model
	scheduleTypeInput textinput.Model
	cursor           int
	inputs           []textinput.Model
	labels           []string
	width            int
	height           int
}

// SaveJobMsg is sent when the job should be saved
type SaveJobMsg struct {
	Job cli.Job
}

// NewJobEditorModel creates a new job editor model
func NewJobEditorModel(job cli.Job) JobEditorModel {
	commandInput := textinput.New()
	commandInput.Placeholder = "Command"
	commandInput.SetValue(job.Command)
	commandInput.Focus()

	executorInput := textinput.New()
	executorInput.Placeholder = "Executor"
	executorInput.SetValue(job.Executor)

	scheduleTypeInput := textinput.New()
	scheduleTypeInput.Placeholder = "Schedule Type"
	scheduleTypeInput.SetValue(job.ScheduleType)

	inputs := []textinput.Model{commandInput, executorInput, scheduleTypeInput}
	labels := []string{"Command", "Executor", "Schedule Type"}

	return JobEditorModel{
		job:    job,
		inputs: inputs,
		labels: labels,
		cursor: 0,
	}
}

// Init initializes the job editor model
func (m JobEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the job editor model
func (m JobEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return BackMsg{} }
		case "tab", "down":
			m.inputs[m.cursor].Blur()
			m.cursor = (m.cursor + 1) % len(m.inputs)
			m.inputs[m.cursor].Focus()
			return m, textinput.Blink
		case "shift+tab", "up":
			m.inputs[m.cursor].Blur()
			m.cursor = (m.cursor - 1 + len(m.inputs)) % len(m.inputs)
			m.inputs[m.cursor].Focus()
			return m, textinput.Blink
		case "enter":
			m.job.Command = m.inputs[0].Value()
			m.job.Executor = m.inputs[1].Value()
			m.job.ScheduleType = m.inputs[2].Value()
			return m, func() tea.Msg {
				return SaveJobMsg{Job: m.job}
			}
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	m.inputs[m.cursor], cmd = m.inputs[m.cursor].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the job editor model
func (m JobEditorModel) View() string {
	var content strings.Builder

	content.WriteString(GetTitleStyle().Render(fmt.Sprintf("Editing Job: %d", m.job.ID)))
	content.WriteString("\n\n")

	for i, input := range m.inputs {
		label := m.labels[i]
		if i == m.cursor {
			content.WriteString(GetSelectedItemStyle().Render(fmt.Sprintf("%-15s", label+":")))
		} else {
			content.WriteString(fmt.Sprintf("%-15s", label+":"))
		}
		content.WriteString(" ")
		content.WriteString(input.View())
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(GetFooterStyle().Render("tab/↑↓: navigate fields • enter: save • esc: cancel"))
	content.WriteString("\n")

	return content.String()
}
