package ui

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
)

// DetailViewModel represents the detail view for an item
type DetailViewModel struct {
	widget *DetailWidget
	width  int
	height int
}

// NewDetailViewModel creates a new detail view model
func NewDetailViewModel() DetailViewModel {
	return DetailViewModel{}
}

// SetContent populates the detail view with data from a RunTemplate
func (m *DetailViewModel) SetRunTemplate(template cli.RunTemplate) {
	config := DetailConfig{
		Title: fmt.Sprintf("Run Template: %s", template.Name),
		Actions: []DetailAction{
			{Label: "Edit", Key: "e", Command: "edit"},
			{Label: "Schedule", Key: "s", Command: "schedule"},
			{Label: "Clone", Key: "c", Command: "clone"},
			{Label: "Delete", Key: "d", Command: "delete"},
		},
	}
	m.widget = NewDetailWidget(config)

	fields := []DetailField{
		{Label: "ID", Value: fmt.Sprintf("%d", template.ID)},
		{Label: "Name", Value: template.Name},
		{Label: "App ID", Value: fmt.Sprintf("%d", template.AppID)},
		{Label: "Company ID", Value: fmt.Sprintf("%d", template.CompanyID)},
		{Label: "Status", Value: fmt.Sprintf("%d", template.Active)},
		{Label: "Interval", Value: template.Interv},
		{Label: "Executor", Value: template.Executor},
	}
	m.widget.SetData(fields, template)
}

// SetContent populates the detail view with data from an Application
func (m *DetailViewModel) SetApplication(app cli.Application) {
	config := DetailConfig{
		Title: fmt.Sprintf("Application: %s", app.Name),
		Actions: []DetailAction{
			{Label: "Edit", Key: "e", Command: "edit"},
			{Label: "Clone", Key: "c", Command: "clone"},
			{Label: "Delete", Key: "d", Command: "delete"},
		},
	}
	m.widget = NewDetailWidget(config)

	fields := []DetailField{
		{Label: "ID", Value: fmt.Sprintf("%d", app.ID)},
		{Label: "Name", Value: app.Name},
		{Label: "Version", Value: app.Version},
		{Label: "Enabled", Value: fmt.Sprintf("%d", app.Enabled)},
	}
	m.widget.SetData(fields, app)
}

// Init initializes the detail view model
func (m DetailViewModel) Init() tea.Cmd {
	return nil
}

// EditItemMsg is sent when an item should be edited
type EditItemMsg struct {
	Data interface{}
}

// Update handles messages for the detail view model
func (m DetailViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		action, goBack := m.widget.HandleKeypress(msg.String())
		if goBack {
			// This will be handled by the main model to switch views
			return m, func() tea.Msg { return BackMsg{} }
		}
		if action == "edit" {
			return m, func() tea.Msg {
				return EditItemMsg{Data: m.widget.GetData()}
			}
		}
		if action != "" {
			// Handle actions like "clone", "delete"
			// For now, we'll just send a message
			return m, func() tea.Msg {
				return StatusMessage{Text: fmt.Sprintf("Action: %s on %T", action, m.widget.GetData())}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the detail view model
func (m DetailViewModel) View() string {
	if m.widget == nil {
		return "No item selected."
	}
	return m.widget.View()
}

// BackMsg is sent to go back to the previous view
type BackMsg struct{}

// StatusMessage is sent to display a status message to the user
type StatusMessage struct {
	Text string
}
