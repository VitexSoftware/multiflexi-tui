package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// QueueModel represents the queue listing screen
type QueueModel struct {
	queue   []cli.Queue
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

// queueLoadedMsg is sent when queue items are loaded successfully
type queueLoadedMsg struct {
	queue []cli.Queue
}

// queueErrorMsg is sent when there's an error loading queue items
type queueErrorMsg struct {
	err error
}

// queueTruncatedMsg is sent when the queue is truncated successfully
type queueTruncatedMsg struct{}

// NewQueueModel creates a new queue model
func NewQueueModel() QueueModel {
	return QueueModel{
		queue:   []cli.Queue{},
		offset:  0,
		limit:   10,
		loading: true,
		cursor:  0,
	}
}

// Init initializes the queue model and loads the first batch of queue items
func (m QueueModel) Init() tea.Cmd {
	return m.loadQueueCmd()
}

// Update handles messages for the queue model
func (m QueueModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case queueLoadedMsg:
		m.loading = false
		m.queue = msg.queue
		m.hasMore = len(msg.queue) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case queueErrorMsg:
		m.loading = false
		m.err = msg.err
		return m, nil

	case queueTruncatedMsg:
		m.loading = true
		return m, m.loadQueueCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.queue)-1 {
				m.cursor++
			}

		case "left", "h":
			// Previous page
			if m.hasPrev && !m.loading {
				m.offset = max(0, m.offset-m.limit)
				m.loading = true
				return m, m.loadQueueCmd()
			}

		case "right", "l":
			// Next page
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.loading = true
				return m, m.loadQueueCmd()
			}

		case "enter":
			// Truncate queue
			return m, m.truncateQueueCmd()
		}
	}

	return m, nil
}

// View renders the queue listing
func (m QueueModel) View() string {
	if m.loading {
		return "Loading queue..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading queue: %v", m.err))
	}

	var content strings.Builder

	// Queue table header
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-8s %-s", "ID", "Message")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", m.width))
	content.WriteString("\n")

	// Queue list
	if len(m.queue) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No queue items found"))
	} else {
		for i, item := range m.queue {
			var style lipgloss.Style
			if i == m.cursor {
				style = GetSelectedItemStyle()
			} else {
				style = GetUnselectedItemStyle()
			}

			line := fmt.Sprintf("%-8d %-s", item.ID, item.Message)
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
	content.WriteString("\n\n")
	content.WriteString(GetButtonStyle().Render("Truncate Queue"))

	return content.String()
}

// loadQueueCmd returns a command that loads queue items
func (m QueueModel) loadQueueCmd() tea.Cmd {
	return func() tea.Msg {
		queue, err := cli.GetQueue(m.limit, m.offset)
		if err != nil {
			return queueErrorMsg{err: err}
		}
		return queueLoadedMsg{queue: queue}
	}
}

// truncateQueueCmd returns a command that truncates the queue
func (m QueueModel) truncateQueueCmd() tea.Cmd {
	return func() tea.Msg {
		err := cli.TruncateQueue()
		if err != nil {
			return queueErrorMsg{err: err}
		}
		return queueTruncatedMsg{}
	}
}
