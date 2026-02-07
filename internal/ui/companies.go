package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CompaniesModel represents the companies listing screen
type CompaniesModel struct {
	companies []cli.Company
	offset    int
	limit     int
	loading   bool
	err       error
	width     int
	height    int
	cursor    int
	hasMore   bool
	hasPrev   bool
}

// companiesLoadedMsg is sent when companies are loaded successfully
type companiesLoadedMsg struct {
	companies []cli.Company
}

// companiesErrorMsg is sent when there's an error loading companies
type companiesErrorMsg struct {
	err error
}

// NewCompaniesModel creates a new companies model
func NewCompaniesModel() CompaniesModel {
	return CompaniesModel{
		companies: []cli.Company{},
		offset:    0,
		limit:     10,
		loading:   true,
		cursor:    0,
	}
}

// Init initializes the companies model and loads the first batch
func (m CompaniesModel) Init() tea.Cmd {
	return m.loadCompaniesCmd()
}

// Update handles messages for the companies model
func (m CompaniesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case companiesLoadedMsg:
		m.loading = false
		m.companies = msg.companies
		m.hasMore = len(msg.companies) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case companiesErrorMsg:
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
			if m.cursor < len(m.companies)-1 {
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
				return m, m.loadCompaniesCmd()
			}

		case "shift+right", "shift+l":
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.loading = true
				return m, m.loadCompaniesCmd()
			}
		}
	}

	return m, nil
}

// View renders the companies listing
func (m CompaniesModel) View() string {
	if m.loading {
		return "Loading companies..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading companies: %v", m.err))
	}

	var content strings.Builder

	// Companies table header - TurboVision style
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-6s %-30s %-15s %-20s", "ID", "Name", "IC", "Status")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("═", 75))
	content.WriteString("\n")

	// Companies list
	if len(m.companies) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No companies found"))
	} else {
		for i, company := range m.companies {
			var style lipgloss.Style
			var prefix string
			if i == m.cursor {
				style = GetSelectedItemStyle()
				prefix = "► " // TurboVision-style focus indicator
			} else {
				style = GetUnselectedItemStyle()
				prefix = "  "
			}

			// Determine status
			status := "Disabled"
			if company.Enabled == 1 {
				status = "Enabled"
			}

			// Truncate long names for display
			name := company.Name
			if len(name) > 28 {
				name = name[:25] + "..."
			}

			line := fmt.Sprintf("%s%-4d %-30s %-15s %-20s", prefix, company.ID, name, company.IC, status)
			content.WriteString(style.Render(line))
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")

	// Pagination controls - TurboVision style
	pageNum := (m.offset / m.limit) + 1

	var prevText, nextText string
	if m.hasPrev {
		prevText = GetSelectedItemStyle().Render("[◄] Prev")
	} else {
		prevText = GetItemDescriptionStyle().Render("[◄] Prev")
	}

	if m.hasMore {
		nextText = GetSelectedItemStyle().Render("[►] Next")
	} else {
		nextText = GetItemDescriptionStyle().Render("[►] Next")
	}

	pageInfo := GetItemDescriptionStyle().Render(fmt.Sprintf("Page %d", pageNum))
	content.WriteString(prevText + "  " + nextText + "    " + pageInfo)
	content.WriteString("\n")

	return content.String()
}

// loadCompaniesCmd returns a command that loads companies
func (m CompaniesModel) loadCompaniesCmd() tea.Cmd {
	return func() tea.Msg {
		companies, err := cli.GetCompanies(m.limit, m.offset)
		if err != nil {
			return companiesErrorMsg{err: err}
		}
		return companiesLoadedMsg{companies: companies}
	}
}
