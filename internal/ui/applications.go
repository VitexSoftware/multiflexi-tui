package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ApplicationsModel represents the applications listing screen
type ApplicationsModel struct {
	apps    []cli.Application
	offset  int
	limit   int
	loading bool
	err     error
	width   int
	height  int
	cursor  int
	hasMore bool
	hasPrev bool
}

// applicationsLoadedMsg is sent when applications are loaded successfully
type applicationsLoadedMsg struct {
	apps []cli.Application
}

// applicationsErrorMsg is sent when there's an error loading applications
type applicationsErrorMsg struct {
	err error
}

// NewApplicationsModel creates a new applications model
func NewApplicationsModel() ApplicationsModel {
	return ApplicationsModel{
		apps:    []cli.Application{},
		offset:  0,
		limit:   10,
		loading: true,
		cursor:  0,
	}
}

// Init initializes the applications model and loads the first batch
func (m ApplicationsModel) Init() tea.Cmd {
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
		m.loading = false
		m.apps = msg.apps
		m.hasMore = len(msg.apps) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case applicationsErrorMsg:
		m.loading = false
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.apps)-1 {
				m.cursor++
			}

		case "tab":
			return m, func() tea.Msg {
				return ShowMenuMsg{}
			}

		case "shift+left", "shift+h":
			if m.hasPrev && !m.loading {
				m.offset = max(0, m.offset-m.limit)
				m.loading = true
				return m, m.loadApplicationsCmd()
			}

		case "shift+right", "shift+l":
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.loading = true
				return m, m.loadApplicationsCmd()
			}
		}
	}

	return m, nil
}

// View renders the applications listing
func (m ApplicationsModel) View() string {
	if m.loading {
		return "Loading applications..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading applications: %v", m.err))
	}

	var content strings.Builder

	// Applications table header
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-6s %-25s %-15s %-20s", "ID", "Name", "Version", "Status")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", 70))
	content.WriteString("\n")

	// Applications list
	if len(m.apps) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No applications found"))
	} else {
		for i, app := range m.apps {
			var style lipgloss.Style
			if i == m.cursor {
				style = GetSelectedItemStyle()
			} else {
				style = GetUnselectedItemStyle()
			}

			// Determine status
			status := "Disabled"
			if app.Enabled == 1 {
				status = "Enabled"
			}

			// Truncate long names for display
			name := app.Name
			if len(name) > 23 {
				name = name[:20] + "..."
			}

			// Format version
			version := app.Version
			if len(version) > 13 {
				version = version[:10] + "..."
			}
			if version == "" {
				version = "N/A"
			}

			line := fmt.Sprintf("%-6d %-25s %-15s %-20s", app.ID, name, version, status)
			content.WriteString(style.Render(line))
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")

	// Pagination controls
	pageNum := (m.offset / m.limit) + 1

	var prevText, nextText string
	if m.hasPrev {
		prevText = GetSelectedItemStyle().Render("[←] Prev")
	} else {
		prevText = GetItemDescriptionStyle().Render("[←] Prev")
	}

	if m.hasMore {
		nextText = GetSelectedItemStyle().Render("[→] Next")
	} else {
		nextText = GetItemDescriptionStyle().Render("[→] Next")
	}

	pageInfo := GetItemDescriptionStyle().Render(fmt.Sprintf("Page %d", pageNum))
	content.WriteString(prevText + "  " + nextText + "    " + pageInfo)
	content.WriteString("\n")

	return content.String()
}

// loadApplicationsCmd returns a command that loads applications
func (m ApplicationsModel) loadApplicationsCmd() tea.Cmd {
	return func() tea.Msg {
		apps, err := cli.GetApplications(m.limit, m.offset)
		if err != nil {
			return applicationsErrorMsg{err: err}
		}
		return applicationsLoadedMsg{apps: apps}
	}
}
