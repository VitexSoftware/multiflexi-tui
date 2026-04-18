package ui

import (
	"fmt"
	"strings"
)

// tableOverhead is the number of non-data lines rendered by View():
// title(1) + col-header(1) + top-sep(1) + bottom-sep(1) + pagination(1) = 5
const tableOverhead = 5

// minTableRows is the minimum number of data rows the table will display.
const minTableRows = 3

// TableWidget renders a paginated table with cursor selection.
type TableWidget struct {
	title    string
	columns  []TableColumn
	rows     []TableRow
	cursor   int
	offset   int
	limit    int
	loading  bool
	err      error
	hasMore  bool
	pageNum  int
	helpText string
}

// NewTableWidget creates a new table widget.
func NewTableWidget(title string, columns []TableColumn, limit int, helpText string) *TableWidget {
	if limit < minTableRows {
		limit = minTableRows
	}
	return &TableWidget{
		title:    title,
		columns:  columns,
		limit:    limit,
		helpText: helpText,
		loading:  true,
	}
}

// SetContentHeight adjusts the row limit to fill the given content-area height.
// Returns true when the limit changed and the caller should re-fetch.
func (t *TableWidget) SetContentHeight(h int) bool {
	newLimit := h - tableOverhead
	if newLimit < minTableRows {
		newLimit = minTableRows
	}
	if newLimit == t.limit {
		return false
	}
	t.limit = newLimit
	return true
}

// SetData updates the table with fresh data.
func (t *TableWidget) SetData(rows []TableRow) {
	t.rows = rows
	t.loading = false
	t.err = nil
	t.hasMore = len(rows) >= t.limit
	t.pageNum = (t.offset / t.limit) + 1
	if t.cursor >= len(rows) && len(rows) > 0 {
		t.cursor = len(rows) - 1
	}
	if len(rows) == 0 {
		t.cursor = 0
	}
}

func (t *TableWidget) SetLoading(l bool) { t.loading = l }
func (t *TableWidget) SetError(e error)  { t.err = e; t.loading = false }
func (t *TableWidget) Cursor() int       { return t.cursor }
func (t *TableWidget) Offset() int       { return t.offset }
func (t *TableWidget) Limit() int        { return t.limit }

// SelectedRow returns the row at the cursor, or nil.
func (t *TableWidget) SelectedRow() *TableRow {
	if len(t.rows) == 0 || t.cursor < 0 || t.cursor >= len(t.rows) {
		return nil
	}
	return &t.rows[t.cursor]
}

// HandleKey processes navigation keys. Returns action flags.
func (t *TableWidget) HandleKey(key string) (refresh, nextPage, prevPage, openDetail, openEditor, openCreate bool) {
	switch key {
	case "up", "k":
		if t.cursor > 0 {
			t.cursor--
		}
	case "down", "j":
		if t.cursor < len(t.rows)-1 {
			t.cursor++
		}
	case "enter", " ":
		if len(t.rows) > 0 {
			return false, false, false, true, false, false
		}
	case "e":
		if len(t.rows) > 0 {
			return false, false, false, false, true, false
		}
	case "n":
		return false, false, false, false, false, true
	case "right", "pgdown":
		if t.hasMore {
			t.offset += t.limit
			t.cursor = 0
			return false, true, false, false, false, false
		}
	case "left", "pgup":
		if t.offset > 0 {
			t.offset -= t.limit
			if t.offset < 0 {
				t.offset = 0
			}
			t.cursor = 0
			return false, false, true, false, false, false
		}
	case "r":
		t.cursor = 0
		return true, false, false, false, false, false
	}
	return false, false, false, false, false, false
}

// View renders the table filling available height.
func (t *TableWidget) View() string {
	var b strings.Builder

	// Title
	if t.title != "" {
		b.WriteString(TitleStyle().Render(t.title))
		b.WriteString("\n")
	}

	if t.loading {
		b.WriteString(DescriptionStyle().Render("  Loading..."))
		b.WriteString("\n")
		return b.String()
	}
	if t.err != nil {
		b.WriteString(ErrorStyle().Render(fmt.Sprintf("  Error: %v", t.err)))
		b.WriteString("\n")
		return b.String()
	}

	// Compute total column width for separators
	totalWidth := 1 // leading indicator char
	for _, col := range t.columns {
		totalWidth += col.Width + 1
	}
	sep := strings.Repeat("─", totalWidth)

	// Column headers
	parts := make([]string, len(t.columns))
	for i, col := range t.columns {
		parts[i] = fmt.Sprintf("%-*s", col.Width, col.Header)
	}
	b.WriteString(" " + strings.Join(parts, " ") + "\n")
	b.WriteString(sep + "\n")

	// Data rows
	if len(t.rows) == 0 {
		b.WriteString(DescriptionStyle().Render("  (no items)") + "\n")
	} else {
		count := len(t.rows)
		if count > t.limit {
			count = t.limit
		}
		for i := 0; i < count; i++ {
			row := t.rows[i]
			rowParts := make([]string, len(t.columns))
			for j, col := range t.columns {
				val := row.Values[col.Field]
				if len(val) > col.Width {
					if col.Width > 3 {
						val = val[:col.Width-3] + "..."
					} else {
						val = val[:col.Width]
					}
				}
				rowParts[j] = fmt.Sprintf("%-*s", col.Width, val)
			}
			indicator := " "
			if i == t.cursor {
				indicator = "►"
			}
			line := indicator + strings.Join(rowParts, " ")
			if i == t.cursor {
				b.WriteString(SelectedStyle().Render(line))
			} else {
				b.WriteString(UnselectedStyle().Render(line))
			}
			b.WriteString("\n")
		}
	}

	// Pagination bar
	b.WriteString(sep + "\n")
	prevStr := DescriptionStyle().Render("[←]")
	if t.offset > 0 {
		prevStr = SelectedStyle().Render("[←]")
	}
	nextStr := DescriptionStyle().Render("[→]")
	if t.hasMore {
		nextStr = SelectedStyle().Render("[→]")
	}
	hint := ""
	if t.helpText != "" {
		hint = "  " + DescriptionStyle().Render(t.helpText)
	}
	b.WriteString(fmt.Sprintf(" %s pg%d  %d items  %s%s\n",
		prevStr, t.pageNum, len(t.rows), nextStr, hint))

	return b.String()
}
