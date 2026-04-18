# Architecture Guide

## Overview

MultiFlexi TUI uses a **data-driven entity registry** pattern built on the [Bubbletea](https://github.com/charmbracelet/bubbletea) Elm-architecture framework.

```
cmd/multiflexi-tui/main.go
        │
        └── app.New(client, menuItems)
                │
                ├── internal/app      — root model, nav stack, menu bar
                ├── internal/entity   — entity definitions + generic views
                ├── internal/cli      — CLI client interface + types
                └── internal/ui       — shared widgets and styles
```

## Layers

### CLI Layer (`internal/cli`)

`Client` is an interface that wraps `multiflexi-cli`:

```go
type Client interface {
    GetStatus() (*StatusInfo, error)
    List(entity string, limit, offset int) ([]byte, error)
    Create(entity string, args []string) error
    Update(entity string, args []string) error
    Delete(entity, action string, id int) error
    Run(args ...string) (string, error)
}
```

`CLIClient` implements it via `exec.Command`. Tests use a mock implementation.

### Entity Registry (`internal/entity`)

Each entity is an `EntityDef` struct with callbacks:

```go
type EntityDef struct {
    Name         string
    CLIEntity    string          // CLI subcommand (e.g. "company")
    DeleteAction string          // "delete" or "remove"
    Limit        int             // default page size
    Columns      []ui.TableColumn
    Fetch        func(cli.Client, int, int) ([]ui.TableRow, error)
    ToDetail     func(interface{}) []ui.DetailField
    ToEditor     func(interface{}) []ui.EditorField   // nil = no edit
    UpdateArgs   func(interface{}, map[string]string) []string
    NewFields    func() []ui.EditorField               // nil = no create
    CreateArgs   func(map[string]string) []string
    GetID        func(interface{}) int
    GetLabel     func(interface{}) string
    Actions      []ui.ActionDef       // row-level actions (shown in DetailView)
    ListActions  []ui.ListActionDef   // list-level actions (global key bindings)
}
```

Each entity file calls `Register(Entry{...})` in its `init()` — no changes needed in any other file.

### Generic Views (`internal/entity`)

| View | Description |
|------|-------------|
| `ListView` | Paginated table for any `EntityDef`. Handles `WindowSizeMsg` → `TableWidget.SetContentHeight` → re-fetch. |
| `DetailView` | Scrollable field list + action buttons for any `EntityDef`. |
| `EditorView` | Multi-field form for create and update modes. |
| `ActionFormView` | Prompted-input form that calls an arbitrary `onSave` callback (used for schedule, save-artifact, etc.). |

### UI Widgets (`internal/ui`)

| Widget | Description |
|--------|-------------|
| `TableWidget` | Paginated table with cursor. `SetContentHeight(h)` adapts row limit to terminal height. |
| `Viewer` | Scrollable text viewer with PgUp/PgDn/g/G keys and percentage indicator. |
| `ConfirmDialog` | Y/N modal for destructive operations. |

### App Layer (`internal/app`)

`App` is the root Bubbletea model (~460 lines). Responsibilities:

- **Menu bar**: horizontal scrollable bar; `adjustMenuViewport()` keeps the focused item visible.
- **Navigation stack**: `Navigator` push/pop for back-navigation.
- **Message routing**: `NavigateToMsg`, `NavigateBackMsg`, `ConfirmMsg/Yes/No`, `StatusMsg`.
- **Window sizing**: on `tea.WindowSizeMsg`, forwards `Height − 5` (chrome lines) to the active child view so it fills the content area exactly.

Chrome accounting:
```
Line 0:  menu title + menu items
Line 1:  hint line
Line 2:  ══ separator
...content area (Height − 5 lines)...
Line H-2: ══ separator
Line H-1: status message (optional)
Line H:   help/key hint line
```

## Adding a New Entity

Create a single file `internal/entity/myentity.go`:

```go
package entity

import (
    "fmt"
    "github.com/VitexSoftware/multiflexi-tui/internal/cli"
    "github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var MyEntityDef = &EntityDef{
    Name:         "My Entity",
    CLIEntity:    "myentity",
    DeleteAction: "delete",
    Limit:        10,
    Columns: []ui.TableColumn{
        {Header: "ID",   Field: "id",   Width: 6},
        {Header: "Name", Field: "name", Width: 30},
    },
    Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
        raw, err := c.List("myentity", limit, offset)
        if err != nil {
            return nil, err
        }
        var items []cli.MyEntity
        if err := json.Unmarshal(raw, &items); err != nil {
            return nil, err
        }
        rows := make([]ui.TableRow, len(items))
        for i, item := range items {
            rows[i] = ui.TableRow{
                Values:   map[string]string{"id": fmt.Sprintf("%d", item.ID), "name": item.Name},
                FullData: item,
            }
        }
        return rows, nil
    },
    ToDetail: func(data interface{}) []ui.DetailField {
        item := data.(cli.MyEntity)
        return []ui.DetailField{
            {Label: "ID",   Value: fmt.Sprintf("%d", item.ID)},
            {Label: "Name", Value: item.Name},
        }
    },
    GetID:    func(data interface{}) int    { return data.(cli.MyEntity).ID },
    GetLabel: func(data interface{}) string { return "MyEntity: " + data.(cli.MyEntity).Name },
    // Optional: ToEditor, UpdateArgs, NewFields, CreateArgs, Actions, ListActions
}

func init() {
    Register(Entry{Label: "MyEntities", Hint: "Manage my entities", Def: MyEntityDef})
}
```

The entity appears in the menu automatically — no other files need changing.

## Message Flow

```
User keypress
    → app.handleKey
        → activeView.Update(KeyMsg)          // e.g. ListView
            → returns NavigateToMsg{detail}
    → app.Update handles NavigateToMsg
        → nav.Push(currentView)
        → activeView = detail
        → detail.Init()

User presses Esc in detail
    → detail.Update → returns NavigateBackMsg
    → app.Update handles NavigateBackMsg
        → nav.Pop() → restore previous view

User presses 'd' in detail → delete
    → detail.executeActionByCommand("delete")
        → returns ConfirmMsg{label, action}
    → app.Update handles ConfirmMsg
        → nav.Push(detail); activeView = ConfirmDialog
    → user presses 'y'
        → ConfirmDialog → ConfirmYesMsg{action}
    → app.Update handles ConfirmYesMsg
        → nav.Pop() → restore detail
        → runs action() as tea.Cmd
```

## Testing

```bash
make test        # go test ./...
```

Key test files:
- `internal/entity/entity_test.go` — exercises ToDetail/ToEditor/UpdateArgs/CreateArgs/GetID/GetLabel for all entities
- `internal/cli/client_test.go` — CLI client parsing tests
- `internal/ui/` — widget tests (ConfirmDialog, Viewer)
