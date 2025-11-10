# MultiFlexi TUI

A terminal user interface (TUI) frontend for the `multiflexi-cli` tool, built with the [Charmbracelet Bubbletea](https://github.com/charmbracelet/bubbletea) framework.

## Features

- **Dynamic Command Discovery**: Automatically loads available commands from `multiflexi-cli describe`
- **Interactive Menu**: Scrollable list of commands with descriptions
- **Help Viewer**: Displays command help text in a scrollable viewer
- **Keyboard Navigation**: Intuitive controls for navigation and selection
- **Clean UI**: Styled with Lipgloss for a professional appearance

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

### Keyboard Controls

- **↑/↓**: Navigate through the command list
- **Enter**: Select a command to view its help
- **q**: Go back to menu (from help view) or quit (from menu)
- **Ctrl+C**: Exit the application
- **/** (in menu): Filter/search commands

## Project Structure

```
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

For issues and questions, please visit the [GitHub repository](https://github.com/VitexSoftware/multiflexi-tui) or contact Vitex Software at info@vitexsoftware.cz.