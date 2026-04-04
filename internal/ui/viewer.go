package ui

import (
	"strings"
	tea "github.com/charmbracelet/bubbletea"
)

// Viewer displays scrollable text content (help, docs).
type Viewer struct {
	title   string
	content string
	err     error
	scroll  int
}

func NewViewer(title string) *Viewer {
	return &Viewer{title: title}
}

func (m *Viewer) SetContent(title, content string) { m.title = title; m.content = content }
func (m *Viewer) SetError(err error)                { m.err = err }

func (m *Viewer) Init() tea.Cmd { return nil }

func (m *Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case HelpLoadedMsg:
		m.title = msg.Command
		m.content = msg.Content
		return m, nil
	case HelpErrorMsg:
		m.err = msg.Err
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.scroll > 0 { m.scroll-- }
		case "down", "j":
			m.scroll++
		case "esc", "q":
			return m, func() tea.Msg { return NavigateBackMsg{} }
		}
	}
	return m, nil
}

func (m *Viewer) View() string {
	var b strings.Builder
	b.WriteString(TitleStyle().Render(m.title))
	b.WriteString("\n\n")
	if m.err != nil {
		b.WriteString(ErrorStyle().Render(m.err.Error()))
		return b.String()
	}
	if m.content == "" {
		b.WriteString(DescriptionStyle().Render("Loading..."))
		return b.String()
	}
	lines := strings.Split(m.content, "\n")
	start := m.scroll
	if start >= len(lines) { start = len(lines) - 1 }
	if start < 0 { start = 0 }
	end := start + 30
	if end > len(lines) { end = len(lines) }
	for _, line := range lines[start:end] {
		b.WriteString(line + "\n")
	}
	b.WriteString("\n" + DescriptionStyle().Render("↑/↓: scroll • esc: back"))
	return b.String()
}
