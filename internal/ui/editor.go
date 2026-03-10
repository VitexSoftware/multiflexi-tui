package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// EditorField defines a single editable field.
type EditorField struct {
	Label       string
	Placeholder string
	Value       string
}

// SaveEditorMsg is emitted when the user presses Enter to save.
// Data contains the original entity; Fields contains updated values keyed by label.
type SaveEditorMsg struct {
	Data   interface{}
	Fields map[string]string
}

// EditorModel is a unified multi-field form editor.
type EditorModel struct {
	title  string
	data   interface{} // original entity being edited
	inputs []textinput.Model
	labels []string
	cursor int
	width  int
	height int
}

// NewEditorModel creates an editor from a title, the original data, and field definitions.
func NewEditorModel(title string, data interface{}, fields []EditorField) EditorModel {
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
	return EditorModel{
		title:  title,
		data:   data,
		inputs: inputs,
		labels: labels,
		cursor: 0,
	}
}

// Init starts the text input blink cursor.
func (m EditorModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles key events for the editor.
func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
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
			fields := make(map[string]string, len(m.inputs))
			for i, input := range m.inputs {
				fields[m.labels[i]] = input.Value()
			}
			data := m.data
			return m, func() tea.Msg {
				return SaveEditorMsg{Data: data, Fields: fields}
			}
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	m.inputs[m.cursor], cmd = m.inputs[m.cursor].Update(msg)
	return m, cmd
}

// View renders the editor form.
func (m EditorModel) View() string {
	var b strings.Builder

	b.WriteString(GetTitleStyle().Render(m.title))
	b.WriteString("\n\n")

	for i, input := range m.inputs {
		label := m.labels[i]
		if i == m.cursor {
			b.WriteString(GetSelectedItemStyle().Render(fmt.Sprintf("%-15s", label+":")))
		} else {
			b.WriteString(fmt.Sprintf("%-15s", label+":"))
		}
		b.WriteString(" ")
		b.WriteString(input.View())
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(GetFooterStyle().Render("tab/↑↓: navigate fields • enter: save • esc: cancel"))
	b.WriteString("\n")

	return b.String()
}
