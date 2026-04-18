package ui

import (
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
)

// NavigateToMsg tells the app to push a new view onto the navigation stack.
type NavigateToMsg struct {
	View tea.Model
}

// NavigateBackMsg tells the app to pop the current view off the stack.
type NavigateBackMsg struct{}

// StatusMsg displays a transient message in the footer.
type StatusMsg struct {
	Text string
}

// ConfirmMsg is sent when the user confirms an action.
type ConfirmMsg struct {
	Label  string
	Action func() tea.Msg
}

// ConfirmYesMsg is sent when user confirms.
type ConfirmYesMsg struct {
	Action func() tea.Msg
}

// ConfirmNoMsg is sent when user cancels.
type ConfirmNoMsg struct{}

// HelpLoadedMsg carries loaded help text.
type HelpLoadedMsg struct {
	Command string
	Content string
}

// HelpErrorMsg carries a help load error.
type HelpErrorMsg struct {
	Command string
	Err     error
}

// DataLoadedMsg carries async-loaded data.
type DataLoadedMsg struct {
	Data interface{}
}

// DataErrorMsg carries a data loading error.
type DataErrorMsg struct {
	Err error
}

// DetailField is a single label-value pair for detail views.
type DetailField struct {
	Label string
	Value string
}

// EditorField defines one form field.
type EditorField struct {
	Label       string
	Placeholder string
	Value       string
	Required    bool
}

// ActionDef defines an action button on a detail view.
type ActionDef struct {
	Label   string
	Key     string // shortcut key
	Command string // "edit", "delete", or custom identifier

	// Handler, if non-nil, is called when this action fires instead of built-in dispatch.
	// It receives the CLI client and the entity's FullData, and returns a tea.Cmd.
	Handler func(client cli.Client, data interface{}) tea.Cmd

	// Confirm, if non-empty, shows a yes/no dialog with this message before calling Handler.
	Confirm string
}

// ListActionDef defines a list-level action (not per-row).
type ListActionDef struct {
	Label   string
	Key     string
	Handler func(client cli.Client) tea.Cmd
	Confirm string // if non-empty, prompts before calling Handler
}

// TableColumn defines a column in a table.
type TableColumn struct {
	Header string
	Width  int
	Field  string
}

// TableRow holds one row of table data.
type TableRow struct {
	ID       int
	Values   map[string]string
	FullData interface{} // original typed entity
}
