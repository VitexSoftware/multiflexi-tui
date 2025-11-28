package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TokensModel represents the tokens listing screen
type TokensModel struct {
	tokens  []cli.Token
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

// tokensLoadedMsg is sent when tokens are loaded successfully
type tokensLoadedMsg struct {
	tokens []cli.Token
}

// tokensErrorMsg is sent when there's an error loading tokens
type tokensErrorMsg struct {
	err error
}

// NewTokensModel creates a new tokens model
func NewTokensModel() TokensModel {
	return TokensModel{
		tokens:  []cli.Token{},
		offset:  0,
		limit:   10,
		loading: true,
		cursor:  0,
	}
}

// Init initializes the tokens model and loads the first batch of tokens
func (m TokensModel) Init() tea.Cmd {
	return m.loadTokensCmd()
}

// Update handles messages for the tokens model
func (m TokensModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tokensLoadedMsg:
		m.loading = false
		m.tokens = msg.tokens
		m.hasMore = len(msg.tokens) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case tokensErrorMsg:
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
			if m.cursor < len(m.tokens)-1 {
				m.cursor++
			}

		case "left", "h":
			// Previous page
			if m.hasPrev && !m.loading {
				m.offset = max(0, m.offset-m.limit)
				m.loading = true
				return m, m.loadTokensCmd()
			}

		case "right", "l":
			// Next page
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.loading = true
				return m, m.loadTokensCmd()
			}
		}
	}

	return m, nil
}

// View renders the tokens listing
func (m TokensModel) View() string {
	if m.loading {
		return "Loading tokens..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading tokens: %v", m.err))
	}

	var content strings.Builder

	// Tokens table header
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-8s %-25s %-30s", "ID", "User", "Token")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", 70))
	content.WriteString("\n")

	// Tokens list
	if len(m.tokens) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No tokens found"))
	} else {
		for i, token := range m.tokens {
			var style lipgloss.Style
			if i == m.cursor {
				style = GetSelectedItemStyle()
			} else {
				style = GetUnselectedItemStyle()
			}

			line := fmt.Sprintf("%-8d %-25s %-30s", token.ID, token.User, token.Token)
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

// loadTokensCmd returns a command that loads tokens
func (m TokensModel) loadTokensCmd() tea.Cmd {
	return func() tea.Msg {
		tokens, err := cli.GetTokens(m.limit, m.offset)
		if err != nil {
			return tokensErrorMsg{err: err}
		}
		return tokensLoadedMsg{tokens: tokens}
	}
}
