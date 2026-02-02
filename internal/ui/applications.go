package ui

import (
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
)

// ApplicationsModel represents the applications listing view
type ApplicationsModel struct {
	table *TableWidget
	width  int
	height int
}

// applicationsLoadedMsg is sent when applications are loaded
type applicationsLoadedMsg struct {
	apps []cli.Application
}

// applicationsErrorMsg is sent when there's an error loading applications
type applicationsErrorMsg struct {
	err error
}

// OpenApplicationDetailMsg is sent when an application should be opened in detail view
type OpenApplicationDetailMsg struct {
	Application cli.Application
}

// NewApplicationsModel creates a new applications model
func NewApplicationsModel() ApplicationsModel {
	config := TableConfig{
		Title: "📦 Applications",
		Columns: []TableColumn{
			{Header: "ID", Width: 5, Field: "id"},
			{Header: "Name", Width: 30, Field: "name"},
			{Header: "Version", Width: 15, Field: "version"},
			{Header: "Status", Width: 10, Field: "status"},
		},
		Limit:    10,
		HelpText: "↑/↓: navigate • ←/→: paginate • r: refresh • enter: view details",
	}

	return ApplicationsModel{
		table: NewTableWidget(config),
	}
}

// convertApplicationsToTableRows converts CLI Application data to table rows
func convertApplicationsToTableRows(apps []cli.Application) []TableRow {
	rows := make([]TableRow, len(apps))

	for i, app := range apps {
		// Status - plain text only
		status := "Disabled"
		if app.Enabled == 1 {
			status = "Enabled"
		}

		rows[i] = TableRow{
			ID: app.ID,
			Values: map[string]interface{}{
				"id":         app.ID,
				"name":       app.Name,
				"version":    app.Version,
				"status":     status,
				"_full_data": app, // Store full application data
			},
		}
	}

	return rows
}

// Init initializes the applications model
func (m ApplicationsModel) Init() tea.Cmd {
	m.table.SetLoading(true)
	return m.loadApplicationsCmd()
}

// Update handles messages for the applications model
func (m ApplicationsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case applicationsLoadedMsg:
		rows := convertApplicationsToTableRows(msg.apps)
		m.table.SetData(rows)
		return m, nil

	case applicationsErrorMsg:
		m.table.SetError(msg.err)
		return m, nil

	case tea.KeyMsg:
		needsRefresh, needsNextPage, needsPrevPage, openDetail := m.table.HandleKeypress(msg.String())

		if openDetail {
			// Get the selected row and extract the full Application data
			selectedRow := m.table.GetSelectedRow()
			if selectedRow != nil {
				if fullData, exists := selectedRow.Values["_full_data"]; exists {
					if app, ok := fullData.(cli.Application); ok {
						return m, func() tea.Msg {
							return OpenApplicationDetailMsg{Application: app}
						}
					}
				}
			}
		}

		if needsRefresh {
			m.table.SetLoading(true)
			return m, m.loadApplicationsCmd()
		}

		if needsNextPage || needsPrevPage {
			m.table.SetLoading(true)
			return m, m.loadApplicationsCmd()
		}
	}

	return m, nil
}

// View renders the applications model
func (m ApplicationsModel) View() string {
	return m.table.View()
}

// loadApplicationsCmd returns a command that loads applications from CLI
func (m ApplicationsModel) loadApplicationsCmd() tea.Cmd {
	return func() tea.Msg {
		apps, err := cli.GetApplications(m.table.GetLimit(), m.table.GetOffset())
		if err != nil {
			return applicationsErrorMsg{err: err}
		}
		return applicationsLoadedMsg{apps: apps}
	}
}
