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

	case ui.BackMsg:
		// Switch back to the previous view
		m.state = m.previousState

	case ui.OpenRunTemplateDetailMsg:
		// Switch to RunTemplate detail view
		m.previousState = m.state
		m.state = DetailView
		m.detailView.SetRunTemplate(msg.RunTemplate)
		return m, nil

	case ui.OpenApplicationDetailMsg:
		// Switch to Application detail view
		m.previousState = m.state
		m.state = DetailView
		m.detailView.SetApplication(msg.Application)
		return m, nil

	case ui.SaveApplicationMsg:
		err := cli.UpdateApplication(msg.App)
		if err != nil {
			m.statusMessage = fmt.Sprintf("Error saving application: %v", err)
		} else {
			m.statusMessage = fmt.Sprintf("Saved application: %s", msg.App.Name)
		}
		m.state = m.previousState
		return m, nil

	case ui.SaveRunTemplateMsg:
		err := cli.UpdateRunTemplate(msg.Template)
		if err != nil {
			m.statusMessage = fmt.Sprintf("Error saving run template: %v", err)
		} else {
			m.statusMessage = fmt.Sprintf("Saved run template: %s", msg.Template.Name)
		}
		m.state = m.previousState
		return m, nil

	case ui.StatusMessage:
		m.statusMessage = msg.Text
		return m, nil

	case ui.EditItemMsg:
		switch data := msg.Data.(type) {
		case cli.RunTemplate:
			m.previousState = m.state
			m.state = RunTemplateEditorView
			m.runTemplateEditor = ui.NewRunTemplateEditorModel(data)
		case cli.Application:
			m.previousState = m.state
			m.state = ApplicationEditorView
			m.applicationEditor = ui.NewApplicationEditorModel(data)
		}
		return m, nil

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
			m.focus = true // Focus on menu when switching to menu view
			return m, nil
		case "tab":
			m.focus = !m.focus
			return m, nil
		}

		// Handle menu or view-specific navigation based on focus
		if m.focus {
			// Menu navigation
			switch msg.String() {
			case "left", "h":
				if m.menuCursor > 0 {
					m.menuCursor--
					m.updateSelectedHint()
				}
			case "right", "l":
				if m.menuCursor < len(m.menuItems)-1 {
					m.menuCursor++
					m.updateSelectedHint()
				}
			case "enter", "space":
				return m.handleMenuSelection()
			}
		} else {
			// Content view navigation
			switch m.state {
			case RunTemplatesView:
				var cmd tea.Cmd
				runTemplatesModel, cmd := m.runTemplates.Update(msg)
				m.runTemplates = runTemplatesModel.(ui.RunTemplatesModel)
				return m, cmd
			case JobsView:
				var cmd tea.Cmd
				jobsModel, cmd := m.jobs.Update(msg)
				m.jobs = jobsModel.(ui.JobsModel)
				return m, cmd
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
			case CredentialsView:
				var cmd tea.Cmd
				credentialsModel, cmd := m.credentials.Update(msg)
				m.credentials = credentialsModel.(ui.CredentialsModel)
				return m, cmd
			case TokensView:
				var cmd tea.Cmd
				tokensModel, cmd := m.tokens.Update(msg)
				m.tokens = tokensModel.(ui.TokensModel)
				return m, cmd
			case UsersView:
				var cmd tea.Cmd
				usersModel, cmd := m.users.Update(msg)
				m.users = usersModel.(ui.UsersModel)
				return m, cmd
			case ArtifactsView:
				var cmd tea.Cmd
				artifactsModel, cmd := m.artifacts.Update(msg)
				m.artifacts = artifactsModel.(ui.ArtifactsModel)
				return m, cmd
			case CredTypesView:
				var cmd tea.Cmd
				credTypesModel, cmd := m.credTypes.Update(msg)
				m.credTypes = credTypesModel.(ui.CredTypesModel)
				return m, cmd
			case CompanyAppsView:
				var cmd tea.Cmd
				companyAppsModel, cmd := m.companyApps.Update(msg)
				m.companyApps = companyAppsModel.(ui.CompanyAppsModel)
				return m, cmd
			case EncryptionView:
				var cmd tea.Cmd
				encryptionModel, cmd := m.encryption.Update(msg)
				m.encryption = encryptionModel.(ui.EncryptionModel)
				return m, cmd
			case QueueView:
				var cmd tea.Cmd
				queueModel, cmd := m.queue.Update(msg)
				m.queue = queueModel.(ui.QueueModel)
				return m, cmd
			case PruneView:
				var cmd tea.Cmd
				pruneModel, cmd := m.prune.Update(msg)
				m.prune = pruneModel.(ui.PruneModel)
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
			case DetailView:
				var cmd tea.Cmd
				detailModel, cmd := m.detailView.Update(msg)
				m.detailView = detailModel.(ui.DetailViewModel)
				return m, cmd
			case RunTemplateEditorView:
				var cmd tea.Cmd
				editorModel, cmd := m.runTemplateEditor.Update(msg)
				m.runTemplateEditor = editorModel.(ui.RunTemplateEditorModel)
				return m, cmd
			case ApplicationEditorView:
				var cmd tea.Cmd
				editorModel, cmd := m.applicationEditor.Update(msg)
				m.applicationEditor = editorModel.(ui.ApplicationEditorModel)
				return m, cmd
			}
		}
	}
