package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CredentialsModel represents the credentials listing screen
type CredentialsModel struct {
	credentials []cli.Credential
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

// credentialsLoadedMsg is sent when credentials are loaded successfully
type credentialsLoadedMsg struct {
	credentials []cli.Credential
}

// credentialsErrorMsg is sent when there's an error loading credentials
type credentialsErrorMsg struct {
	err error
}

// NewCredentialsModel creates a new credentials model
func NewCredentialsModel() CredentialsModel {
	return CredentialsModel{
		credentials: []cli.Credential{},
		offset:      0,
		limit:       10,
		loading:     true,
		cursor:      0,
	}
}

// Init initializes the credentials model and loads the first batch of credentials
func (m CredentialsModel) Init() tea.Cmd {
	return m.loadCredentialsCmd()
}

// Update handles messages for the credentials model
func (m CredentialsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case credentialsLoadedMsg:
		m.loading = false
		m.credentials = msg.credentials
		m.hasMore = len(msg.credentials) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case credentialsErrorMsg:
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
			if m.cursor < len(m.credentials)-1 {
				m.cursor++
			}

		case "left", "h":
			// Previous page
			if m.hasPrev && !m.loading {
				m.offset = max(0, m.offset-m.limit)
				m.loading = true
				return m, m.loadCredentialsCmd()
			}

		case "right", "l":
			// Next page
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.loading = true
				return m, m.loadCredentialsCmd()
			}
		}
	}

	return m, nil
}

// View renders the credentials listing
func (m CredentialsModel) View() string {
	if m.loading {
		return "Loading credentials..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading credentials: %v", m.err))
	}

	var content strings.Builder

	// Credentials table header
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-8s %-25s %-15s %-20s", "ID", "Name", "Company ID", "Credential Type ID")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", 70))
	content.WriteString("\n")

	// Credentials list
	if len(m.credentials) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No credentials found"))
	} else {
		for i, credential := range m.credentials {
			var style lipgloss.Style
			if i == m.cursor {
				style = GetSelectedItemStyle()
			} else {
				style = GetUnselectedItemStyle()
			}

			line := fmt.Sprintf("%-8d %-25s %-15d %-20d", credential.ID, credential.Name, credential.CompanyID, credential.CredentialTypeID)
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

// loadCredentialsCmd returns a command that loads credentials
func (m CredentialsModel) loadCredentialsCmd() tea.Cmd {
	return func() tea.Msg {
		credentials, err := cli.GetCredentials(m.limit, m.offset)
		if err != nil {
			return credentialsErrorMsg{err: err}
		}
		return credentialsLoadedMsg{credentials: credentials}
	}
}
