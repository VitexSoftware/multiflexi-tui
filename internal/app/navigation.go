package app

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// updateSelectedHint updates the hint text based on current menu selection
func (m *Model) updateSelectedHint() {
	hints := []string{
		"View system dashboard with status information",
		"View and manage run templates with pagination controls",
		"View and manage running jobs with pagination controls",
		"Browse available MultiFlexi applications and their status",
		"View registered companies and their configuration",
		"View and manage credentials",
		"View and manage tokens",
		"View and manage users",
		"View and manage artifacts",
		"View and manage credential types",
		"View and manage credential prototypes",
		"View and manage company-application relations",
		"Manage encryption settings",
		"Manage the job queue",
		"Prune logs and jobs",
		"Browse available MultiFlexi commands and their documentation",
		"View help and documentation for using this interface",
		"Exit the MultiFlexi TUI application",
	}
	if m.menuCursor >= 0 && m.menuCursor < len(hints) {
		m.selectedHint = hints[m.menuCursor]
	} else {
		m.selectedHint = "Navigation: ←/→ to move, Enter to select"
	}
}

// handleMenuSelection handles menu item selection
func (m Model) handleMenuSelection() (tea.Model, tea.Cmd) {
	m.activeMenuItem = m.menuCursor

	switch m.menuCursor {
	case 0: // Status
		m.viewKind = ViewHome
		m.activeView = nil
		m.focus = false
		return m, nil
	case 1: // RunTemplates
		return m.switchToEntityList(ViewRunTemplates, newRunTemplatesList())
	case 2: // Jobs
		return m.switchToEntityList(ViewJobs, newJobsList())
	case 3: // Applications
		return m.switchToEntityList(ViewApplications, newApplicationsList())
	case 4: // Companies
		return m.switchToEntityList(ViewCompanies, newCompaniesList())
	case 5: // Credentials
		return m.switchToEntityList(ViewCredentials, newCredentialsList())
	case 6: // Tokens
		return m.switchToEntityList(ViewTokens, newTokensList())
	case 7: // Users
		return m.switchToEntityList(ViewUsers, newUsersList())
	case 8: // Artifacts
		return m.switchToEntityList(ViewArtifacts, newArtifactsList())
	case 9: // CredTypes
		return m.switchToEntityList(ViewCredTypes, newCredTypesList())
	case 10: // CrPrototypes
		return m.switchToEntityList(ViewCrPrototypes, newCrPrototypesList())
	case 11: // CompanyApps
		return m.switchToEntityList(ViewCompanyApps, newCompanyAppsList())
	case 12: // Encryption
		v := newEncryptionView(m.statusInfo.Encryption)
		m.viewKind = ViewEncryption
		m.activeView = v
		m.focus = false
		return m, nil
	case 13: // Queue
		return m.switchToEntityList(ViewQueue, newQueueList())
	case 14: // Prune
		v := newPruneView()
		m.viewKind = ViewPrune
		m.activeView = v
		m.focus = false
		return m, nil
	case 15: // Commands
		m.viewKind = ViewCommands
		m.activeView = nil
		m.focus = true
		return m, nil
	case 16: // Help
		m.viewKind = ViewHelp
		m.activeView = nil
		m.focus = false
		return m, m.loadHelpCmd("help")
	case 17: // Quit
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) switchToEntityList(kind ViewKind, list ui.EntityListModel) (tea.Model, tea.Cmd) {
	m.viewKind = kind
	m.activeView = list
	m.focus = false
	return m, list.Init()
}

// handleMouse processes mouse events
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.MouseLeft:
		if msg.Y == 0 {
			title := "MultiFlexi TUI"
			xPos := len(title) + 4 + 1
			startIdx := m.menuOffset
			if m.menuOffset > 0 {
				xPos += 4
			}
			for i := startIdx; i < len(m.menuItems); i++ {
				itemWidth := len(m.menuItems[i]) + 3
				if msg.X >= xPos && msg.X < xPos+itemWidth {
					m.menuCursor = i
					m.updateSelectedHint()
					return m.handleMenuSelection()
				}
				xPos += itemWidth
			}
		} else if msg.Y >= 3 {
			if m.focus {
				m.focus = false
			}
		}
	case tea.MouseWheelUp:
		if !m.focus {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			return m.Update(keyMsg)
		}
	case tea.MouseWheelDown:
		if !m.focus {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			return m.Update(keyMsg)
		}
	}
	return m, nil
}

