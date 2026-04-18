package entity

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// EditorView is a generic create/update form driven by an EntityDef.
type EditorView struct {
	client   cli.Client
	def      *EntityDef
	data     interface{} // original entity (nil for create)
	isCreate bool
	title    string

	inputs []textinput.Model
	labels []string
	cursor int
}

// NewEditorView creates an editor for updating or creating an entity.
func NewEditorView(c cli.Client, def *EntityDef, data interface{}, isCreate bool) *EditorView {
	var fields []ui.EditorField
	var title string

	if isCreate {
		title = fmt.Sprintf("New %s", def.Name)
		if def.NewFields != nil {
			fields = def.NewFields()
		}
	} else {
		label := def.Name
		if def.GetLabel != nil {
			label = def.GetLabel(data)
		}
		title = fmt.Sprintf("Edit %s", label)
		if def.ToEditor != nil {
			fields = def.ToEditor(data)
		}
	}

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

	return &EditorView{
		client:   c,
		def:      def,
		data:     data,
		isCreate: isCreate,
		title:    title,
		inputs:   inputs,
		labels:   labels,
	}
}

func (m *EditorView) Init() tea.Cmd {
	return textinput.Blink
}

func (m *EditorView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return ui.NavigateBackMsg{} }
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
			return m.save()
		}
	}

	var cmd tea.Cmd
	if len(m.inputs) > 0 {
		m.inputs[m.cursor], cmd = m.inputs[m.cursor].Update(msg)
	}
	return m, cmd
}

func (m *EditorView) save() (tea.Model, tea.Cmd) {
	fields := make(map[string]string, len(m.inputs))
	for i, input := range m.inputs {
		fields[m.labels[i]] = input.Value()
	}

	client := m.client
	def := m.def
	data := m.data
	isCreate := m.isCreate

	return m, func() tea.Msg {
		var err error
		var label string

		if isCreate {
			if def.CreateArgs != nil {
				args := def.CreateArgs(fields)
				_, err = client.Create(def.CLIEntity, args...)
				label = fmt.Sprintf("New %s", def.Name)
			}
		} else {
			if def.UpdateArgs != nil {
				args := def.UpdateArgs(data, fields)
				err = client.Update(def.CLIEntity, args...)
				if def.GetLabel != nil {
					label = def.GetLabel(data)
				} else {
					label = def.Name
				}
			}
		}

		if err != nil {
			return ui.StatusMsg{Text: fmt.Sprintf("Error saving %s: %v", label, err)}
		}

		return ui.NavigateBackAndRefreshMsg{}
	}
}

func (m *EditorView) View() string {
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
	b.WriteString(ui.FooterStyle().Render("tab/↑↓: fields • enter: save • esc: cancel"))
	b.WriteString("\n")

	return b.String()
}
