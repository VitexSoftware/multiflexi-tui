package entity

import (
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

const defaultHelp = "↑/↓: navigate • ←/→: paginate • r: refresh • enter: detail • e: edit • n: new • (entity actions shown below)"

// ListView is a generic list view driven by an EntityDef.
type ListView struct {
	client cli.Client
	def    *EntityDef
	table  *ui.TableWidget
}

// NewListView creates a new ListView for the given entity.
func NewListView(c cli.Client, def *EntityDef) *ListView {
	limit := def.Limit
	if limit == 0 {
		limit = 10
	}
	return &ListView{
		client: c,
		def:    def,
		table:  ui.NewTableWidget(def.Name, def.Columns, limit, defaultHelp),
	}
}

func (m *ListView) Init() tea.Cmd {
	m.table.SetLoading(true)
	return m.fetchCmd()
}

func (m *ListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case ui.DataLoadedMsg:
		rows := msg.Data.([]ui.TableRow)
		m.table.SetData(rows)
		return m, nil

	case ui.DataErrorMsg:
		m.table.SetError(msg.Err)
		return m, nil

	case tea.KeyMsg:
		for _, la := range m.def.ListActions {
			if la.Key != msg.String() {
				continue
			}
			if la.Confirm != "" {
				handler := la.Handler
				client := m.client
				label := la.Confirm
				return m, func() tea.Msg {
					return ui.ConfirmMsg{
						Label: label,
						Action: func() tea.Msg {
							cmd := handler(client)
							if cmd != nil {
								return cmd()
							}
							return nil
						},
					}
				}
			}
			return m, la.Handler(m.client)
		}

		refresh, nextPage, prevPage, openDetail, openEditor, openCreate := m.table.HandleKey(msg.String())

		if openDetail {
			row := m.table.SelectedRow()
			if row != nil && row.FullData != nil {
				detail := NewDetailView(m.client, m.def, row.FullData)
				return m, func() tea.Msg { return ui.NavigateToMsg{View: detail} }
			}
		}

		if openEditor && m.def.ToEditor != nil {
			row := m.table.SelectedRow()
			if row != nil && row.FullData != nil {
				editor := NewEditorView(m.client, m.def, row.FullData, false)
				return m, func() tea.Msg { return ui.NavigateToMsg{View: editor} }
			}
		}

		if openCreate && m.def.NewFields != nil {
			editor := NewEditorView(m.client, m.def, nil, true)
			return m, func() tea.Msg { return ui.NavigateToMsg{View: editor} }
		}

		if refresh || nextPage || prevPage {
			m.table.SetLoading(true)
			return m, m.fetchCmd()
		}
	}
	return m, nil
}

func (m *ListView) View() string {
	return m.table.View()
}

func (m *ListView) fetchCmd() tea.Cmd {
	limit := m.table.Limit()
	offset := m.table.Offset()
	fetch := m.def.Fetch
	client := m.client
	return func() tea.Msg {
		rows, err := fetch(client, limit, offset)
		if err != nil {
			return ui.DataErrorMsg{Err: err}
		}
		return ui.DataLoadedMsg{Data: rows}
	}
}
