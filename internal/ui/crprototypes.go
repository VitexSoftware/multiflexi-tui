package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CrPrototypesModel represents the crprototypes listing screen
type CrPrototypesModel struct {
	crprototypes []cli.CrPrototype
	offset       int
	limit        int
	loading      bool
	err          error
	width        int
	height       int
	cursor       int
	hasMore      bool
	hasPrev      bool
}

// crprototypesLoadedMsg is sent when crprototypes are loaded successfully
type crprototypesLoadedMsg struct {
	crprototypes []cli.CrPrototype
}

// crprototypesErrorMsg is sent when there's an error loading crprototypes
type crprototypesErrorMsg struct {
	err error
}

// NewCrPrototypesModel creates a new crprototypes model
func NewCrPrototypesModel() CrPrototypesModel {
	return CrPrototypesModel{
		crprototypes: []cli.CrPrototype{},
		offset:       0,
		limit:        10,
		loading:      true,
		cursor:       0,
	}
}

// Init initializes the crprototypes model and loads the first batch of crprototypes
func (m CrPrototypesModel) Init() tea.Cmd {
	return m.loadCrPrototypesCmd()
}

// Update handles messages for the crprototypes model
func (m CrPrototypesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case crprototypesLoadedMsg:
		m.loading = false
		m.crprototypes = msg.crprototypes
		m.hasMore = len(msg.crprototypes) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case crprototypesErrorMsg:
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
			if m.cursor < len(m.crprototypes)-1 {
				m.cursor++
			}

		case "left", "h":
			// Previous page
			if m.hasPrev && !m.loading {
				m.offset -= m.limit
				if m.offset < 0 {
					m.offset = 0
				}
				m.cursor = 0
				m.loading = true
				return m, m.loadCrPrototypesCmd()
			}

		case "right", "l":
			// Next page
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.cursor = 0
				m.loading = true
				return m, m.loadCrPrototypesCmd()
			}

		case "r":
			// Refresh current page
			if !m.loading {
				m.loading = true
				m.cursor = 0
				return m, m.loadCrPrototypesCmd()
			}
		}
	}

	return m, nil
}

// View renders the crprototypes listing
func (m CrPrototypesModel) View() string {
	if m.loading {
		return "Loading credential prototypes..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading credential prototypes: %v", m.err))
	}

	var content strings.Builder

	// CrPrototypes table header
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-8s %-35s %-40s %-10s", "ID", "Name", "Description", "Version")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", 95))
	content.WriteString("\n")

	// CrPrototypes list
	if len(m.crprototypes) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No credential prototypes found"))
	} else {
		for i, crprototype := range m.crprototypes {
			var style lipgloss.Style
			if i == m.cursor {
				style = GetSelectedItemStyle()
			} else {
				style = GetUnselectedItemStyle()
			}

			// Truncate long descriptions for display
			description := crprototype.Description
			if len(description) > 38 {
				description = description[:35] + "..."
			}

			line := fmt.Sprintf("%-8d %-35s %-40s %-10s", crprototype.ID, crprototype.Name, description, crprototype.Version)
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

// loadCrPrototypesCmd returns a command that loads crprototypes
func (m CrPrototypesModel) loadCrPrototypesCmd() tea.Cmd {
	return func() tea.Msg {
		crprototypes, err := cli.GetCrPrototypes(m.limit, m.offset)
		if err != nil {
			return crprototypesErrorMsg{err: err}
		}
		return crprototypesLoadedMsg{crprototypes: crprototypes}
	}
}