// --- Entity list factory functions ---

var defaultHelpText = "↑/↓: navigate • ←/→: paginate • r: refresh • enter: view details"

func newRunTemplatesList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "📋 Run Templates",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 5, Field: "id"},
			{Header: "Name", Width: 25, Field: "name"},
			{Header: "App ID", Width: 10, Field: "app_id"},
			{Header: "Company", Width: 10, Field: "company"},
			{Header: "Status", Width: 8, Field: "status"},
			{Header: "Interval", Width: 10, Field: "interval"},
			{Header: "Executor", Width: 15, Field: "executor"},
		},
		Limit:        10,
		HelpText:     defaultHelpText,
		SupportsEdit: true,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetRunTemplates(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			templates := data.([]cli.RunTemplate)
			rows := make([]ui.TableRow, len(templates))
			for i, t := range templates {
				name := t.Name
				if name == "" {
					name = "<unnamed>"
				}
				status := "Active"
				if t.Active == 0 {
					status = "Inactive"
				}
				intervalText := t.Interv
				switch t.Interv {
				case "d":
					intervalText = "Daily"
				case "w":
					intervalText = "Weekly"
				case "m":
					intervalText = "Monthly"
				case "n":
					intervalText = "Never"
				}
				rows[i] = ui.TableRow{
					ID: t.ID,
					Values: map[string]interface{}{
						"id": t.ID, "name": name, "app_id": t.AppID,
						"company": t.CompanyID, "status": status,
						"interval": intervalText, "executor": t.Executor,
						"_full_data": t,
					},
				}
			}
			return rows
		},
	})
}

func newJobsList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "💼 Jobs",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 8, Field: "id"},
			{Header: "Command", Width: 25, Field: "command"},
			{Header: "Status", Width: 12, Field: "status"},
			{Header: "Schedule", Width: 20, Field: "schedule"},
		},
		Limit:        10,
		HelpText:     defaultHelpText,
		SupportsEdit: true,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetJobs(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			jobs := data.([]cli.Job)
			rows := make([]ui.TableRow, len(jobs))
			for i, j := range jobs {
				status := "Running"
				if j.PID == 0 {
					if j.Exitcode == -1 {
						status = "Scheduled"
					} else if j.Exitcode == 0 {
						status = "Success"
					} else {
						status = "Failed"
					}
				}
				schedule := j.Schedule
				if len(schedule) >= 16 {
					schedule = schedule[11:16]
				}
				rows[i] = ui.TableRow{
					ID: j.ID,
					Values: map[string]interface{}{
						"id": j.ID, "command": j.Command,
						"status": status, "schedule": schedule,
						"_full_data": j,
					},
				}
			}
			return rows
		},
	})
}

func newApplicationsList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "📦 Applications",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 5, Field: "id"},
			{Header: "Name", Width: 30, Field: "name"},
			{Header: "Version", Width: 15, Field: "version"},
			{Header: "Status", Width: 10, Field: "status"},
		},
		Limit:        10,
		HelpText:     defaultHelpText,
		SupportsEdit: true,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetApplications(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			apps := data.([]cli.Application)
			rows := make([]ui.TableRow, len(apps))
			for i, app := range apps {
				status := "Disabled"
				if app.Enabled == 1 {
					status = "Enabled"
				}
				rows[i] = ui.TableRow{
					ID: app.ID,
					Values: map[string]interface{}{
						"id": app.ID, "name": app.Name,
						"version": app.Version, "status": status,
						"_full_data": app,
					},
				}
			}
			return rows
		},
	})
}

