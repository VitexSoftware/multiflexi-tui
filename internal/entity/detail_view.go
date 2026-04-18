package entity

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// DetailView shows the details of a single entity with action buttons.
type DetailView struct {
	client         cli.Client
	def            *EntityDef
	data           interface{}
	fields         []ui.DetailField
	actions        []ui.ActionDef
	selectedAction int
}

// NewDetailView creates a detail view for the given entity data.
func NewDetailView(c cli.Client, def *EntityDef, data interface{}) *DetailView {
	var fields []ui.DetailField
	if def.ToDetail != nil {
		fields = def.ToDetail(data)
	}
	return &DetailView{
		client:  c,
		def:     def,
		data:    data,
		fields:  fields,
		actions: def.Actions,
	}
}

func (m *DetailView) Init() tea.Cmd { return nil }

func (m *DetailView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "esc", "q":
			return m, func() tea.Msg { return ui.NavigateBackMsg{} }
		case "tab", "right":
			if len(m.actions) > 0 {
				m.selectedAction = (m.selectedAction + 1) % len(m.actions)
			}
		case "left":
			if len(m.actions) > 0 {
				m.selectedAction = (m.selectedAction - 1 + len(m.actions)) % len(m.actions)
			}
		case "enter":
			return m.executeAction()
		default:
			// Check shortcut keys
			for _, action := range m.actions {
				if action.Key == key {
					return m.executeActionByCommand(action.Command)
				}
			}
		}
	}
	return m, nil
}

func (m *DetailView) executeAction() (tea.Model, tea.Cmd) {
	if m.selectedAction < len(m.actions) {
		return m.executeActionByCommand(m.actions[m.selectedAction].Command)
	}
	return m, nil
}

func (m *DetailView) executeActionByCommand(command string) (tea.Model, tea.Cmd) {
	switch command {
	case "edit":
		if m.def.ToEditor != nil {
			editor := NewEditorView(m.client, m.def, m.data, false)
			return m, func() tea.Msg { return ui.NavigateToMsg{View: editor} }
		}
	case "delete":
		if m.def.GetID != nil && m.def.GetLabel != nil {
			id := m.def.GetID(m.data)
			label := m.def.GetLabel(m.data)
			cliEntity := m.def.CLIEntity
			deleteAction := m.def.DeleteAction
			client := m.client
			return m, func() tea.Msg {
				return ui.ConfirmMsg{
					Label: fmt.Sprintf("Delete %s?", label),
					Action: func() tea.Msg {
						err := client.Delete(cliEntity, deleteAction, id)
						if err != nil {
							return ui.StatusMsg{Text: fmt.Sprintf("Error deleting %s: %v", label, err)}
						}
						return ui.StatusMsg{Text: fmt.Sprintf("Deleted %s", label)}
					},
				}
			}
		}
	default:
		for _, action := range m.actions {
			if action.Command != command || action.Handler == nil {
				continue
			}
			if action.Confirm != "" {
				handler := action.Handler
				client := m.client
				data := m.data
				label := action.Confirm
				return m, func() tea.Msg {
					return ui.ConfirmMsg{
						Label: label,
						Action: func() tea.Msg {
							cmd := handler(client, data)
							if cmd != nil {
								return cmd()
							}
							return nil
						},
					}
				}
			}
			return m, action.Handler(m.client, m.data)
		}
		return m, func() tea.Msg {
			return ui.StatusMsg{Text: fmt.Sprintf("Action '%s' not implemented", command)}
		}
	}
	return m, nil
}

func (m *DetailView) View() string {
	var b strings.Builder

	label := m.def.Name
	if m.def.GetLabel != nil {
		label = m.def.GetLabel(m.data)
	}
	b.WriteString(ui.TitleStyle().Render(label))
	b.WriteString("\n\n")

	// Fields
	maxW := 0
	for _, f := range m.fields {
		if len(f.Label) > maxW {
			maxW = len(f.Label)
		}
	}
	for _, f := range m.fields {
		b.WriteString(fmt.Sprintf("%-*s: %s\n", maxW, f.Label, f.Value))
	}
	b.WriteString("\n")

	// Action buttons
	if len(m.actions) > 0 {
		var buttons []string
		for i, a := range m.actions {
			btn := fmt.Sprintf("[%s] %s", a.Key, a.Label)
			if i == m.selectedAction {
				buttons = append(buttons, ui.SelectedStyle().Render(btn))
			} else {
				buttons = append(buttons, ui.UnselectedStyle().Render(btn))
			}
		}
		b.WriteString(strings.Join(buttons, "   "))
		b.WriteString("\n\n")
	}

	b.WriteString("ESC/q: Back\n")
	return b.String()
}
