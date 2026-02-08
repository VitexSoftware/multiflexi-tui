package ui

import "github.com/charmbracelet/lipgloss"

var (
	// TurboVision Classic Colors
	tvBlue      = lipgloss.Color("#0000AA") // Classic blue background
	tvCyan      = lipgloss.Color("#00AAAA") // Cyan for highlights
	tvWhite     = lipgloss.Color("#FFFFFF") // White text
	tvLightGray = lipgloss.Color("#CCCCCC") // Light gray text
	tvDarkGray  = lipgloss.Color("#808080") // Dark gray for disabled
	tvYellow    = lipgloss.Color("#FFFF00") // Yellow for selections
	tvRed       = lipgloss.Color("#FF0000") // Red for errors
	tvGreen     = lipgloss.Color("#00FF00") // Green for success

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Foreground(tvWhite).
			Background(tvBlue)

	// Title styles - White on blue with double border
	titleStyle = lipgloss.NewStyle().
			Foreground(tvWhite).
			Background(tvBlue).
			Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(tvWhite).
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	// List styles - Classic TurboVision window style
	listStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(tvWhite).
			Background(tvBlue).
			Foreground(tvWhite).
			Padding(1, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(tvBlue).
				Background(tvYellow).
				Bold(true).
				Padding(0, 1)

	unselectedItemStyle = lipgloss.NewStyle().
				Foreground(tvWhite).
				Background(tvBlue).
				Padding(0, 1)

	itemDescriptionStyle = lipgloss.NewStyle().
				Foreground(tvLightGray).
				Background(tvBlue)

	// Viewer styles - TurboVision dialog box style
	viewerStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(tvWhite).
			Background(tvBlue).
			Foreground(tvWhite).
			Padding(1, 2).
			Margin(0, 0, 1, 0)

	// Footer styles - Status bar style
	footerStyle = lipgloss.NewStyle().
			Foreground(tvBlue).
			Background(tvCyan).
			Bold(true).
			Margin(1, 0, 0, 0)

	// Error styles - Classic red on white
	errorStyle = lipgloss.NewStyle().
			Foreground(tvRed).
			Background(tvWhite).
			Bold(true)

	// Button styles - TurboVision button style
	buttonStyle = lipgloss.NewStyle().
			Foreground(tvBlue).
			Background(tvLightGray).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(tvWhite)

	disabledButtonStyle = lipgloss.NewStyle().
				Foreground(tvDarkGray).
				Background(tvBlue).
				Padding(0, 1).
				Border(lipgloss.NormalBorder()).
				BorderForeground(tvDarkGray)

	// Status table styles - TurboVision window style
	statusTableStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(tvWhite).
				Background(tvBlue).
				Foreground(tvWhite).
				Padding(1, 2)

	activeStatusStyle = lipgloss.NewStyle().
				Foreground(tvGreen).
				Background(tvBlue).
				Bold(true)

	disabledStatusStyle = lipgloss.NewStyle().
				Foreground(tvRed).
				Background(tvBlue).
				Bold(true)

	statusLabelStyle = lipgloss.NewStyle().
				Foreground(tvWhite).
				Background(tvBlue).
				Bold(true)

	statusValueStyle = lipgloss.NewStyle().
				Foreground(tvLightGray).
				Background(tvBlue)
)

// GetListStyle returns the list container style
func GetListStyle() lipgloss.Style {
	return listStyle
}

// GetTitleStyle returns the title style
func GetTitleStyle() lipgloss.Style {
	return titleStyle
}

// GetSelectedItemStyle returns the selected item style
func GetSelectedItemStyle() lipgloss.Style {
	return selectedItemStyle
}

// GetUnselectedItemStyle returns the unselected item style
func GetUnselectedItemStyle() lipgloss.Style {
	return unselectedItemStyle
}

// GetItemDescriptionStyle returns the item description style
func GetItemDescriptionStyle() lipgloss.Style {
	return itemDescriptionStyle
}

// GetViewerStyle returns the viewer style
func GetViewerStyle() lipgloss.Style {
	return viewerStyle
}

// GetFooterStyle returns the footer style
func GetFooterStyle() lipgloss.Style {
	return footerStyle
}

// GetErrorStyle returns the error style
func GetErrorStyle() lipgloss.Style {
	return errorStyle
}

// GetButtonStyle returns the button style
func GetButtonStyle() lipgloss.Style {
	return buttonStyle
}

// GetDisabledButtonStyle returns the disabled button style
func GetDisabledButtonStyle() lipgloss.Style {
	return disabledButtonStyle
}

// GetStatusTableStyle returns the status table style
func GetStatusTableStyle() lipgloss.Style {
	return statusTableStyle
}

// GetActiveStatusStyle returns the active status style (green)
func GetActiveStatusStyle() lipgloss.Style {
	return activeStatusStyle
}

// GetDisabledStatusStyle returns the disabled status style (red)
func GetDisabledStatusStyle() lipgloss.Style {
	return disabledStatusStyle
}

// GetStatusLabelStyle returns the status label style
func GetStatusLabelStyle() lipgloss.Style {
	return statusLabelStyle
}

// GetStatusValueStyle returns the status value style
func GetStatusValueStyle() lipgloss.Style {
	return statusValueStyle
}
