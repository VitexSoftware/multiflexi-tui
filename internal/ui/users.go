package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UsersModel represents the users listing screen
type UsersModel struct {
	users   []cli.User
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

// usersLoadedMsg is sent when users are loaded successfully
type usersLoadedMsg struct {
	users []cli.User
}

// usersErrorMsg is sent when there's an error loading users
type usersErrorMsg struct {
	err error
}

// NewUsersModel creates a new users model
func NewUsersModel() UsersModel {
	return UsersModel{
		users:   []cli.User{},
		offset:  0,
		limit:   10,
		loading: true,
		cursor:  0,
	}
}

// Init initializes the users model and loads the first batch of users
func (m UsersModel) Init() tea.Cmd {
	return m.loadUsersCmd()
}

// Update handles messages for the users model
func (m UsersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case usersLoadedMsg:
		m.loading = false
		m.users = msg.users
		m.hasMore = len(msg.users) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case usersErrorMsg:
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
			if m.cursor < len(m.users)-1 {
				m.cursor++
			}

		case "left", "h":
			// Previous page
			if m.hasPrev && !m.loading {
				m.offset = max(0, m.offset-m.limit)
				m.loading = true
				return m, m.loadUsersCmd()
			}

		case "right", "l":
			// Next page
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.loading = true
				return m, m.loadUsersCmd()
			}
		}
	}

	return m, nil
}

// View renders the users listing
func (m UsersModel) View() string {
	if m.loading {
		return "Loading users..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading users: %v", m.err))
	}

	var content strings.Builder

	// Users table header
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-8s %-20s %-20s %-20s %-30s", "ID", "Login", "First Name", "Last Name", "Email")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", 100))
	content.WriteString("\n")

	// Users list
	if len(m.users) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No users found"))
	} else {
		for i, user := range m.users {
			var style lipgloss.Style
			if i == m.cursor {
				style = GetSelectedItemStyle()
			} else {
				style = GetUnselectedItemStyle()
			}

			line := fmt.Sprintf("%-8d %-20s %-20s %-20s %-30s", user.ID, user.Login, user.Firstname, user.Lastname, user.Email)
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

// loadUsersCmd returns a command that loads users
func (m UsersModel) loadUsersCmd() tea.Cmd {
	return func() tea.Msg {
		users, err := cli.GetUsers(m.limit, m.offset)
		if err != nil {
			return usersErrorMsg{err: err}
		}
		return usersLoadedMsg{users: users}
	}
}
