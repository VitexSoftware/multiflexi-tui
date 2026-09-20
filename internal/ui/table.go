package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// tableOverhead is the number of non-data lines rendered by View():
// title(1) + col-header(1) + top-sep(1) + bottom-sep(1) + pagination(1) = 5
const tableOverhead = 5

// minTableRows is the minimum number of data rows the table will display.
const minTableRows = 3

// indicatorWidth is the leading cursor glyph column (" " or "►").
const indicatorWidth = 1

// TableWidget renders a paginated table with cursor selection.
type TableWidget struct {
	title    string
	columns  []TableColumn
	rows     []TableRow
	cursor   int
	offset   int
	limit    int
	width    int // available content-area width (0 = use column mins only)
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

// SetContentWidth stores the available content-area width so flex columns can expand.
func (t *TableWidget) SetContentWidth(w int) {
	if w < 0 {
		w = 0
	}
	t.width = w
}

// effectiveWidths returns per-column widths after distributing leftover terminal
// space across flex columns (explicit Flex, or auto-picked wide columns).
func (t *TableWidget) effectiveWidths() []int {
	n := len(t.columns)
	widths := make([]int, n)
	if n == 0 {
		return widths
	}

	flex := make([]int, 0, n)
	for i, c := range t.columns {
		w := c.Width
		if w < len(c.Header) {
			w = len(c.Header)
		}
		if w < 1 {
			w = 1
		}
		widths[i] = w
		if c.Flex {
			flex = append(flex, i)
		}
	}
	if len(flex) == 0 {
		// Auto-flex: every column with Width >= 20, else the single widest.
		maxW, maxI := 0, 0
		for i, c := range t.columns {
			if c.Width > maxW {
				maxW, maxI = c.Width, i
			}
			if c.Width >= 20 {
				flex = append(flex, i)
			}
		}
		if len(flex) == 0 {
			flex = []int{maxI}
		}
	}

	// indicator + sum(widths) + space between columns
	used := indicatorWidth + n // n gaps: after indicator and between cols... actually
	// layout is: indicator + col0 + " " + col1 + " " + ...
	// so spaces = n (one after indicator is baked into joining: indicator+join(parts," "))
	// View does: indicator + strings.Join(rowParts, " ") → spaces between cols = n-1, plus indicator = 1 char
	used = indicatorWidth
	for i, w := range widths {
		used += w
		if i < n-1 {
			used++ // gap between columns
		}
	}

	extra := t.width - used
	if extra <= 0 || len(flex) == 0 {
		return widths
	}

	share := extra / len(flex)
	rem := extra % len(flex)
	for i, idx := range flex {
		widths[idx] += share
		if i < rem {
			widths[idx]++
		}
	}
	return widths
}

// truncateCell truncates/pads by visible width (ANSI/emoji-aware).
func truncateCell(val string, width int) string {
	if width <= 0 {
		return ""
	}
	w := lipgloss.Width(val)
	if w == width {
		return val
	}
	if w < width {
		return val + strings.Repeat(" ", width-w)
	}
	if width <= 3 {
		return lipgloss.NewStyle().MaxWidth(width).Render(val)
	}
	trimmed := lipgloss.NewStyle().MaxWidth(width - 3).Render(val)
	return trimmed + "..."
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

// View renders the table filling available height and width.
func (t *TableWidget) View() string {
	var b strings.Builder

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

	widths := t.effectiveWidths()
	totalWidth := indicatorWidth
	for i, w := range widths {
		totalWidth += w
		if i < len(widths)-1 {
			totalWidth++
		}
	}
	if t.width > totalWidth {
		totalWidth = t.width
	}
	sep := strings.Repeat("─", totalWidth)

	parts := make([]string, len(t.columns))
	for i, col := range t.columns {
		parts[i] = fmt.Sprintf("%-*s", widths[i], truncateCell(col.Header, widths[i]))
	}
	b.WriteString(" " + strings.Join(parts, " ") + "\n")
	b.WriteString(sep + "\n")

	if len(t.rows) == 0 {
		b.WriteString(DescriptionStyle().Render("  (no items)") + "\n")
	} else {
		count := len(t.rows)
		if count > t.limit {
			count = t.limit
		}
		for i := 0; i < count; i++ {
			row := t.rows[i]
			selected := i == t.cursor
			indicator := " "
			if selected {
				indicator = SelectedStyle().Render("►")
			}
			rowParts := make([]string, len(t.columns))
			for j, col := range t.columns {
				val := truncateCell(row.Values[col.Field], widths[j])
				// Preserve per-cell colours (e.g. exit-code badges); style plain cells.
				if !hasANSI(val) {
					if selected {
						val = SelectedStyle().Render(val)
					} else {
						val = UnselectedStyle().Render(val)
					}
				}
				rowParts[j] = val
			}
			b.WriteString(indicator + strings.Join(rowParts, " "))
			b.WriteString("\n")
		}
	}

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
