
package app

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


// Init initializes the application model
func (m Model) Init() tea.Cmd {
	return m.loadStatusCmd() // Only load status, not jobs initially
}

// Update handles messages for the application
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Forward size message to current view

	case StatusLoadedMsg:
		m.statusInfo = msg.status

	case ui.ShowHelpMsg:
		// Switch to help view and load help content
		m.state = HelpView
		return m, m.loadHelpCmd(msg.Command)

	case ui.ShowMenuMsg:
		// Switch to menu view
		m.state = MenuView

	case ui.BackToMenuMsg:
		// Switch back to menu view
		m.state = MenuView

	case ui.OpenRunTemplateDetailMsg:
		// Switch to RunTemplate detail view
		m.state = RunTemplateDetailView
		return m, m.loadRunTemplateDetail(msg.RunTemplate)

	case helpLoadedMsg:
		// Update viewer with help content
		m.viewer.SetContent(msg.command, msg.content)
		// Removed stray return statement outside of any function

	case helpErrorMsg:
		// Display error in viewer
		m.viewer.SetError(msg.err)
		return m, nil

	case tea.KeyMsg:
		// Handle global keys that apply to all views
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "f10":
			m.state = MenuView
			return m, nil
		}

		// Handle view-specific key navigation
		switch m.state {
		case HomeView:
			// Handle menu navigation for home view
			switch msg.String() {
			case "left", "h":
				if m.menuCursor > 0 {
					m.menuCursor--
					m.updateSelectedHint()
				}
				return m, nil
			case "right", "l":
				if m.menuCursor < len(m.menuItems)-1 {
					m.menuCursor++
					m.updateSelectedHint()
				}
				return m, nil
			case "enter", "space":
				return m.handleMenuSelection()
				}
			}
		case ApplicationsView:
			var cmd tea.Cmd
			applicationsModel, cmd := m.applications.Update(msg)
			m.applications = applicationsModel.(ui.ApplicationsModel)
			return m, cmd
		case CompaniesView:
			var cmd tea.Cmd
			companiesModel, cmd := m.companies.Update(msg)
			m.companies = companiesModel.(ui.CompaniesModel)
			return m, cmd
		case MenuView:
			var cmd tea.Cmd
			menuModel, cmd := m.menu.Update(msg)
			m.menu = menuModel.(ui.MenuModel)
			return m, cmd
		case HelpView:
			var cmd tea.Cmd
			viewerModel, cmd := m.viewer.Update(msg)
			m.viewer = viewerModel.(ui.ViewerModel)
			return m, cmd
		}
	}

	// Forward other messages to current view

	return m, nil
}

// View renders the current view with the new layout
func (m Model) View() string {
	// Top menu bar
	menuBar := m.renderMenuBar()

	// Content area
	var content string
	var statusPanel string

	switch m.state {
	case HomeView:
		content = m.renderSystemStatus()
		statusPanel = m.renderHelpFooter()
	case RunTemplatesView:
		content = m.runTemplates.View()
		statusPanel = m.renderHelpFooter()
	case RunTemplateDetailView:
		content = m.runTemplateDetail.View()
		statusPanel = m.renderHelpFooter()
	case JobsView:
		content = m.jobs.View()
		statusPanel = m.renderHelpFooter()
	case ApplicationsView:
		content = m.applications.View()
		statusPanel = m.renderHelpFooter()
	case CompaniesView:
		content = m.companies.View()
		statusPanel = m.renderHelpFooter()
	case MenuView:
		content = m.menu.View()
		statusPanel = m.renderHelpFooter()
	case HelpView:
		content = m.viewer.View()
		statusPanel = m.renderHelpFooter()
	default:
		content = "Unknown view state"
		statusPanel = m.renderHelpFooter()
	}

	return menuBar + content + statusPanel
}

// renderMenuBar renders the top menu bar with hints
func (m Model) renderMenuBar() string {
	var menuItems []string

	for i, item := range m.menuItems {
		if i == m.menuCursor {
			menuItems = append(menuItems, ui.GetSelectedItemStyle().Render(" "+item+" "))
		} else {
			menuItems = append(menuItems, ui.GetUnselectedItemStyle().Render(" "+item+" "))
		}
	}

	width := m.width
	if width == 0 {
		width = 80 // Default width if not set
	}

	menuLine := ui.GetTitleStyle().Render("MultiFlexi TUI") + "  " + strings.Join(menuItems, " ")
	hintLine := ui.GetItemDescriptionStyle().Render(m.selectedHint)
	separator := strings.Repeat("─", width)

	return menuLine + "\n" + hintLine + "\n" + separator + "\n"
}

