package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmDialog presents a Y/N confirmation prompt.
type ConfirmDialog struct {
	label  string
	action func() tea.Msg
}

// NewConfirmDialog creates a confirmation dialog.
func NewConfirmDialog(label string, action func() tea.Msg) *ConfirmDialog {
	return &ConfirmDialog{label: label, action: action}
}

func (m *ConfirmDialog) Init() tea.Cmd { return nil }

func (m *ConfirmDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			return m, func() tea.Msg { return ConfirmYesMsg{Action: m.action} }
		case "n", "N", "esc":
			return m, func() tea.Msg { return ConfirmNoMsg{} }
		}
	}
	return m, nil
}

func (m *ConfirmDialog) View() string {
	var b strings.Builder
	b.WriteString(TitleStyle().Render("⚠️  Confirm"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("%s\n\n", m.label))
	b.WriteString(SelectedStyle().Render("[Y] Yes") + "   " + UnselectedStyle().Render("[N] No"))
	b.WriteString("\n")
	return b.String()
}
