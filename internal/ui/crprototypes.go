package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CrPrototypesModel represents the crprototypes listing screen with enhanced data handling
type CrPrototypesModel struct {
	crprototypes    []cli.CrPrototype
	listingManager  *ListingManager
	config          ListingConfig
	state           ListingState
	offset          int
	limit           int
	loading         bool
	err             error
	width           int
	height          int
	cursor          int
	hasMore         bool
	hasPrev         bool
	lastRefresh     time.Time
	autoRefreshTick tea.Cmd
}

// crprototypesLoadedMsg is sent when crprototypes are loaded successfully
type crprototypesLoadedMsg struct {
	crprototypes []cli.CrPrototype
}

// crprototypesErrorMsg is sent when there's an error loading crprototypes
type crprototypesErrorMsg struct {
	err error
}

// NewCrPrototypesModel creates a new crprototypes model with enhanced data handling
func NewCrPrototypesModel() CrPrototypesModel {
	config := DefaultConfigs["crprototype"]
	return CrPrototypesModel{
		crprototypes:   []cli.CrPrototype{},
		listingManager: NewListingManager(),
		config:         config,
		offset:         0,
		limit:          config.DefaultLimit,
		loading:        true,
		cursor:         0,
		lastRefresh:    time.Now(),
	}
}

// Init initializes the crprototypes model and loads the first batch with enhanced data handling
func (m CrPrototypesModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadCrPrototypesCmd(),
		m.listingManager.SetupAutoRefresh(m.config, m.limit, m.offset),
	)
}

// Update handles messages for the crprototypes model
func (m CrPrototypesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case DataLoadedMsg:
		if msg.EntityType == "crprototype" {
			if crprototypes, ok := msg.Data.([]cli.CrPrototype); ok {
				m.loading = false
				m.crprototypes = crprototypes
				m.state = msg.State
				m.hasMore = msg.State.HasMore
				m.hasPrev = msg.State.HasPrev
				m.lastRefresh = time.Now()
				m.err = nil
				return m, nil
			}
		}

	case DataErrorMsg:
		if msg.EntityType == "crprototype" {
			m.loading = false
			m.err = msg.Error
			return m, nil
		}

	case RefreshDataMsg:
		if msg.EntityType == "crprototype" {
			m.loading = true
			return m, m.listingManager.RefreshData(m.config, m.limit, m.offset, msg.Force)
		}

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
				return m, m.listingManager.LoadData(m.config, m.limit, m.offset)
			}

		case "right", "l":
			// Next page
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.cursor = 0
				m.loading = true
				return m, m.listingManager.LoadData(m.config, m.limit, m.offset)
			}

		case "r":
			// Manual refresh with force flag
			if !m.loading {
				m.loading = true
				m.cursor = 0
				return m, m.listingManager.RefreshData(m.config, m.limit, m.offset, true)
			}

		case "c":
			// Clear cache
			m.listingManager.ClearCache()
			return m, nil
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
	refreshText := GetItemDescriptionStyle().Render("[r] Refresh")
	clearCacheText := GetItemDescriptionStyle().Render("[c] Clear Cache")

	content.WriteString(prevText + "  " + nextText + "    " + pageInfo + "    " + refreshText + "  " + clearCacheText)
	content.WriteString("\n")

	// Show cache and performance information
	if m.state.Cached {
		cacheInfo := GetItemDescriptionStyle().Render(fmt.Sprintf("✓ Cached • Fetch: %v • Updated: %s",
			m.state.FetchTime, m.lastRefresh.Format("15:04:05")))
		content.WriteString(cacheInfo)
		content.WriteString("\n")
	} else if m.state.FetchTime > 0 {
		perfInfo := GetItemDescriptionStyle().Render(fmt.Sprintf("Fetch: %v • Updated: %s",
			m.state.FetchTime, m.lastRefresh.Format("15:04:05")))
		content.WriteString(perfInfo)
		content.WriteString("\n")
	}

	return content.String()
}

// loadCrPrototypesCmd returns a command that loads crprototypes using enhanced data handling
func (m CrPrototypesModel) loadCrPrototypesCmd() tea.Cmd {
	return m.listingManager.LoadData(m.config, m.limit, m.offset)
}