func newCompaniesList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "🏢 Companies",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 6, Field: "id"},
			{Header: "Name", Width: 30, Field: "name"},
			{Header: "IC", Width: 15, Field: "ic"},
			{Header: "Status", Width: 10, Field: "status"},
		},
		Limit:        10,
		HelpText:     defaultHelpText,
		SupportsEdit: true,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetCompanies(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			companies := data.([]cli.Company)
			rows := make([]ui.TableRow, len(companies))
			for i, c := range companies {
				status := "Disabled"
				if c.Enabled == 1 {
					status = "Enabled"
				}
				rows[i] = ui.TableRow{
					ID: c.ID,
					Values: map[string]interface{}{
						"id": c.ID, "name": c.Name,
						"ic": c.IC, "status": status,
						"_full_data": c,
					},
				}
			}
			return rows
		},
	})
}

func newCredentialsList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "🔑 Credentials",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 8, Field: "id"},
			{Header: "Name", Width: 25, Field: "name"},
			{Header: "Company ID", Width: 12, Field: "company_id"},
			{Header: "Type ID", Width: 12, Field: "type_id"},
		},
		Limit:    10,
		HelpText: defaultHelpText,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetCredentials(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			creds := data.([]cli.Credential)
			rows := make([]ui.TableRow, len(creds))
			for i, c := range creds {
				rows[i] = ui.TableRow{
					ID: c.ID,
					Values: map[string]interface{}{
						"id": c.ID, "name": c.Name,
						"company_id": c.CompanyID, "type_id": c.CredentialTypeID,
						"_full_data": c,
					},
				}
			}
			return rows
		},
	})
}

func newTokensList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "🎟️ Tokens",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 8, Field: "id"},
			{Header: "User", Width: 25, Field: "user"},
			{Header: "Token", Width: 30, Field: "token"},
		},
		Limit:    10,
		HelpText: defaultHelpText,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetTokens(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			tokens := data.([]cli.Token)
			rows := make([]ui.TableRow, len(tokens))
			for i, t := range tokens {
				rows[i] = ui.TableRow{
					ID: t.ID,
					Values: map[string]interface{}{
						"id": t.ID, "user": t.User, "token": t.Token,
						"_full_data": t,
					},
				}
			}
			return rows
		},
	})
}

func newUsersList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "👤 Users",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 6, Field: "id"},
			{Header: "Login", Width: 20, Field: "login"},
			{Header: "First Name", Width: 20, Field: "firstname"},
			{Header: "Last Name", Width: 20, Field: "lastname"},
			{Header: "Email", Width: 30, Field: "email"},
		},
		Limit:    10,
		HelpText: defaultHelpText,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetUsers(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			users := data.([]cli.User)
			rows := make([]ui.TableRow, len(users))
			for i, u := range users {
				rows[i] = ui.TableRow{
					ID: u.ID,
					Values: map[string]interface{}{
						"id": u.ID, "login": u.Login,
						"firstname": u.Firstname, "lastname": u.Lastname,
						"email": u.Email, "_full_data": u,
					},
				}
			}
			return rows
		},
	})
}

func newArtifactsList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "📎 Artifacts",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 8, Field: "id"},
			{Header: "Job ID", Width: 10, Field: "job_id"},
			{Header: "Filename", Width: 40, Field: "filename"},
		},
		Limit:    10,
		HelpText: defaultHelpText,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetArtifacts(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			artifacts := data.([]cli.Artifact)
			rows := make([]ui.TableRow, len(artifacts))
			for i, a := range artifacts {
				rows[i] = ui.TableRow{
					ID: a.ID,
					Values: map[string]interface{}{
						"id": a.ID, "job_id": a.JobID,
						"filename": a.Filename, "_full_data": a,
					},
				}
			}
			return rows
		},
	})
}

