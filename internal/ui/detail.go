package ui

import (
	"fmt"
	"strings"
)

// DetailField represents a single field in the detail view
type DetailField struct {
	Label string
	Value string
	Width int // Optional width for formatting, 0 means auto
}

// DetailAction represents an action button in the detail view
type DetailAction struct {
	Label   string
	Key     string // Key to trigger the action
	Command string // Command identifier for handling
}

// DetailWidget is a reusable detail view component
type DetailWidget struct {
	title   string
	fields  []DetailField
	actions []DetailAction
	loading bool
	err     error
	data    interface{} // Store the original data object
}

// DetailConfig contains configuration for creating a detail view
type DetailConfig struct {
	Title   string
	Actions []DetailAction
}

// NewDetailWidget creates a new detail widget
func NewDetailWidget(config DetailConfig) *DetailWidget {
	return &DetailWidget{
		title:   config.Title,
		actions: config.Actions,
		fields:  []DetailField{},
		loading: true,
	}
}

// SetData updates the detail view with new data
func (d *DetailWidget) SetData(fields []DetailField, data interface{}) {
	d.fields = fields
	d.data = data
	d.loading = false
	d.err = nil
}

// SetLoading sets the loading state
func (d *DetailWidget) SetLoading(loading bool) {
	d.loading = loading
}

// SetError sets the error state
func (d *DetailWidget) SetError(err error) {
	d.err = err
	d.loading = false
}

// GetData returns the stored data object
func (d *DetailWidget) GetData() interface{} {
	return d.data
}

// HandleKeypress processes keyboard input for actions
func (d *DetailWidget) HandleKeypress(key string) (actionCommand string, shouldGoBack bool) {
	switch key {
	case "escape", "q":
		return "", true
	default:
		// Check if key matches any action
		for _, action := range d.actions {
			if action.Key == key {
				return action.Command, false
			}
		}
	}
	return "", false
}

// View renders the detail widget
func (d *DetailWidget) View() string {
	var content strings.Builder

	// Title
	if d.title != "" {
		content.WriteString(GetTitleStyle().Render(d.title))
		content.WriteString("\n\n")
	}

	if d.loading {
		content.WriteString(GetItemDescriptionStyle().Render("Loading details..."))
		content.WriteString("\n")
		return content.String()
	}

	if d.err != nil {
		content.WriteString(GetItemDescriptionStyle().Render(fmt.Sprintf("Error: %v", d.err)))
		content.WriteString("\n")
		return content.String()
	}

	if len(d.fields) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No details available"))
		content.WriteString("\n")
		return content.String()
	}

	// Calculate max label width for alignment
	maxLabelWidth := 0
	for _, field := range d.fields {
		if len(field.Label) > maxLabelWidth {
			maxLabelWidth = len(field.Label)
		}
	}

	// Detail fields
	for _, field := range d.fields {
		label := fmt.Sprintf("%-*s:", maxLabelWidth, field.Label)
		content.WriteString(fmt.Sprintf("%s %s", label, field.Value))
		content.WriteString("\n")
	}

	content.WriteString("\n")

	// Action buttons
	if len(d.actions) > 0 {
		var actionButtons []string
		for _, action := range d.actions {
			actionButtons = append(actionButtons, fmt.Sprintf("[%s] %s", action.Key, action.Label))
		}
		content.WriteString(strings.Join(actionButtons, "   "))
		content.WriteString("\n\n")
	}

	// Navigation help
	content.WriteString("ESC/q: Back to list")
	content.WriteString("\n")

	return content.String()
}
