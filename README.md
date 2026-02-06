# MultiFlexi TUI

A modern terminal user interface (TUI) frontend for the `multiflexi-cli` tool, built with the [Charmbracelet Bubbletea](https://github.com/charmbracelet/bubbletea) framework.

## Features

- **Comprehensive CLI Coverage**: Complete parity with all `multiflexi-cli --format=json` functionality
- **Entity Management**: View and manage all MultiFlexi entities with full field coverage
- **Advanced Pagination**: Navigate through large datasets with limit/offset controls
- **Real-time Status Panel**: Live system information from `multiflexi-cli status`
- **Dynamic Menu System**: Access all entities through intuitive navigation
- **Enhanced Data Views**: Complete field coverage for Artifacts, CredTypes, and new CrPrototypes
- **Professional UI**: Responsive three-panel design with clean styling
- **Keyboard Navigation**: Efficient controls for all operations

### Key Features

#### 🎯 **Complete Entity Management**

- **All MultiFlexi Entities**: Full coverage of Applications, Companies, Jobs, Users, RunTemplates, Credentials, Tokens, Artifacts, CredTypes, CrPrototypes, CompanyApps, Encryption, Queue, and Prune operations
- **Enhanced Data Views**: Complete field coverage with proper JSON parsing
- **Pagination Controls**: Navigate through large datasets with configurable limits
- **Real-time Updates**: Refresh data with 'r' key for current state

#### 🧭 **Comprehensive Navigation Menu**

- **Complete Entity Access**: Status | RunTemplates | Jobs | Applications | Companies | Credentials | Tokens | Users | Artifacts | CredTypes | CrPrototypes | CompanyApps | Encryption | Queue | Prune | Commands | Help | Quit
- **Context-Aware Hints**: Dynamic descriptions for each entity and operation
- **Seamless Navigation**: Arrow key navigation with visual feedback and focus management

#### 📊 **Real-time Status Panel**

- **System Information**: CLI version, database migration status, user info
- **Live Updates**: Refresh with 'r' key for current system state
- **JSON Parsing**: Clean, formatted display of status data

#### 🎨 **Professional UI Design**

- **Three-Panel Layout**: Menu at top, content in middle, status at bottom
- **Responsive Design**: Adapts to different terminal sizes
- **Color-Coded Elements**: Clear visual hierarchy and status indication
- **Consistent Styling**: Professional appearance throughout

### Navigation Summary

| Key | Action | Context |
|-----|--------|---------|
| `←/→` or `h/l` | Navigate top menu | Global navigation |
| `Enter` or `Space` | Select menu item | Menu selection |
| `↑/↓` or `k/j` | Navigate within lists | Content navigation |
| `Shift+←/→` | Previous/next job pages | Jobs pagination |
| `Tab` | Switch between views | View switching |
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

## Usage

Simply run the application:

```bash
multiflexi-tui
```

The application will launch with the Jobs dashboard as the default view. Use the keyboard navigation to explore different sections.

## Project Structure

```text
multiflexi-tui/
├── cmd/
│   └── multiflexi-tui/
│       └── main.go          # Application entry point
├── internal/
│   ├── app/
│   │   └── app.go           # Application coordination and state management
│   ├── cli/
│   │   └── cli.go           # MultiFlexi CLI integration
│   └── ui/
│       ├── menu.go          # Command list interface
│       ├── viewer.go        # Help text viewer
│       └── styles.go        # UI styling with Lipgloss
├── debian/                  # Debian packaging files
│   ├── control
│   ├── rules
│   ├── install
│   └── changelog
├── go.mod                   # Go module definition
├── Makefile                 # Build automation
└── README.md               # This file
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

1. Application starts and loads commands using `multiflexi-cli describe`
2. Commands are displayed in an interactive list
3. User selects a command to view its help
4. Help is loaded using `multiflexi-cli <command> --help`
5. Help text is displayed in a scrollable viewer
6. User can return to menu or exit

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