// ... (rest of the code)
	case DetailView:
		content = m.detailView.View()
		statusPanel = m.renderHelpFooter()
	case RunTemplateEditorView:
		content = m.runTemplateEditor.View()
		statusPanel = m.renderHelpFooter()
	case ApplicationEditorView:
		content = m.applicationEditor.View()
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
	case CredentialsView:
		content = m.credentials.View()
		statusPanel = m.renderHelpFooter()
	case TokensView:
		content = m.tokens.View()
		statusPanel = m.renderHelpFooter()
	case UsersView:
		content = m.users.View()
		statusPanel = m.renderHelpFooter()
	case ArtifactsView:
		content = m.artifacts.View()
		statusPanel = m.renderHelpFooter()
	case CredTypesView:
		content = m.credTypes.View()
		statusPanel = m.renderHelpFooter()
	case CompanyAppsView:
		content = m.companyApps.View()
		statusPanel = m.renderHelpFooter()
	case EncryptionView:
		content = m.encryption.View()
		statusPanel = m.renderHelpFooter()
	case QueueView:
		content = m.queue.View()
		statusPanel = m.renderHelpFooter()
	case PruneView:
		content = m.prune.View()
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
	style := ui.GetUnselectedItemStyle()
	if m.focus {
		style = ui.GetSelectedItemStyle()
	}

	for i, item := range m.menuItems {
		if i == m.menuCursor {
			menuItems = append(menuItems, style.Render(" "+item+" "))
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
func (m *Model) renderHelpFooter() string {
	width := m.width
	if width == 0 {
		width = 80 // Default width if not set
	}

	separator := strings.Repeat("─", width)
	var helpLine string
	if m.focus {
		helpLine = ui.GetFooterStyle().Render("←/→: navigate menu • enter: select • tab: switch to content • q: quit")
	} else {
		helpLine = ui.GetFooterStyle().Render("↑/↓: navigate list • ←/→: paginate • tab: switch to menu • q: quit")
	}

	statusLine := ""
	if m.statusMessage != "" {
		statusLine = ui.GetFooterStyle().Render(m.statusMessage)
		m.statusMessage = ""
	}

	return separator + "\n" + statusLine + "\n" + helpLine
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
	case 5: // Credentials
		m.selectedHint = "View and manage credentials"
	case 6: // Tokens
		m.selectedHint = "View and manage tokens"
	case 7: // Users
		m.selectedHint = "View and manage users"
	case 8: // Artifacts
		m.selectedHint = "View and manage artifacts"
	case 9: // CredTypes
		m.selectedHint = "View and manage credential types"
	case 10: // CompanyApps
		m.selectedHint = "View and manage company-application relations"
	case 11: // Encryption
		m.selectedHint = "Manage encryption settings"
	case 12: // Queue
		m.selectedHint = "Manage the job queue"
	case 13: // Prune
		m.selectedHint = "Prune logs and jobs"
	case 14: // Commands
		m.selectedHint = "Browse available MultiFlexi commands and their documentation"
	case 15: // Help
		m.selectedHint = "View help and documentation for using this interface"
	case 16: // Quit
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
	case 5: // Credentials
		m.state = CredentialsView
		// Reset credentials model and trigger loading
		m.credentials = ui.NewCredentialsModel()
		return m, m.credentials.Init()
	case 6: // Tokens
		m.state = TokensView
		// Reset tokens model and trigger loading
		m.tokens = ui.NewTokensModel()
		return m, m.tokens.Init()
	case 7: // Users
		m.state = UsersView
		// Reset users model and trigger loading
		m.users = ui.NewUsersModel()
		return m, m.users.Init()
	case 8: // Artifacts
		m.state = ArtifactsView
		// Reset artifacts model and trigger loading
		m.artifacts = ui.NewArtifactsModel()
		return m, m.artifacts.Init()
	case 9: // CredTypes
		m.state = CredTypesView
		// Reset credTypes model and trigger loading
		m.credTypes = ui.NewCredTypesModel()
		return m, m.credTypes.Init()
	case 10: // CompanyApps
		m.state = CompanyAppsView
		// Reset companyApps model and trigger loading
		m.companyApps = ui.NewCompanyAppsModel()
		return m, m.companyApps.Init()
	case 11: // Encryption
		m.state = EncryptionView
		m.encryption = ui.NewEncryptionModel(m.statusInfo.Encryption)
		return m, nil
	case 12: // Queue
		m.state = QueueView
		// Reset queue model and trigger loading
		m.queue = ui.NewQueueModel()
		return m, m.queue.Init()
	case 13: // Prune
		m.state = PruneView
		m.prune = ui.NewPruneModel()
		return m, nil
	case 14: // Commands
		m.state = MenuView
		return m, nil
	case 15: // Help
		m.state = HelpView
		return m, m.loadHelpCmd("help")
	case 16: // Quit
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

// Run starts the Bubbletea program
func Run() error {
	model := NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to start program: %w", err)
	}

	return nil
}
