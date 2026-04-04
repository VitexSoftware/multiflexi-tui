package entity

import (
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// EntityDef defines everything needed to manage one entity type.
// All type-specific logic is captured in closures.
type EntityDef struct {
	Name         string // display name
	CLIEntity    string // multiflexi-cli subcommand
	DeleteAction string // "delete" or "remove"

	Columns []ui.TableColumn
	Limit   int

	// Fetch returns raw data and converts to TableRows.
	Fetch   func(c cli.Client, limit, offset int) ([]ui.TableRow, error)

	// Detail fields from a row's FullData.
	ToDetail func(data interface{}) []ui.DetailField

	// Editor fields from a row's FullData (nil = read-only entity).
	ToEditor func(data interface{}) []ui.EditorField

	// Build CLI args for update from edited fields. Returns args like ["--name", "foo", "--id", "5"].
	UpdateArgs func(data interface{}, fields map[string]string) []string

	// Editor fields for creating a new entity (nil = create not supported).
	NewFields func() []ui.EditorField

	// Build CLI args for create from editor fields.
	CreateArgs func(fields map[string]string) []string

	// GetID extracts the entity ID from FullData.
	GetID func(data interface{}) int

	// GetLabel returns a display label for the entity instance.
	GetLabel func(data interface{}) string

	// Actions for the detail view (edit, delete, plus entity-specific ones).
	Actions []ui.ActionDef
}

// Entry is a menu-compatible wrapper around an EntityDef.
type Entry struct {
	Label string
	Hint  string
	Def   *EntityDef
}

// All registered entity entries, populated by init() in each entity file.
var All []Entry

// Register adds an entity entry to the global registry.
func Register(e Entry) {
	All = append(All, e)
}

// NewListViewForEntity creates a ListView tea.Model for the given entity definition.
func NewListViewForEntity(c cli.Client, def *EntityDef) tea.Model {
	return NewListView(c, def)
}
