package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ArtifactsModel represents the artifacts listing screen
type ArtifactsModel struct {
	artifacts []cli.Artifact
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

// artifactsLoadedMsg is sent when artifacts are loaded successfully
type artifactsLoadedMsg struct {
	artifacts []cli.Artifact
}

// artifactsErrorMsg is sent when there's an error loading artifacts
type artifactsErrorMsg struct {
	err error
}

// NewArtifactsModel creates a new artifacts model
func NewArtifactsModel() ArtifactsModel {
	return ArtifactsModel{
		artifacts: []cli.Artifact{},
		offset:    0,
		limit:     10,
		loading:   true,
		cursor:    0,
	}
}

// Init initializes the artifacts model and loads the first batch of artifacts
func (m ArtifactsModel) Init() tea.Cmd {
	return m.loadArtifactsCmd()
}

// Update handles messages for the artifacts model
func (m ArtifactsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case artifactsLoadedMsg:
		m.loading = false
		m.artifacts = msg.artifacts
		m.hasMore = len(msg.artifacts) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case artifactsErrorMsg:
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
			if m.cursor < len(m.artifacts)-1 {
				m.cursor++
			}

		case "left", "h":
			// Previous page
			if m.hasPrev && !m.loading {
				m.offset = max(0, m.offset-m.limit)
				m.loading = true
				return m, m.loadArtifactsCmd()
			}

		case "right", "l":
			// Next page
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.loading = true
				return m, m.loadArtifactsCmd()
			}
		}
	}

	return m, nil
}

// View renders the artifacts listing
func (m ArtifactsModel) View() string {
	if m.loading {
		return "Loading artifacts..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading artifacts: %v", m.err))
	}

	var content strings.Builder

	// Artifacts table header
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-8s %-15s %-50s", "ID", "Job ID", "File")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", 75))
	content.WriteString("\n")

	// Artifacts list
	if len(m.artifacts) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No artifacts found"))
	} else {
		for i, artifact := range m.artifacts {
			var style lipgloss.Style
			if i == m.cursor {
				style = GetSelectedItemStyle()
			} else {
				style = GetUnselectedItemStyle()
			}

			line := fmt.Sprintf("%-8d %-15d %-50s", artifact.ID, artifact.Job_ID, artifact.Filename)
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

// loadArtifactsCmd returns a command that loads artifacts
func (m ArtifactsModel) loadArtifactsCmd() tea.Cmd {
	return func() tea.Msg {
		artifacts, err := cli.GetArtifacts(m.limit, m.offset)
		if err != nil {
			return artifactsErrorMsg{err: err}
		}
		return artifactsLoadedMsg{artifacts: artifacts}
	}
}
