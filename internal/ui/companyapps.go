package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CompanyAppsModel represents the companyapps listing screen
type CompanyAppsModel struct {
	companyapps []cli.CompanyApp
	offset      int
	limit       int
	loading     bool
	err         error
	width       int
	height      int
	cursor      int
	hasMore     bool
	hasPrev     bool
}

// companyappsLoadedMsg is sent when companyapps are loaded successfully
type companyappsLoadedMsg struct {
	companyapps []cli.CompanyApp
}

// companyappsErrorMsg is sent when there's an error loading companyapps
type companyappsErrorMsg struct {
	err error
}

// NewCompanyAppsModel creates a new companyapps model
func NewCompanyAppsModel() CompanyAppsModel {
	return CompanyAppsModel{
		companyapps: []cli.CompanyApp{},
		offset:      0,
		limit:       10,
		loading:     true,
		cursor:      0,
	}
}

// Init initializes the companyapps model and loads the first batch of companyapps
func (m CompanyAppsModel) Init() tea.Cmd {
	return m.loadCompanyAppsCmd()
}

// Update handles messages for the companyapps model
func (m CompanyAppsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case companyappsLoadedMsg:
		m.loading = false
		m.companyapps = msg.companyapps
		m.hasMore = len(msg.companyapps) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case companyappsErrorMsg:
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
			if m.cursor < len(m.companyapps)-1 {
				m.cursor++
			}

		case "left", "h":
			// Previous page
			if m.hasPrev && !m.loading {
				m.offset = max(0, m.offset-m.limit)
				m.loading = true
				return m, m.loadCompanyAppsCmd()
			}

		case "right", "l":
			// Next page
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.loading = true
				return m, m.loadCompanyAppsCmd()
			}
		}
	}

	return m, nil
}

// View renders the companyapps listing
func (m CompanyAppsModel) View() string {
	if m.loading {
		return "Loading company-application relations..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading company-application relations: %v", m.err))
	}

	var content strings.Builder

	// CompanyApps table header
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-8s %-15s %-15s", "ID", "Company ID", "App ID")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", 40))
	content.WriteString("\n")

	// CompanyApps list
	if len(m.companyapps) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No company-application relations found"))
	} else {
		for i, companyapp := range m.companyapps {
			var style lipgloss.Style
			if i == m.cursor {
				style = GetSelectedItemStyle()
			} else {
				style = GetUnselectedItemStyle()
			}

			line := fmt.Sprintf("%-8d %-15d %-15d", companyapp.ID, companyapp.CompanyID, companyapp.AppID)
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

// loadCompanyAppsCmd returns a command that loads companyapps
func (m CompanyAppsModel) loadCompanyAppsCmd() tea.Cmd {
	return func() tea.Msg {
		companyapps, err := cli.GetCompanyApps(m.limit, m.offset)
		if err != nil {
			return companyappsErrorMsg{err: err}
		}
		return companyappsLoadedMsg{companyapps: companyapps}
	}
}
