package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CredTypesModel represents the credtypes listing screen
type CredTypesModel struct {
	credtypes []cli.CredType
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

// credtypesLoadedMsg is sent when credtypes are loaded successfully
type credtypesLoadedMsg struct {
	credtypes []cli.CredType
}

// credtypesErrorMsg is sent when there's an error loading credtypes
type credtypesErrorMsg struct {
	err error
}

// NewCredTypesModel creates a new credtypes model
func NewCredTypesModel() CredTypesModel {
	return CredTypesModel{
		credtypes: []cli.CredType{},
		offset:    0,
		limit:     10,
		loading:   true,
		cursor:    0,
	}
}

// Init initializes the credtypes model and loads the first batch of credtypes
func (m CredTypesModel) Init() tea.Cmd {
	return m.loadCredTypesCmd()
}

// Update handles messages for the credtypes model
func (m CredTypesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case credtypesLoadedMsg:
		m.loading = false
		m.credtypes = msg.credtypes
		m.hasMore = len(msg.credtypes) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case credtypesErrorMsg:
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
			if m.cursor < len(m.credtypes)-1 {
				m.cursor++
			}

		case "left", "h":
			// Previous page
			if m.hasPrev && !m.loading {
				m.offset = max(0, m.offset-m.limit)
				m.loading = true
				return m, m.loadCredTypesCmd()
			}

		case "right", "l":
			// Next page
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.loading = true
				return m, m.loadCredTypesCmd()
			}
		}
	}

	return m, nil
}

// View renders the credtypes listing
func (m CredTypesModel) View() string {
	if m.loading {
		return "Loading credential types..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading credential types: %v", m.err))
	}

	var content strings.Builder

	// CredTypes table header
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-8s %-40s %-30s", "ID", "UUID", "Name")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", 80))
	content.WriteString("\n")

	// CredTypes list
	if len(m.credtypes) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No credential types found"))
	} else {
		for i, credtype := range m.credtypes {
			var style lipgloss.Style
			if i == m.cursor {
				style = GetSelectedItemStyle()
			} else {
				style = GetUnselectedItemStyle()
			}

			line := fmt.Sprintf("%-8d %-40s %-30s", credtype.ID, credtype.UUID, credtype.Name)
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

// loadCredTypesCmd returns a command that loads credtypes
func (m CredTypesModel) loadCredTypesCmd() tea.Cmd {
	return func() tea.Msg {
		credtypes, err := cli.GetCredTypes(m.limit, m.offset)
		if err != nil {
			return credtypesErrorMsg{err: err}
		}
		return credtypesLoadedMsg{credtypes: credtypes}
	}
}
