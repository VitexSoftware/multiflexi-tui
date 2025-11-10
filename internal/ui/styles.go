package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#04B575")
	accentColor    = lipgloss.Color("#F25D94")
	textColor      = lipgloss.Color("#FAFAFA")
	subtleColor    = lipgloss.Color("#626262")
	
	// Base styles
	baseStyle = lipgloss.NewStyle().
		Foreground(textColor)
	
	// Title styles
	titleStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Margin(0, 0, 1, 0)
	
	// List styles
	listStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2)
	
	selectedItemStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Padding(0, 1)
	
	unselectedItemStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Padding(0, 1)
	
	itemDescriptionStyle = lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true)
	
	// Viewer styles
	viewerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Margin(0, 0, 1, 0)
	
	// Footer styles
	footerStyle = lipgloss.NewStyle().
		Foreground(subtleColor).
		Margin(1, 0, 0, 0)
	
	// Error styles
	errorStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)
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