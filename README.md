# MultiFlexi TUI

A modern terminal user interface (TUI) frontend for the `multiflexi-cli` tool, built with the [Charmbracelet Bubbletea](https://github.com/charmbracelet/bubbletea) framework.

## Features

- **Comprehensive CLI Coverage**: Complete parity with all `multiflexi-cli --format=json` functionality
- **Entity Management**: View, edit, and delete all MultiFlexi entities with full field coverage
- **Inline Editors**: Edit Jobs, Companies, Applications, and RunTemplates directly in the TUI
- **Delete with Confirmation**: Safe deletion with Y/N confirmation dialog for all entity types
- **Advanced Pagination**: Navigate through large datasets with limit/offset controls
- **Real-time Status Panel**: Live system information from `multiflexi-cli status`
- **Dynamic Menu System**: Access all entities through intuitive navigation with active section indicator
- **Mouse Support**: Click menu items, scroll lists with mouse wheel (GPM/xterm compatible)
- **Professional UI**: TurboVision-inspired design with green active-section highlighting
- **Keyboard Navigation**: Efficient controls for all operations

### Key Features

#### 🎯 **Complete Entity Management**

- **All MultiFlexi Entities**: Full coverage of Applications, Companies, Jobs, Users, RunTemplates, Credentials, Tokens, Artifacts, CredTypes, CrPrototypes, CompanyApps, Encryption, Queue, and Prune operations
- **Detail Views**: Press Enter/Space to view full record details with action buttons
- **Inline Editors**: Press 'e' to edit records (Jobs, Companies, Applications, RunTemplates) with multi-field forms
- **Delete with Confirmation**: Press 'd' in detail view to delete with Y/N safety prompt
- **Pagination Controls**: Navigate through large datasets with configurable limits

![Overview](docs/title-screenshot.png?raw=true)

#### 🧭 **Comprehensive Navigation Menu**

- **Complete Entity Access**: Status | RunTemplates | Jobs | Applications | Companies | Credentials | Tokens | Users | Artifacts | CredTypes | CrPrototypes | CompanyApps | Encryption | Queue | Prune | Commands | Help | Quit
- **Active Section Indicator**: Currently active menu item highlighted in green at all times
- **Context-Aware Hints**: Dynamic descriptions for each entity and operation
- **Seamless Navigation**: Arrow key and mouse click navigation with visual feedback

#### 🖱️ **Mouse Support**

- **Menu Clicks**: Click menu items to navigate directly
- **Content Focus**: Click content area to switch focus from menu
- **Scroll Wheel**: Scroll through list items with mouse wheel
- **GPM Compatible**: Works in Linux console via GPM and all xterm-compatible terminals

#### 📊 **Real-time Status Panel**

- **System Information**: CLI version, database migration status, user info
- **Live Updates**: Refresh with 'r' key for current system state
- **JSON Parsing**: Clean, formatted display of status data

#### 🎨 **Professional UI Design**

- **Three-Panel Layout**: Menu at top, content in middle, status at bottom
- **Active Section Highlight**: Green menu item shows which section is in use
- **Responsive Design**: Adapts to different terminal sizes
- **Color-Coded Elements**: Clear visual hierarchy and status indication
- **Consistent Styling**: TurboVision-inspired appearance throughout

### Navigation Summary

| Key | Action | Context |
|-----|--------|---------|
| `←/→` or `h/l` | Navigate top menu | Menu focused |
| `Enter` or `Space` | Select menu item / Open detail | Menu / List |
| `↑/↓` or `k/j` | Navigate within lists | Content focused |
| `←/→` | Previous/next page | Content focused |
| `e` | Open editor for selected record | List / Detail view |
| `d` | Delete record (with confirmation) | Detail view |
| `Tab` | Switch focus between menu and content | Global |
| `Esc` | Go back (detail→list, editor→list) | Detail / Editor |
| `Mouse click` | Select menu item or focus content | Global |
| `Mouse wheel` | Scroll list up/down | Content focused |
| `r` | Refresh status and data | Data refresh |
| `q` or `Ctrl+C` | Quit application | Exit |

#### Menu Options & Entity Coverage

- **Status**: System dashboard with live information
- **RunTemplates**: Execution templates with pagination
- **Jobs**: Running and historical job management
- **Applications**: MultiFlexi application catalog
- **Companies**: Registered company management
- **Credentials**: Authentication credential management
- **Tokens**: API token management
- **Users**: User account management
- **Artifacts**: Job artifacts and files (enhanced: 7 fields)
- **CredTypes**: Credential type definitions (enhanced: 8 fields)
- **CrPrototypes**: Credential prototypes (new entity: 10 fields)
- **CompanyApps**: Company-application relationships
- **Encryption**: System encryption management
- **Queue**: Job queue management
- **Prune**: Log and data cleanup operations
- **Commands**: CLI command documentation
- **Help**: Application usage documentation

## Prerequisites

- Go 1.21 or later
- `multiflexi-cli` installed and available in PATH
- For Debian packaging: `debhelper-compat`, `golang-any`, `dpkg-dev`

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/VitexSoftware/multiflexi-tui.git
cd multiflexi-tui

# Build the binary
make build

# Install locally (optional)
make install
```

### Debian Package

```bash
# Build Debian package
make deb

# Or use dpkg-buildpackage directly
dpkg-buildpackage -us -uc

# Install the generated .deb package
sudo dpkg -i ../multiflexi-tui_1.0.0-1_amd64.deb
```

## Usage

Simply run the application:

```bash
multiflexi-tui
```

The application will launch with the Status dashboard as the default view. Use keyboard or mouse to navigate between sections. The active section is highlighted in green in the menu bar.

## Project Structure

```text
multiflexi-tui/
├── cmd/
│   └── multiflexi-tui/
│       └── main.go              # Application entry point
├── internal/
│   ├── app/
│   │   ├── app.go               # Application coordination, routing, mouse handling
│   │   └── model.go             # Application state model and view states
│   ├── cli/
│   │   └── cli.go               # MultiFlexi CLI integration (CRUD operations)
│   └── ui/
│       ├── styles.go            # UI styling with Lipgloss (TurboVision theme)
│       ├── detail.go            # Reusable detail widget
│       ├── detailview.go        # Detail view model with edit/delete actions
│       ├── job_editor.go        # Job editor (command, executor, schedule_type)
│       ├── company_editor.go    # Company editor (name, email, IC, slug)
│       ├── application_editor.go # Application editor
│       ├── runtemplate_editor.go # RunTemplate editor
│       ├── confirm_dialog.go    # Delete confirmation dialog (Y/N)
│       ├── jobs.go              # Jobs listing view
│       ├── companies.go         # Companies listing view
│       ├── applications.go      # Applications listing view
│       ├── runtemplates.go      # RunTemplates listing view
│       ├── menu.go              # Command list interface
│       ├── viewer.go            # Help text viewer
│       └── ...                  # Other entity listing views
├── debian/                      # Debian packaging files
├── go.mod                       # Go module definition
├── Makefile                     # Build automation
└── README.md                    # This file
```

## Development

### Building

```bash
# Development build (with debug info)
make dev

# Production build (optimized)
make build

# Run tests
make test

# Clean build artifacts
make clean
```

### Dependencies

The project uses only packages available in Debian stable repositories:

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - TUI components
- `github.com/charmbracelet/lipgloss` - Styling library

### Makefile Targets

- `make build` - Build optimized binary
- `make dev` - Build development binary
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make install` - Install binary locally
- `make deb` - Build Debian package
- `make deps` - Download and tidy dependencies
- `make run` - Build and run the application
- `make check-deps` - Verify required tools are installed

## Architecture

The application follows a clean modular architecture:

1. **CLI Layer** (`internal/cli`): Handles communication with `multiflexi-cli`
2. **UI Layer** (`internal/ui`): Implements Bubbletea models and views
3. **App Layer** (`internal/app`): Coordinates between UI components and manages state
4. **Main** (`cmd/multiflexi-tui`): Application entry point

### Flow

1. Application starts and loads system status from `multiflexi-cli status`
2. User navigates the menu bar to select an entity type (Jobs, Companies, etc.)
3. Entity listings are loaded from `multiflexi-cli <entity> list --format=json`
4. User can view details (Enter), edit (e), or delete (d) records
5. Editors provide multi-field forms; deletions require Y/N confirmation
6. All CRUD operations are executed via the corresponding `multiflexi-cli` commands

## Package Information

- **Source**: multiflexi-tui
- **Section**: utils
- **Priority**: optional
- **Maintainer**: Vitex Software <info@vitexsoftware.cz>
- **Homepage**: https://github.com/VitexSoftware/multiflexi-cli
- **Dependencies**: `multiflexi-cli`

## License

This project follows the same license as the MultiFlexi project.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make test` to ensure tests pass
6. Submit a pull request

## Support

For issues and questions, please visit the [GitHub repository](https://github.com/VitexSoftware/multiflexi-tui) or contact Vitex Software at <info@vitexsoftware.cz>.
