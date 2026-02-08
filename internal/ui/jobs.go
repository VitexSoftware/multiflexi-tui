package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// JobsModel represents the jobs listing screen
type JobsModel struct {
	jobs    []cli.Job
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

// jobsLoadedMsg is sent when jobs are loaded successfully
type jobsLoadedMsg struct {
	jobs []cli.Job
}

// jobsErrorMsg is sent when there's an error loading jobs
type jobsErrorMsg struct {
	err error
}

// OpenJobDetailMsg is sent when a job detail view should be opened
type OpenJobDetailMsg struct {
	Job cli.Job
}

// OpenJobEditorMsg is sent when a job editor should be opened
type OpenJobEditorMsg struct {
	Job cli.Job
}

// ShowMenuMsg is a message to show the menu
type ShowMenuMsg struct{}

// NewJobsModel creates a new jobs model
func NewJobsModel() JobsModel {
	return JobsModel{
		jobs:    []cli.Job{},
		offset:  0,
		limit:   10,
		loading: true,
		cursor:  0,
	}
}

// Init initializes the jobs model and loads the first batch of jobs
func (m JobsModel) Init() tea.Cmd {
	return m.loadJobsCmd()
}

// Update handles messages for the jobs model
func (m JobsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case jobsLoadedMsg:
		m.loading = false
		m.jobs = msg.jobs
		m.hasMore = len(msg.jobs) == m.limit
		m.hasPrev = m.offset > 0
		return m, nil

	case jobsErrorMsg:
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
			if m.cursor < len(m.jobs)-1 {
				m.cursor++
			}

		case "left", "h":
			// Previous page (for job navigation specifically)
			if m.hasPrev && !m.loading {
				m.offset = max(0, m.offset-m.limit)
				m.loading = true
				return m, m.loadJobsCmd()
			}

		case "right", "l":
			// Next page (for job navigation specifically)
			if m.hasMore && !m.loading {
				m.offset += m.limit
				m.loading = true
				return m, m.loadJobsCmd()
			}

		case " ":
			// Space: Open detail view for selected job
			if len(m.jobs) > 0 && m.cursor >= 0 && m.cursor < len(m.jobs) {
				selectedJob := m.jobs[m.cursor]
				return m, func() tea.Msg {
					return OpenJobDetailMsg{Job: selectedJob}
				}
			}

		case "enter", "return":
			// Enter: Open detail view for selected job
			if len(m.jobs) > 0 && m.cursor >= 0 && m.cursor < len(m.jobs) {
				selectedJob := m.jobs[m.cursor]
				return m, func() tea.Msg {
					return OpenJobDetailMsg{Job: selectedJob}
				}
			}

		case "e":
			// E: Open editor for selected job
			if len(m.jobs) > 0 && m.cursor >= 0 && m.cursor < len(m.jobs) {
				selectedJob := m.jobs[m.cursor]
				return m, func() tea.Msg {
					return OpenJobEditorMsg{Job: selectedJob}
				}
			}
		}
	}

	return m, nil
}

// View renders the jobs listing
func (m JobsModel) View() string {
	if m.loading {
		return "Loading jobs..."
	}

	if m.err != nil {
		return GetErrorStyle().Render(fmt.Sprintf("Error loading jobs: %v", m.err))
	}

	var content strings.Builder

	// Jobs table header - TurboVision style
	headerStyle := GetSelectedItemStyle().Copy().Bold(true)
	content.WriteString(headerStyle.Render(fmt.Sprintf("%-8s %-25s %-15s %-20s", "ID", "Command", "Status", "Schedule")))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("═", 70))
	content.WriteString("\n")

	// Jobs list
	if len(m.jobs) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No jobs found"))
	} else {
		for i, job := range m.jobs {
			var style lipgloss.Style
			var prefix string
			if i == m.cursor {
				style = GetSelectedItemStyle()
				prefix = "► " // TurboVision-style focus indicator
			} else {
				style = GetUnselectedItemStyle()
				prefix = "  "
			}

			// Determine status based on exitcode and PID
			status := "Running"
			if job.PID == 0 {
				if job.Exitcode == -1 {
					status = "Scheduled"
				} else if job.Exitcode == 0 {
					status = "Success"
				} else {
					status = "Failed"
				}
			}

			// Truncate long command names for display
			command := job.Command
			if len(command) > 23 {
				command = command[:20] + "..."
			}

			// Format schedule time
			schedule := job.Schedule
			if len(schedule) > 18 {
				// Extract just the time part if it's a full datetime
				if len(schedule) >= 16 {
					schedule = schedule[11:16] // Extract HH:MM from YYYY-MM-DD HH:MM:SS
				}
			}

			line := fmt.Sprintf("%s%-6d %-25s %-15s %-20s", prefix, job.ID, command, status, schedule)
			content.WriteString(style.Render(line))
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")

	// Pagination controls - TurboVision style
	pageNum := (m.offset / m.limit) + 1

	var prevText, nextText string
	if m.hasPrev {
		prevText = GetSelectedItemStyle().Render("[◄] Prev")
	} else {
		prevText = GetItemDescriptionStyle().Render("[◄] Prev")
	}

	if m.hasMore {
		nextText = GetSelectedItemStyle().Render("[►] Next")
	} else {
		nextText = GetItemDescriptionStyle().Render("[►] Next")
	}

	pageInfo := GetItemDescriptionStyle().Render(fmt.Sprintf("Page %d", pageNum))

	// Simple concatenation to ensure single line
	content.WriteString(prevText + "  " + nextText + "    " + pageInfo)
	content.WriteString("\n")

	return content.String()
}

// loadJobsCmd returns a command that loads jobs
func (m JobsModel) loadJobsCmd() tea.Cmd {
	return func() tea.Msg {
		jobs, err := cli.GetJobs(m.limit, m.offset)
		if err != nil {
			return jobsErrorMsg{err: err}
		}
		return jobsLoadedMsg{jobs: jobs}
	}
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