// renderHelpFooter renders just the help footer
func (m Model) renderHelpFooter() string {
	width := m.width
	if width == 0 {
		width = 80 // Default width if not set
	}

	separator := strings.Repeat("─", width)
	helpLine := ui.GetFooterStyle().Render("←/→: navigate menu • enter: select • r: refresh • q: quit")

	return separator + "\n" + helpLine
}

// renderSystemStatus renders the system status as main content
func (m Model) renderSystemStatus() string {
	var content strings.Builder

	// Title
	content.WriteString(ui.GetTitleStyle().Render("🖥️  MultiFlexi System Dashboard"))
	content.WriteString("\n\n")

	if m.statusInfo == nil {
		content.WriteString(ui.GetItemDescriptionStyle().Render("Loading system status..."))
		content.WriteString("\n")
	} else {
		// Create status rows with emoticons and colors
		rows := []struct {
			icon   string
			label  string
			value  string
			status string
		}{
			{"🔧", "CLI Version", m.statusInfo.Version, "info"},
			{"👤", "User", m.statusInfo.User, "info"},
			{"🐘", "PHP", m.statusInfo.PHP, "info"},
			{"💻", "OS", m.statusInfo.OS, "info"},
			{"🏢", "Companies", fmt.Sprintf("%d", m.statusInfo.Companies), "info"},
			{"📱", "Applications", fmt.Sprintf("%d", m.statusInfo.Apps), "info"},
			{"📄", "Templates", fmt.Sprintf("%d", m.statusInfo.Templates), "info"},
			{"⚙️", "Executor", m.statusInfo.Executor, m.statusInfo.Executor},
			{"📅", "Scheduler", m.statusInfo.Scheduler, m.statusInfo.Scheduler},
			{"🔐", "Encryption", m.statusInfo.Encryption, m.statusInfo.Encryption},
			{"📊", "Zabbix", m.statusInfo.Zabbix, m.statusInfo.Zabbix},
			{"📈", "Telemetry", m.statusInfo.Telemetry, m.statusInfo.Telemetry},
		}

		// Calculate column widths
		labelWidth := 15

		// Render each row
		for _, row := range rows {
			var valueStyle lipgloss.Style
			switch row.status {
			case "active":
				valueStyle = ui.GetActiveStatusStyle()
			case "disabled":
				valueStyle = ui.GetDisabledStatusStyle()
			default:
				valueStyle = ui.GetItemDescriptionStyle()
			}

			line := fmt.Sprintf("%s %-*s %s",
				row.icon,
				labelWidth,
				row.label+":",
				valueStyle.Render(row.value))

			content.WriteString(line)
			content.WriteString("\n")
		}

		// Add database info if available
		if m.statusInfo.Database != "" {
			content.WriteString("\n")
			content.WriteString("🗄️  Database Information:")
			content.WriteString("\n")
			// Truncate database info if too long
			dbInfo := m.statusInfo.Database
			if len(dbInfo) > 80 {
				dbInfo = dbInfo[:77] + "..."
			}
			content.WriteString(ui.GetItemDescriptionStyle().Render("   " + dbInfo))
			content.WriteString("\n")
		}
	}

	return content.String()
}

// helpLoadedMsg is sent when help content is loaded successfully
type helpLoadedMsg struct {
	command string
	content string
}

// helpErrorMsg is sent when there's an error loading help content
type helpErrorMsg struct {
	command string
	err     error
}

// loadHelpCmd returns a command that loads help for a command
func (m Model) loadHelpCmd(command string) tea.Cmd {
	return func() tea.Msg {
		content, err := cli.GetCommandHelp(command)
		if err != nil {
			return helpErrorMsg{command: command, err: err}
		}
		return helpLoadedMsg{command: command, content: content}
	}
}

