package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmDeleteYesMsg is sent when the user confirms deletion
type ConfirmDeleteYesMsg struct {
	Data interface{}
}

// ConfirmDeleteNoMsg is sent when the user cancels deletion
type ConfirmDeleteNoMsg struct{}

// ConfirmDialogModel represents a confirmation dialog
type ConfirmDialogModel struct {
	label  string
	data   interface{}
	width  int
	height int
}

// NewConfirmDialogModel creates a new confirmation dialog
func NewConfirmDialogModel(label string, data interface{}) ConfirmDialogModel {
	return ConfirmDialogModel{
		label: label,
		data:  data,
	}
}

// Init initializes the confirmation dialog
func (m ConfirmDialogModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the confirmation dialog
func (m ConfirmDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			data := m.data
			return m, func() tea.Msg {
				return ConfirmDeleteYesMsg{Data: data}
			}
		case "n", "N", "esc", "q":
			return m, func() tea.Msg {
				return ConfirmDeleteNoMsg{}
			}
		}
	}
	return m, nil
}

// View renders the confirmation dialog
func (m ConfirmDialogModel) View() string {
	var content strings.Builder

	content.WriteString("\n")
	content.WriteString(GetTitleStyle().Render("⚠  Confirm Delete"))
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("  Are you sure you want to delete %s?\n", m.label))
	content.WriteString("\n")
	content.WriteString("  " + GetSelectedItemStyle().Render(" [y] Yes ") + "   " + GetUnselectedItemStyle().Render(" [n] No "))
	content.WriteString("\n\n")
	content.WriteString(GetItemDescriptionStyle().Render("  This action cannot be undone."))
	content.WriteString("\n")

	return content.String()
}