func newCredTypesList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "🏷️ Credential Types",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 8, Field: "id"},
			{Header: "UUID", Width: 38, Field: "uuid"},
			{Header: "Name", Width: 30, Field: "name"},
		},
		Limit:    10,
		HelpText: defaultHelpText,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetCredTypes(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			types := data.([]cli.CredType)
			rows := make([]ui.TableRow, len(types))
			for i, t := range types {
				rows[i] = ui.TableRow{
					ID: t.ID,
					Values: map[string]interface{}{
						"id": t.ID, "uuid": t.UUID, "name": t.Name,
						"_full_data": t,
					},
				}
			}
			return rows
		},
	})
}

func newCrPrototypesList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "🧬 Credential Prototypes",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 8, Field: "id"},
			{Header: "Name", Width: 30, Field: "name"},
			{Header: "Description", Width: 40, Field: "description"},
			{Header: "Version", Width: 10, Field: "version"},
		},
		Limit:    10,
		HelpText: defaultHelpText,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetCrPrototypes(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			protos := data.([]cli.CrPrototype)
			rows := make([]ui.TableRow, len(protos))
			for i, p := range protos {
				desc := p.Description
				if len(desc) > 38 {
					desc = desc[:35] + "..."
				}
				rows[i] = ui.TableRow{
					ID: p.ID,
					Values: map[string]interface{}{
						"id": p.ID, "name": p.Name,
						"description": desc, "version": p.Version,
						"_full_data": p,
					},
				}
			}
			return rows
		},
	})
}

func newCompanyAppsList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "🔗 Company-Application Relations",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 8, Field: "id"},
			{Header: "Company ID", Width: 15, Field: "company_id"},
			{Header: "App ID", Width: 15, Field: "app_id"},
		},
		Limit:    10,
		HelpText: defaultHelpText,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetCompanyApps(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			apps := data.([]cli.CompanyApp)
			rows := make([]ui.TableRow, len(apps))
			for i, ca := range apps {
				rows[i] = ui.TableRow{
					ID: ca.ID,
					Values: map[string]interface{}{
						"id": ca.ID, "company_id": ca.CompanyID,
						"app_id": ca.AppID, "_full_data": ca,
					},
				}
			}
			return rows
		},
	})
}

func newQueueList() ui.EntityListModel {
	return ui.NewEntityListModel(ui.EntityListConfig{
		Title: "📬 Queue",
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 8, Field: "id"},
			{Header: "Job", Width: 8, Field: "job"},
			{Header: "Type", Width: 12, Field: "type"},
			{Header: "App", Width: 20, Field: "app"},
			{Header: "Company", Width: 20, Field: "company"},
			{Header: "After", Width: 20, Field: "after"},
		},
		Limit:    10,
		HelpText: defaultHelpText,
		Fetch: func(limit, offset int) (interface{}, error) {
			return cli.GetQueue(limit, offset)
		},
		Convert: func(data interface{}) []ui.TableRow {
			queue := data.([]cli.Queue)
			rows := make([]ui.TableRow, len(queue))
			for i, q := range queue {
				appName := q.AppName
				if len(appName) > 18 {
					appName = appName[:15] + "..."
				}
				companyName := q.CompanyName
				if len(companyName) > 18 {
					companyName = companyName[:15] + "..."
				}
				rows[i] = ui.TableRow{
					ID: q.ID,
					Values: map[string]interface{}{
						"id": q.ID, "job": q.Job,
						"type": q.ScheduleType, "app": appName,
						"company": companyName, "after": q.After,
						"_full_data": q,
					},
				}
			}
			return rows
		},
	})
}

// Ensure all ViewKind constants are used (suppress unused warnings)
var _ = fmt.Sprint(ViewHome, ViewRunTemplates, ViewJobs, ViewApplications, ViewCompanies,
	ViewCredentials, ViewTokens, ViewUsers, ViewArtifacts, ViewCredTypes,
	ViewCrPrototypes, ViewCompanyApps, ViewEncryption, ViewQueue, ViewPrune,
	ViewCommands, ViewHelp, ViewDetail, ViewEditor, ViewConfirmDelete)