// loadStatusCmd returns a command that loads status from CLI
func (m Model) loadStatusCmd() tea.Cmd {
	return func() tea.Msg {
		status, err := cli.GetStatusInfo()
		if err != nil {
			// Create minimal error status
			errorLen := len(err.Error())
			if errorLen > 20 {
				errorLen = 20
			}
			return StatusLoadedMsg{status: &cli.StatusInfo{
				Version: "Error",
				User:    err.Error()[:errorLen],
			}}
		}
		return StatusLoadedMsg{status: status}
	}
}

// updateSelectedHint updates the hint text based on current menu selection
func (m *Model) updateSelectedHint() {
	switch m.menuCursor {
	case 0: // Status
		m.selectedHint = "View system dashboard with status information"
	case 1: // RunTemplates
		m.selectedHint = "View and manage run templates with pagination controls"
	case 2: // Jobs
		m.selectedHint = "View and manage running jobs with pagination controls"
	case 3: // Applications
		m.selectedHint = "Browse available MultiFlexi applications and their status"
	case 4: // Companies
		m.selectedHint = "View registered companies and their configuration"
	case 5: // Commands
		m.selectedHint = "Browse available MultiFlexi commands and their documentation"
	case 6: // Help
		m.selectedHint = "View help and documentation for using this interface"
	case 7: // Quit
		m.selectedHint = "Exit the MultiFlexi TUI application"
	default:
		m.selectedHint = "Navigation: ←/→ to move, Enter to select"
	}
}

// handleMenuSelection handles menu item selection
func (m Model) handleMenuSelection() (tea.Model, tea.Cmd) {
	switch m.menuCursor {
	case 0: // Status
		m.state = HomeView
		return m, nil
	case 1: // RunTemplates
		m.state = RunTemplatesView
		// Reset runTemplates model and trigger loading
		m.runTemplates = ui.NewRunTemplatesModel()
		return m, m.runTemplates.Init()
	case 2: // Jobs
		m.state = JobsView
		// Reset jobs model and trigger loading
		m.jobs = ui.NewJobsModel()
		return m, m.jobs.Init()
	case 3: // Applications
		m.state = ApplicationsView
		// Reset applications model and trigger loading
		m.applications = ui.NewApplicationsModel()
		return m, m.applications.Init()
	case 4: // Companies
		m.state = CompaniesView
		// Reset companies model and trigger loading
		m.companies = ui.NewCompaniesModel()
		return m, m.companies.Init()
	case 5: // Commands
		m.state = MenuView
		return m, nil
	case 6: // Help
		m.state = HelpView
		return m, m.loadHelpCmd("help")
	case 7: // Quit
		return m, tea.Quit
	default:
		return m, nil
	}
}

// handleDetailAction handles action commands from detail views
func (m Model) handleDetailAction(actionCommand string) tea.Cmd {
	switch actionCommand {
	case "edit":
		// Handle edit action - for now just show a message
		// In a real implementation, this would open an edit form
		return func() tea.Msg {
			return nil // Placeholder
		}
	case "schedule":
		// Handle schedule action - for now just show a message
		// In a real implementation, this would open a scheduling dialog
		return func() tea.Msg {
			return nil // Placeholder
		}
	default:
		return nil
	}
}

// loadRunTemplateDetail loads the detail view for a RunTemplate
func (m Model) loadRunTemplateDetail(template cli.RunTemplate) tea.Cmd {
	return func() tea.Msg {
		// Convert RunTemplate to detail fields
		fields := []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", template.ID)},
			{Label: "Name", Value: template.Name},
			{Label: "App ID", Value: fmt.Sprintf("%d", template.AppID)},
			{Label: "Company ID", Value: fmt.Sprintf("%d", template.CompanyID)},
			{Label: "Interval", Value: template.Interv},
			{Label: "Active", Value: fmt.Sprintf("%d", template.Active)},
			{Label: "Executor", Value: template.Executor},
			{Label: "Cron", Value: template.Cron},
			{Label: "Last Schedule", Value: template.LastSchedule},
			{Label: "Next Schedule", Value: template.NextSchedule},
		}

		m.runTemplateDetail.SetData(fields, template)
		return nil
	}
}

// Run starts the Bubbletea program
func Run() error {
	model, err := NewModel()
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	_, err = p.Run()
	if err != nil {
		return fmt.Errorf("failed to start program: %w", err)
	}

	return nil
}
