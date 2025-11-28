package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// PruneModel represents the prune view
type PruneModel struct {
	logs    bool
	jobs    bool
	keep    textinput.Model
	err     error
	pruning bool
	cursor  int
}

// pruneFinishedMsg is sent when pruning is finished successfully
type pruneFinishedMsg struct{}

// pruneErrorMsg is sent when there's an error pruning
type pruneErrorMsg struct {
	err error
}

// NewPruneModel creates a new prune model
func NewPruneModel() PruneModel {
	ti := textinput.New()
	ti.SetValue("1000")
	ti.Focus()
	return PruneModel{
		keep: ti,
	}
}

// Init initializes the prune model
func (m PruneModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the prune model
func (m PruneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pruneFinishedMsg:
		m.pruning = false
		return m, nil

	case pruneErrorMsg:
		m.pruning = false
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 2 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case 0:
				m.logs = !m.logs
			case 1:
				m.jobs = !m.jobs
			case 2:
				if !m.pruning {
					m.pruning = true
					keep, _ := strconv.Atoi(m.keep.Value())
					return m, m.pruneCmd(m.logs, m.jobs, keep)
				}
			}
		}
	}
	m.keep, cmd = m.keep.Update(msg)
	return m, cmd
}

// View renders the prune view
func (m PruneModel) View() string {
	var content strings.Builder

	content.WriteString(GetTitleStyle().Render("Prune Logs and Jobs"))
	content.WriteString("\n\n")

	if m.err != nil {
		content.WriteString(GetErrorStyle().Render(fmt.Sprintf("Error: %v", m.err)))
		content.WriteString("\n\n")
	}

	// Checkboxes
	logsCheckbox := "[ ]"
	if m.logs {
		logsCheckbox = "[x]"
	}
	jobsCheckbox := "[ ]"
	if m.jobs {
		jobsCheckbox = "[x]"
	}

	// Render options
	options := []string{
		fmt.Sprintf("%s Prune logs", logsCheckbox),
		fmt.Sprintf("%s Prune jobs", jobsCheckbox),
		fmt.Sprintf("Keep: %s", m.keep.View()),
	}

	for i, option := range options {
		if i == m.cursor {
			content.WriteString(GetSelectedItemStyle().Render(option))
		} else {
			content.WriteString(option)
		}
		content.WriteString("\n")
	}

	content.WriteString("\n")

	if m.pruning {
		content.WriteString("Pruning...")
	} else {
		content.WriteString(GetButtonStyle().Render("Prune"))
	}

	return content.String()
}

// pruneCmd returns a command that prunes logs and jobs
func (m PruneModel) pruneCmd(logs, jobs bool, keep int) tea.Cmd {
	return func() tea.Msg {
		err := cli.Prune(logs, jobs, keep)
		if err != nil {
			return pruneErrorMsg{err: err}
		}
		return pruneFinishedMsg{}
	}
}
