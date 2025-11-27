package ui

import (
	"fmt"
	"strings"
)

// TableColumn represents a column in the table
type TableColumn struct {
	Header string
	Width  int
	Field  string // Field name or identifier
}

// TableRow represents a single row of data
type TableRow struct {
	ID     int                    // Unique identifier for the row
	Values map[string]interface{} // Column values keyed by Field name
}

// TableWidget is a reusable table component for displaying listings
type TableWidget struct {
	title      string
	columns    []TableColumn
	rows       []TableRow
	cursor     int
	offset     int
	limit      int
	loading    bool
	err        error
	hasMore    bool
	pageNum    int
	totalShown int
	helpText   string
}

// TableConfig contains configuration for creating a table
type TableConfig struct {
	Title    string
	Columns  []TableColumn
	Limit    int
	HelpText string
}

// NewTableWidget creates a new table widget
func NewTableWidget(config TableConfig) *TableWidget {
	return &TableWidget{
		title:    config.Title,
		columns:  config.Columns,
		rows:     []TableRow{},
		cursor:   0,
		offset:   0,
		limit:    config.Limit,
		loading:  true,
		helpText: config.HelpText,
	}
}

// SetData updates the table with new data
func (t *TableWidget) SetData(rows []TableRow) {
	t.rows = rows
	t.loading = false
	t.err = nil
	t.hasMore = len(rows) >= t.limit
	t.pageNum = (t.offset / t.limit) + 1
	t.totalShown = len(rows)

	// Reset cursor if it's beyond the data
	if t.cursor >= len(rows) && len(rows) > 0 {
		t.cursor = len(rows) - 1
	}
	if len(rows) == 0 {
		t.cursor = 0
	}
}

// SetLoading sets the loading state
func (t *TableWidget) SetLoading(loading bool) {
	t.loading = loading
}

// SetError sets the error state
func (t *TableWidget) SetError(err error) {
	t.err = err
	t.loading = false
}

// GetCursor returns the current cursor position
func (t *TableWidget) GetCursor() int {
	return t.cursor
}

// GetSelectedRow returns the currently selected row, or nil if none
func (t *TableWidget) GetSelectedRow() *TableRow {
	if len(t.rows) == 0 || t.cursor >= len(t.rows) || t.cursor < 0 {
		return nil
	}
	return &t.rows[t.cursor]
}

// GetOffset returns the current offset
func (t *TableWidget) GetOffset() int {
	return t.offset
}

// GetLimit returns the current limit
func (t *TableWidget) GetLimit() int {
	return t.limit
}

// HasMore returns whether there are more pages
func (t *TableWidget) HasMore() bool {
	return t.hasMore
}

// HandleKeypress processes keyboard input for navigation
func (t *TableWidget) HandleKeypress(key string) (needsRefresh bool, needsNextPage bool, needsPrevPage bool, openDetail bool) {
	switch key {
	case "up", "k":
		if t.cursor > 0 {
			t.cursor--
		}
	case "down", "j":
		if t.cursor < len(t.rows)-1 {
			t.cursor++
		}
	case "enter":
		// Open detail view for selected row
		if len(t.rows) > 0 && t.cursor >= 0 && t.cursor < len(t.rows) {
			return false, false, false, true
		}
	case "n", "right", "pgdown":
		// Next page
		if t.hasMore {
			t.offset += t.limit
			t.cursor = 0 // Reset cursor for new page
			return false, true, false, false
		}
	case "p", "left", "pgup":
		// Previous page
		if t.offset > 0 {
			t.offset -= t.limit
			if t.offset < 0 {
				t.offset = 0
			}
			t.cursor = 0 // Reset cursor for new page
			return false, false, true, false
		}
	case "r":
		// Refresh
		t.cursor = 0
		return true, false, false, false
	}
	return false, false, false, false
}

// View renders the table widget
func (t *TableWidget) View() string {
	var content strings.Builder

	// Title
	if t.title != "" {
		content.WriteString(GetTitleStyle().Render(t.title))
		content.WriteString("\n\n")
	}

	if t.loading {
		content.WriteString(GetItemDescriptionStyle().Render("Loading..."))
		content.WriteString("\n")
		return content.String()
	}

	if t.err != nil {
		content.WriteString(GetItemDescriptionStyle().Render(fmt.Sprintf("Error: %v", t.err)))
		content.WriteString("\n")
		return content.String()
	}

	if len(t.rows) == 0 {
		content.WriteString(GetItemDescriptionStyle().Render("No data found"))
		content.WriteString("\n")
		return content.String()
	}

	// Calculate total width for separator
	totalWidth := 0
	for _, col := range t.columns {
		totalWidth += col.Width + 1 // +1 for spacing
	}

	// Table header - plain text
	headerParts := make([]string, len(t.columns))
	for i, col := range t.columns {
		format := fmt.Sprintf("%%-%ds", col.Width)
		headerParts[i] = fmt.Sprintf(format, col.Header)
	}
	content.WriteString(strings.Join(headerParts, " "))
	content.WriteString("\n")

	// Separator
	content.WriteString(strings.Repeat("─", totalWidth))
	content.WriteString("\n")

	// Table rows - plain text with focus indicator
	for i, row := range t.rows {
		rowParts := make([]string, len(t.columns))

		for j, col := range t.columns {
			value := ""
			if val, exists := row.Values[col.Field]; exists {
				value = fmt.Sprintf("%v", val)
			}

			// Truncate if too long
			if len(value) > col.Width {
				if col.Width > 3 {
					value = value[:col.Width-3] + "..."
				} else {
					value = value[:col.Width]
				}
			}

			format := fmt.Sprintf("%%-%ds", col.Width)
			rowParts[j] = fmt.Sprintf(format, value)
		}

		// Add focus indicator
		focusIndicator := " "
		if i == t.cursor {
			focusIndicator = "→"
		}

		line := focusIndicator + strings.Join(rowParts, " ")
		content.WriteString(line)
		content.WriteString("\n")
	}

	content.WriteString("\n")

	// Pagination info
	content.WriteString(fmt.Sprintf("Page %d • Showing %d items • Offset %d",
		t.pageNum, t.totalShown, t.offset))
	content.WriteString("\n\n")

	// Pagination buttons
	var prevButton, nextButton string
	if t.offset > 0 {
		prevButton = "← Previous"
	} else {
		prevButton = "  Previous"
	}

	if t.hasMore {
		nextButton = "Next →"
	} else {
		nextButton = "Next  "
	}

	content.WriteString(fmt.Sprintf("[%s]   [%s]   [r] Refresh", prevButton, nextButton))
	content.WriteString("\n")
	content.WriteString("Enter: View details • PgUp/PgDn: Page navigation • ←/→: Page navigation • ↑/↓: Row selection")
	content.WriteString("\n")

	// Help text
	if t.helpText != "" {
		content.WriteString("\n")
		content.WriteString(GetItemDescriptionStyle().Render(t.helpText))
		content.WriteString("\n")
	}

	return content.String()
}
