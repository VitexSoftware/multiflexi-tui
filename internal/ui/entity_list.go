package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// RowConverter converts raw data items into TableRows.
// The function receives a slice (as interface{}) and must return []TableRow.
type RowConverter func(data interface{}) []TableRow

// FetchFunc fetches a page of data from the CLI. Returns (data, error).
type FetchFunc func(limit, offset int) (interface{}, error)

// EntityListConfig defines how an entity listing behaves.
type EntityListConfig struct {
	Title        string
	Columns      []TableColumn
	Limit        int
	HelpText     string
	Fetch        FetchFunc
	Convert      RowConverter
	SupportsEdit bool // whether 'e' key opens an editor
}

// EntityListModel is a generic, reusable listing model for any entity type.
type EntityListModel struct {
	config  EntityListConfig
	table   *TableWidget
	rawData interface{} // the last fetched raw data slice
	width   int
	height  int
}

// entityDataLoadedMsg carries the raw fetched data
type entityDataLoadedMsg struct {
	data interface{}
}

// entityDataErrorMsg carries a fetch error
type entityDataErrorMsg struct {
	err error
}

// NewEntityListModel creates a new entity list model from config.
func NewEntityListModel(config EntityListConfig) EntityListModel {
	tableConfig := TableConfig{
		Title:    config.Title,
		Columns:  config.Columns,
		Limit:    config.Limit,
		HelpText: config.HelpText,
	}
	return EntityListModel{
		config: config,
		table:  NewTableWidget(tableConfig),
	}
}

// Init starts the initial data load.
func (m EntityListModel) Init() tea.Cmd {
	m.table.SetLoading(true)
	return m.fetchCmd()
}

// Update handles all messages for the entity list.
func (m EntityListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case entityDataLoadedMsg:
		m.rawData = msg.data
		rows := m.config.Convert(msg.data)
		m.table.SetData(rows)
		return m, nil

	case entityDataErrorMsg:
		m.table.SetError(msg.err)
		return m, nil

	case tea.KeyMsg:
		needsRefresh, needsNextPage, needsPrevPage, openDetail, openEditor := m.table.HandleKeypress(msg.String())

		if openDetail {
			row := m.table.GetSelectedRow()
			if row != nil {
				if fullData, ok := row.Values["_full_data"]; ok {
					return m, func() tea.Msg { return OpenDetailMsg{Data: fullData} }
				}
			}
		}

		if openEditor && m.config.SupportsEdit {
			row := m.table.GetSelectedRow()
			if row != nil {
				if fullData, ok := row.Values["_full_data"]; ok {
					return m, func() tea.Msg { return EditItemMsg{Data: fullData} }
				}
			}
		}

		if needsRefresh || needsNextPage || needsPrevPage {
			m.table.SetLoading(true)
			return m, m.fetchCmd()
		}
	}

	return m, nil
}

// View renders the entity list.
func (m EntityListModel) View() string {
	return m.table.View()
}

// fetchCmd returns a tea.Cmd that fetches data using the configured FetchFunc.
func (m EntityListModel) fetchCmd() tea.Cmd {
	limit := m.table.GetLimit()
	offset := m.table.GetOffset()
	fetch := m.config.Fetch
	return func() tea.Msg {
		data, err := fetch(limit, offset)
		if err != nil {
			return entityDataErrorMsg{err: err}
		}
		return entityDataLoadedMsg{data: data}
	}
}

// --- Helpers for max() used by old code ---

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// --- Convenience: truncate string for display ---

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen > 3 {
		return s[:maxLen-3] + "..."
	}
	return s[:maxLen]
}

// --- Convenience: format nullable string pointer ---

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// --- Convenience: format int as string ---

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
