# MultiFlexi TUI

A modern terminal user interface for [multiflexi-cli](https://github.com/VitexSoftware/multiflexi-cli), built with the [Charmbracelet Bubbletea](https://github.com/charmbracelet/bubbletea) framework.

![Title Screen](docs/title-screenshot.png?raw=true)

## Features

- **Full Entity CRUD**: List, view details, create, edit, and delete all 14 MultiFlexi entity types
- **Entity-specific Actions**: Beyond CRUD — schedule runs, view stdout/stderr, generate tokens, test event sources, save artifacts, sync credential prototypes, and more
- **Dynamic Terminal Viewport**: Tables and viewers fill the full terminal height automatically (Midnight Commander style); all views reflow on resize
- **Scrollable Menu Bar**: When menu items exceed screen width, the bar scrolls to keep the focused item visible
- **Navigation Stack**: Full back-navigation history (list → detail → editor → confirm → back)
- **Delete with Confirmation**: Y/N confirmation dialog for all destructive operations
- **Pagination**: Navigate large datasets with limit/offset controls auto-sized to terminal height
- **Status Dashboard**: Live system information from `multiflexi-cli status`
- **Mouse Support**: Click menu items, scroll lists with mouse wheel
- **TurboVision Theme**: Classic TurboVision-inspired colour scheme

## Entity Types and Capabilities

| Entity | List | Detail | Create | Edit | Delete | Special Actions |
|--------|------|--------|--------|------|--------|-----------------|
| Companies | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| Applications | ✅ | ✅ | ✅ | ✅ | ✅ | Show Config (`s`) |
| RunTemplates | ✅ | ✅ | ✅ | ✅ | ✅ | Schedule (`s`) |
| Jobs | ✅ | ✅ | ✅ | ✅ | ✅ | View Stdout (`o`), View Stderr (`e`) |
| Credentials | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| Tokens | ✅ | ✅ | ✅ | ✅ | ✅ | Generate (`g`) |
| Users | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| Artifacts | ✅ | ✅ | — | — | ✅ | Save to file (`s`) |
| CredTypes | ✅ | ✅ | ✅ | ✅ | — | — |
| CrPrototypes | ✅ | ✅ | ✅ | ✅ | ✅ | Sync All (list: `S`) |
| CompanyApps | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| Queue | ✅ | ✅ | — | — | — | Fix (`f`), Truncate (`T`) |
| EventSources | ✅ | ✅ | ✅ | ✅ | ✅ | Test (`t`) |
| EventRules | ✅ | ✅ | ✅ | ✅ | ✅ | — |

## Keyboard Reference

### Global

| Key | Action |
|-----|--------|
| `Tab` | Toggle focus between menu bar and content |
| `Esc` | Go back to previous view |
| `Ctrl+C` | Quit |
| `q` | Quit (when menu focused) |

### Menu Bar (focused)

| Key | Action |
|-----|--------|
| `←/→` or `h/l` | Navigate menu items |
| `Enter` or `Space` | Open selected menu item |

### List View (content focused)

| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Move cursor up/down |
| `←/→` or `PgUp/PgDn` | Previous/next page |
| `Enter` or `Space` | Open detail view |
| `e` | Edit selected record |
| `n` | Create new record |
| `r` | Refresh / reload data |
| Entity-specific keys | See table above |

### Detail View

| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Scroll field list (when fields exceed screen) |
| `PgUp/PgDn` | Scroll field list by page |
| `Tab` or `→` | Cycle to next action button |
| `←` | Cycle to previous action button |
| `Enter` | Execute selected action |
| `d` | Delete (with confirmation) |
| `e` | Edit |
| Entity-specific keys | See table above |

### Editor / Form View

| Key | Action |
|-----|--------|
| `Tab` or `↓` | Next field |
| `Shift+Tab` or `↑` | Previous field |
| `Enter` | Save / submit |
| `Esc` | Cancel, go back |

### Viewer (stdout, stderr, help, config output)

| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Scroll one line |
| `PgUp/PgDn` | Scroll one page |
| `g` / `G` | Jump to top / bottom |
| `Esc` or `q` | Go back |

## Prerequisites

- Go 1.21 or later
- `multiflexi-cli` installed and available in `PATH`

## Installation

### From Source

```bash
git clone https://github.com/VitexSoftware/multiflexi-tui.git
cd multiflexi-tui
make build
sudo make install
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

The application opens with the Status dashboard. Use the keyboard or mouse to navigate.

## Project Structure

```
multiflexi-tui/
├── cmd/multiflexi-tui/
│   └── main.go              # Entry point — wires client, registry, and menu
├── internal/
│   ├── app/
│   │   ├── app.go           # Root model: menu bar, nav stack, message routing
│   │   ├── navigator.go     # Navigation stack (push/pop view states)
│   │   └── menu.go          # MenuItem type
│   ├── cli/
│   │   ├── client.go        # Client interface + CLIClient (exec.Command wrapper)
│   │   └── types.go         # All entity structs (14 types + StatusInfo)
│   ├── entity/
│   │   ├── registry.go      # EntityDef struct + global registry
│   │   ├── list_view.go     # Generic ListView — works with any EntityDef
│   │   ├── detail_view.go   # Generic DetailView with scrollable fields + action buttons
│   │   ├── editor_view.go   # Generic EditorView (create + update modes)
│   │   ├── action_form.go   # Generic action form (prompted input → CLI command)
│   │   ├── company.go
│   │   ├── job.go
│   │   ├── application.go
│   │   ├── runtemplate.go
│   │   ├── credential.go
│   │   ├── token.go
│   │   ├── user.go
│   │   ├── artifact.go
│   │   ├── credtype.go
│   │   ├── crprototype.go
│   │   ├── companyapp.go
│   │   ├── queue.go
│   │   ├── eventsource.go
│   │   └── eventrule.go
│   └── ui/
│       ├── messages.go      # Shared message types (NavigateToMsg, ConfirmMsg, …)
│       ├── styles.go        # TurboVision theme (lipgloss)
│       ├── table.go         # Paginated table widget — height-adaptive
│       ├── confirm_dialog.go
│       └── viewer.go        # Scrollable text viewer — height-adaptive
├── docs/
│   ├── ARCHITECTURE.md      # Architecture and entity extension guide
│   └── CLI_STRUCTURES.md    # CLI data structures reference
├── debian/                  # Debian packaging
├── go.mod
├── Makefile
└── README.md
```

## Architecture

The application uses a **data-driven entity registry** pattern:

1. **CLI Layer** (`internal/cli`): `Client` interface wraps `multiflexi-cli`. All operations use `--format=json`. Fully mockable for tests.
2. **Entity Registry** (`internal/entity`): Each entity is a self-contained `EntityDef` with callbacks for fetch, detail rendering, editor fields, CLI arg building, and action handlers. Each entity registers itself via `init()`.
3. **Generic Views**: `ListView`, `DetailView`, and `EditorView` are entity-agnostic — they render any `EntityDef` with no type switches.
4. **Dynamic Viewport**: Every view handles `tea.WindowSizeMsg`. The app passes a content-area height (`terminal height − 5 chrome lines`) so tables and viewers fill available space exactly.
5. **App Layer** (`internal/app`): Lean coordinator (~460 lines) with navigation stack, scrollable menu bar, and message routing.

See [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) for the full guide including how to add a new entity.

## Development

```bash
make build      # Production build (optimized)
make dev        # Development build
make test       # Run tests
make run        # Build and run
make clean      # Clean build artifacts
make deps       # Download and tidy dependencies
make deb        # Build Debian package
make check-deps # Verify required tools
```

## Package Information

- **Section**: utils
- **Priority**: optional
- **Maintainer**: Vitex Software <info@vitexsoftware.cz>
- **Homepage**: https://github.com/VitexSoftware/multiflexi-tui
- **Runtime dependency**: `multiflexi-cli`

## License

Same license as the MultiFlexi project.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make test` to ensure tests pass
5. Submit a pull request

## Support

For issues and questions visit the [GitHub repository](https://github.com/VitexSoftware/multiflexi-tui) or contact Vitex Software at <info@vitexsoftware.cz>.
