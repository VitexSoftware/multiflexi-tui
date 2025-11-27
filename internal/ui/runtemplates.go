package ui

import (
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
)

// RunTemplatesModel represents the run templates listing view
type RunTemplatesModel struct {
	table  *TableWidget
	width  int
	height int
}

// runTemplatesLoadedMsg is sent when run templates are loaded
type runTemplatesLoadedMsg struct {
	runTemplates []cli.RunTemplate
}

// runTemplatesErrorMsg is sent when there's an error loading run templates
type runTemplatesErrorMsg struct {
	err error
}

// OpenRunTemplateDetailMsg is sent when a run template should be opened in detail view
type OpenRunTemplateDetailMsg struct {
	RunTemplate cli.RunTemplate
}

// NewRunTemplatesModel creates a new run templates model
func NewRunTemplatesModel() RunTemplatesModel {
	config := TableConfig{
		Title: "📋 Run Templates",
		Columns: []TableColumn{
			{Header: "ID", Width: 5, Field: "id"},
			{Header: "Name", Width: 25, Field: "name"},
			{Header: "App ID", Width: 10, Field: "app_id"},
			{Header: "Company", Width: 10, Field: "company"},
			{Header: "Status", Width: 8, Field: "status"},
			{Header: "Interval", Width: 10, Field: "interval"},
			{Header: "Executor", Width: 15, Field: "executor"},
		},
		Limit:    10,
		HelpText: "↑↓: navigate • ←→/PgUp/PgDn: prev/next page • r: refresh",
	}

	return RunTemplatesModel{
		table: NewTableWidget(config),
	}
}

// convertRunTemplatesToTableRows converts CLI RunTemplate data to table rows
func convertRunTemplatesToTableRows(runTemplates []cli.RunTemplate) []TableRow {
	rows := make([]TableRow, len(runTemplates))

	for i, template := range runTemplates {
		// Handle null name
		name := template.Name
		if name == "" {
			name = "<unnamed>"
		}

		// Status - plain text only
		status := "Active"
		if template.Active == 0 {
			status = "Inactive"
		}

		// Interval display
		intervalText := template.Interv
		switch template.Interv {
		case "d":
			intervalText = "Daily"
		case "w":
			intervalText = "Weekly"
		case "m":
			intervalText = "Monthly"
		case "n":
			intervalText = "Never"
		}

		rows[i] = TableRow{
			ID: template.ID,
			Values: map[string]interface{}{
				"id":         template.ID,
				"name":       name,
				"app_id":     template.AppID,
				"company":    template.CompanyID,
				"status":     status,
				"interval":   intervalText,
				"executor":   template.Executor,
				"_full_data": template, // Store full template data
			},
		}
	}

	return rows
} // Init initializes the run templates model
func (m RunTemplatesModel) Init() tea.Cmd {
	m.table.SetLoading(true)
	return m.loadRunTemplatesCmd()
}

// Update handles messages for the run templates model
func (m RunTemplatesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case runTemplatesLoadedMsg:
		rows := convertRunTemplatesToTableRows(msg.runTemplates)
		m.table.SetData(rows)
		return m, nil

	case runTemplatesErrorMsg:
		m.table.SetError(msg.err)
		return m, nil

	case tea.KeyMsg:
		needsRefresh, needsNextPage, needsPrevPage, openDetail := m.table.HandleKeypress(msg.String())

		if openDetail {
			// Get the selected row and extract the full RunTemplate data
			selectedRow := m.table.GetSelectedRow()
			if selectedRow != nil {
				if fullData, exists := selectedRow.Values["_full_data"]; exists {
					if template, ok := fullData.(cli.RunTemplate); ok {
						return m, func() tea.Msg {
							return OpenRunTemplateDetailMsg{RunTemplate: template}
						}
					}
				}
			}
		}

		if needsRefresh {
			m.table.SetLoading(true)
			return m, m.loadRunTemplatesCmd()
		}

		if needsNextPage || needsPrevPage {
			m.table.SetLoading(true)
			return m, m.loadRunTemplatesCmd()
		}
	}

	return m, nil
}

// View renders the run templates model
func (m RunTemplatesModel) View() string {
	return m.table.View()
}

// loadRunTemplatesCmd returns a command that loads run templates from CLI
func (m RunTemplatesModel) loadRunTemplatesCmd() tea.Cmd {
	return func() tea.Msg {
		runTemplates, err := cli.GetRunTemplates(m.table.GetLimit(), m.table.GetOffset())
		if err != nil {
			return runTemplatesErrorMsg{err: err}
		}
		return runTemplatesLoadedMsg{runTemplates: runTemplates}
	}
}
