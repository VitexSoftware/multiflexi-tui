package entity

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ActionFormView is a generic form that calls onSave with field values when submitted.
// Used for actions that require user input before executing (e.g., schedule, save-to-file).
type ActionFormView struct {
	title  string
	labels []string
	inputs []textinput.Model
	cursor int
	onSave func(fields map[string]string) tea.Cmd
}

// NewActionFormView creates an action form with the given title, fields, and save callback.
func NewActionFormView(title string, fields []ui.EditorField, onSave func(map[string]string) tea.Cmd) *ActionFormView {
	inputs := make([]textinput.Model, len(fields))
	labels := make([]string, len(fields))
	for i, f := range fields {
		ti := textinput.New()
		ti.Placeholder = f.Placeholder
		ti.SetValue(f.Value)
		if i == 0 {
			ti.Focus()
		}
		inputs[i] = ti
		labels[i] = f.Label
	}
	return &ActionFormView{
		title:  title,
		labels: labels,
		inputs: inputs,
		onSave: onSave,
	}
}

func (m *ActionFormView) Init() tea.Cmd { return textinput.Blink }

func (m *ActionFormView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return ui.NavigateBackMsg{} }
		case "tab", "down":
			if len(m.inputs) > 0 {
				m.inputs[m.cursor].Blur()
				m.cursor = (m.cursor + 1) % len(m.inputs)
				m.inputs[m.cursor].Focus()
			}
			return m, textinput.Blink
		case "shift+tab", "up":
			if len(m.inputs) > 0 {
				m.inputs[m.cursor].Blur()
				m.cursor = (m.cursor - 1 + len(m.inputs)) % len(m.inputs)
				m.inputs[m.cursor].Focus()
			}
			return m, textinput.Blink
		case "enter":
			fields := make(map[string]string, len(m.inputs))
			for i, inp := range m.inputs {
				fields[m.labels[i]] = inp.Value()
			}
			return m, m.onSave(fields)
		}
	}
	if len(m.inputs) > 0 {
		var cmd tea.Cmd
		m.inputs[m.cursor], cmd = m.inputs[m.cursor].Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *ActionFormView) View() string {
	var b strings.Builder
	b.WriteString(ui.TitleStyle().Render(m.title))
	b.WriteString("\n\n")
	for i, input := range m.inputs {
		label := m.labels[i]
		if i == m.cursor {
			b.WriteString(ui.SelectedStyle().Render(fmt.Sprintf("%-15s", label+":")))
		} else {
			b.WriteString(fmt.Sprintf("%-15s", label+":"))
		}
		b.WriteString(" ")
		b.WriteString(input.View())
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(ui.FooterStyle().Render("tab/↑↓: fields • enter: confirm • esc: cancel"))
	b.WriteString("\n")
	return b.String()
}
