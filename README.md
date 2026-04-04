# MultiFlexi TUI

A modern terminal user interface (TUI) frontend for the `multiflexi-cli` tool, built with the [Charmbracelet Bubbletea](https://github.com/charmbracelet/bubbletea) framework.

## Features

- **Complete Entity Management**: List, view details, create, edit, and delete all MultiFlexi entities
- **14 Entity Types**: Applications, Companies, Jobs, RunTemplates, Credentials, Tokens, Users, Artifacts, CredTypes, CrPrototypes, CompanyApps, Queue, EventSources, EventRules
- **Create Operations**: Create new entities directly from the TUI (Companies, Jobs, Applications, RunTemplates)
- **Inline Editors**: Multi-field forms for editing records with Tab/Arrow key navigation
- **Delete with Confirmation**: Safe deletion with Y/N confirmation dialog
- **Pagination**: Navigate through large datasets with limit/offset controls
- **Status Dashboard**: Live system information from `multiflexi-cli status`
- **Data-Driven Menu**: Dynamic menu bar with active section highlighting
- **Navigation Stack**: Full back-navigation history (list → detail → editor → confirm)
- **Mouse Support**: Click menu items, scroll lists with mouse wheel (GPM/xterm compatible)
- **TurboVision Theme**: Professional TurboVision-inspired design

### Navigation

| Key | Action | Context |
|-----|--------|---------|
| `←/→` or `h/l` | Navigate top menu | Menu focused |
| `Enter` or `Space` | Select menu item / Open detail | Menu / List |
| `↑/↓` or `k/j` | Navigate within lists | Content focused |
| `←/→` | Previous/next page | Content focused |
| `e` | Edit selected record | List / Detail view |
| `n` | Create new record | List view |
| `d` | Delete record (with confirmation) | Detail view |
| `Tab` | Switch focus between menu and content | Global |
| `Esc` | Go back to previous view | Any nested view |
| `Mouse click` | Select menu item or focus content | Global |
| `Mouse wheel` | Scroll list up/down | Content focused |
| `r` | Refresh data | List view |
| `q` or `Ctrl+C` | Quit application | Global |

### Menu Options

- **Status**: System dashboard with live information
- **Companies**: Company management (list, detail, create, edit, delete)
- **Jobs**: Job management (list, detail, create, edit, delete)
- **Applications**: Application catalog (list, detail, create, edit, delete)
- **RunTemplates**: Execution templates (list, detail, create, edit, delete)
- **Credentials**: Credential management (list, detail, delete)
- **Tokens**: API token management (list, detail, delete)
- **Users**: User account management (list, detail, delete)
- **Artifacts**: Job artifacts and output files (list, detail)
- **CredTypes**: Credential type definitions (list, detail)
- **CrPrototypes**: Credential prototypes (list, detail)
- **CompanyApps**: Company-application relationships (list, detail, delete)
- **Queue**: Job queue (list, detail)
- **EventSources**: Event source adapters (list, detail, delete)
- **EventRules**: Event-to-RunTemplate mappings (list, detail, delete)
- **Help**: Application usage documentation
- **Quit**: Exit the application

## Prerequisites

- Go 1.21 or later
- `multiflexi-cli` installed and available in PATH
- For Debian packaging: `debhelper-compat`, `golang-any`, `dpkg-dev`

## Installation

### From Source

```bash
git clone https://github.com/VitexSoftware/multiflexi-tui.git
cd multiflexi-tui
make build
make install  # optional
```

### Debian Package

```bash
make deb
sudo dpkg -i ../multiflexi-tui_*.deb
```

## Usage

```bash
multiflexi-tui
```

The application launches with the Status dashboard. Use keyboard or mouse to navigate.

## Project Structure

```text
multiflexi-tui/
├── cmd/multiflexi-tui/
│   └── main.go                  # Entry point — wires CLI client, entity registry, and menu
├── internal/
│   ├── app/
│   │   ├── app.go               # Lean app model (~370 lines) — Update/View/menu/mouse/keys
│   │   ├── navigator.go         # Navigation stack for back-navigation
│   │   └── menu.go              # MenuItem data type
│   ├── cli/
│   │   ├── client.go            # Client interface + CLIClient (exec.Command wrapper)
│   │   └── types.go             # All entity structs (14 types + StatusInfo)
│   ├── entity/
│   │   ├── registry.go          # EntityDef + global registry
│   │   ├── list_view.go         # Generic ListView (works with any EntityDef)
│   │   ├── detail_view.go       # Generic DetailView with action buttons
│   │   ├── editor_view.go       # Generic EditorView (create + update mode)
│   │   ├── company.go           # Company entity definition
│   │   ├── job.go               # Job entity definition
│   │   ├── application.go       # Application entity definition
│   │   ├── runtemplate.go       # RunTemplate entity definition
│   │   ├── credential.go        # Credential entity definition
│   │   ├── token.go             # Token entity definition
│   │   ├── user.go              # User entity definition
│   │   ├── artifact.go          # Artifact entity definition
│   │   ├── credtype.go          # CredType entity definition
│   │   ├── crprototype.go       # CrPrototype entity definition
│   │   ├── companyapp.go        # CompanyApp entity definition
│   │   ├── queue.go             # Queue entity definition
│   │   ├── eventsource.go       # EventSource entity definition (NEW)
│   │   └── eventrule.go         # EventRule entity definition (NEW)
│   └── ui/
│       ├── messages.go          # Shared message types (NavigateToMsg, etc.)
│       ├── styles.go            # TurboVision theme styles
│       ├── table.go             # Reusable table widget with pagination
│       ├── confirm_dialog.go    # Y/N confirmation dialog
│       └── viewer.go            # Scrollable help text viewer
├── debian/                      # Debian packaging files
├── go.mod                       # Go module definition
├── Makefile                     # Build automation
└── README.md                    # This file
```

## Architecture

The application uses a **data-driven entity registry** pattern:

1. **CLI Layer** (`internal/cli`): `Client` interface wraps `multiflexi-cli`. All operations use `--format=json`. Testable via mock implementation.
2. **Entity Registry** (`internal/entity`): Each entity is a self-contained `EntityDef` with callbacks for fetch, detail rendering, editor fields, CLI args building, and actions. Adding a new entity requires only one ~40-line file.
3. **Generic Views** (`internal/entity`): `ListView`, `DetailView`, and `EditorView` are entity-agnostic — they render any `EntityDef` with zero type switches.
4. **App Layer** (`internal/app`): Lean coordinator (~370 lines) with a navigation stack, data-driven menu bar, and message routing.
5. **UI Layer** (`internal/ui`): Shared widgets (table, confirm dialog, viewer) and styles.

### Adding a New Entity

Create a single file in `internal/entity/` (e.g. `myentity.go`):

```go
package entity

var MyEntityDef = &EntityDef{
    Name: "My Entity", CLIEntity: "myentity", DeleteAction: "delete", Limit: 10,
    Columns: []ui.TableColumn{...},
    Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) { ... },
    ToDetail: func(data interface{}) []ui.DetailField { ... },
    // Optional: ToEditor, UpdateArgs, NewFields, CreateArgs for edit/create support
    GetID: func(data interface{}) int { ... },
    GetLabel: func(data interface{}) string { ... },
    Actions: []ui.ActionDef{...},
}

func init() {
    Register(Entry{Label: "MyEntity", Hint: "Description", Def: MyEntityDef})
}
```

The entity automatically appears in the menu — no changes needed elsewhere.

## Development

```bash
make dev        # Development build
make build      # Production build (optimized)
make test       # Run tests
make clean      # Clean build artifacts
make deps       # Download and tidy dependencies
make run        # Build and run
make deb        # Build Debian package
make check-deps # Verify required tools
```

### Dependencies

Packages available in Debian stable repositories:

- `github.com/charmbracelet/bubbletea` — TUI framework
- `github.com/charmbracelet/bubbles` — TUI components (textinput)
- `github.com/charmbracelet/lipgloss` — Styling library

## Package Information

- **Source**: multiflexi-tui
- **Section**: utils
- **Priority**: optional
- **Maintainer**: Vitex Software <info@vitexsoftware.cz>
- **Homepage**: https://github.com/VitexSoftware/multiflexi-tui
- **Dependencies**: `multiflexi-cli`

## License

This project follows the same license as the MultiFlexi project.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make test` to ensure tests pass
5. Submit a pull request

## Support

For issues and questions, visit the [GitHub repository](https://github.com/VitexSoftware/multiflexi-tui) or contact Vitex Software at <info@vitexsoftware.cz>.
