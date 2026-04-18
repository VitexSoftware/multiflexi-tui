package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// viewerOverhead: title(1) + help-bar(1) = 2 fixed lines.
const viewerOverhead = 2

// Viewer displays scrollable text content (help, job output, action results).
type Viewer struct {
	title         string
	content       string
	err           error
	scroll        int
	height        int  // available content-area height (set via WindowSizeMsg)
	RefreshOnBack bool // if true, Esc/q returns RefreshCurrentMsg instead of NavigateBackMsg
}

func NewViewer(title string) *Viewer {
	return &Viewer{title: title, height: 30}
}

func (m *Viewer) SetContent(title, content string) {
	m.title = title
	m.content = content
	m.scroll = 0
}

func (m *Viewer) SetError(err error) { m.err = err }

func (m *Viewer) Init() tea.Cmd { return nil }

func (m *Viewer) visibleLines() int {
	v := m.height - viewerOverhead
	if v < 1 {
		v = 1
	}
	return v
}

func (m *Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		return m, nil

	case HelpLoadedMsg:
		m.title = msg.Command
		m.content = msg.Content
		m.scroll = 0
		return m, nil

	case HelpErrorMsg:
		m.err = msg.Err
		return m, nil

	case tea.KeyMsg:
		lines := strings.Split(m.content, "\n")
		vis := m.visibleLines()
		maxScroll := len(lines) - vis
		if maxScroll < 0 {
			maxScroll = 0
		}
		switch msg.String() {
		case "up", "k":
			if m.scroll > 0 {
				m.scroll--
			}
		case "down", "j":
			if m.scroll < maxScroll {
				m.scroll++
			}
		case "pgup":
			m.scroll -= vis
			if m.scroll < 0 {
				m.scroll = 0
			}
		case "pgdown":
			m.scroll += vis
			if m.scroll > maxScroll {
				m.scroll = maxScroll
			}
		case "home", "g":
			m.scroll = 0
		case "end", "G":
			m.scroll = maxScroll
		case "esc", "q":
			if m.RefreshOnBack {
				return m, func() tea.Msg { return NavigateBackAndRefreshMsg{} }
			}
			return m, func() tea.Msg { return NavigateBackMsg{} }
		}
	}
	return m, nil
}

func (m *Viewer) View() string {
	var b strings.Builder
	b.WriteString(TitleStyle().Render(m.title))
	b.WriteString("\n")

	if m.err != nil {
		b.WriteString(ErrorStyle().Render(m.err.Error()) + "\n")
		b.WriteString(FooterStyle().Render("esc: back") + "\n")
		return b.String()
	}
	if m.content == "" {
		b.WriteString(DescriptionStyle().Render("Loading...") + "\n")
		return b.String()
	}

	lines := strings.Split(m.content, "\n")
	vis := m.visibleLines()

	start := m.scroll
	if start >= len(lines) {
		start = len(lines) - 1
	}
	if start < 0 {
		start = 0
	}
	end := start + vis
	if end > len(lines) {
		end = len(lines)
	}

	for _, line := range lines[start:end] {
		b.WriteString(line + "\n")
	}

	// Help/scroll indicator
	if len(lines) > vis {
		pct := 0
		if len(lines)-vis > 0 {
			pct = (m.scroll * 100) / (len(lines) - vis)
		}
		b.WriteString(FooterStyle().Render(
			fmt.Sprintf(" ↑/↓/PgUp/PgDn: scroll  g/G: top/end  esc: back  [%3d%%]", pct)))
	} else {
		b.WriteString(FooterStyle().Render(" ↑/↓: scroll  esc: back"))
	}
	b.WriteString("\n")

	return b.String()
}
