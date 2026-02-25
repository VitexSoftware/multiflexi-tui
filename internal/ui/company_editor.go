package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// CompanyEditorModel represents the editor for a company
type CompanyEditorModel struct {
	company    cli.Company
	cursor     int
	inputs     []textinput.Model
	labels     []string
	width      int
	height     int
}

// SaveCompanyMsg is sent when the company should be saved
type SaveCompanyMsg struct {
	Company cli.Company
}

// NewCompanyEditorModel creates a new company editor model
func NewCompanyEditorModel(company cli.Company) CompanyEditorModel {
	nameInput := textinput.New()
	nameInput.Placeholder = "Company Name"
	nameInput.SetValue(company.Name)
	nameInput.Focus()

	emailInput := textinput.New()
	emailInput.Placeholder = "Email"
	emailInput.SetValue(company.Email)

	icInput := textinput.New()
	icInput.Placeholder = "IC"
	icInput.SetValue(company.IC)

	slugInput := textinput.New()
	slugInput.Placeholder = "Slug"
	slugInput.SetValue(company.Slug)

	inputs := []textinput.Model{nameInput, emailInput, icInput, slugInput}
	labels := []string{"Name", "Email", "IC", "Slug"}

	return CompanyEditorModel{
		company: company,
		inputs:  inputs,
		labels:  labels,
		cursor:  0,
	}
}

// Init initializes the company editor model
func (m CompanyEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the company editor model
func (m CompanyEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return BackMsg{} }
		case "tab", "down":
			m.inputs[m.cursor].Blur()
			m.cursor = (m.cursor + 1) % len(m.inputs)
			m.inputs[m.cursor].Focus()
			return m, textinput.Blink
		case "shift+tab", "up":
			m.inputs[m.cursor].Blur()
			m.cursor = (m.cursor - 1 + len(m.inputs)) % len(m.inputs)
			m.inputs[m.cursor].Focus()
			return m, textinput.Blink
		case "enter":
			m.company.Name = m.inputs[0].Value()
			m.company.Email = m.inputs[1].Value()
			m.company.IC = m.inputs[2].Value()
			m.company.Slug = m.inputs[3].Value()
			return m, func() tea.Msg {
				return SaveCompanyMsg{Company: m.company}
			}
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	m.inputs[m.cursor], cmd = m.inputs[m.cursor].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the company editor model
func (m CompanyEditorModel) View() string {
	var content strings.Builder

	content.WriteString(GetTitleStyle().Render(fmt.Sprintf("Editing Company: %s", m.company.Name)))
	content.WriteString("\n\n")

	for i, input := range m.inputs {
		label := m.labels[i]
		if i == m.cursor {
			content.WriteString(GetSelectedItemStyle().Render(fmt.Sprintf("%-15s", label+":")))
		} else {
			content.WriteString(fmt.Sprintf("%-15s", label+":"))
		}
		content.WriteString(" ")
		content.WriteString(input.View())
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(GetFooterStyle().Render("tab/↑↓: navigate fields • enter: save • esc: cancel"))
	content.WriteString("\n")

	return content.String()
}
