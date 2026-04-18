package ui

import "github.com/charmbracelet/lipgloss"

var (
	// TurboVision-inspired color scheme
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("21")).
			Foreground(lipgloss.Color("15"))

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("14"))

	activeMenuItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("10"))

	unselectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("7"))

	itemDescriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("4")).
			Padding(0, 1)

	activeStatusStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("10")).
				Bold(true)

	disabledStatusStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("9"))

	debugStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")) // bright yellow — visible but clearly secondary
)

// Public accessors for styles.
func TitleStyle() lipgloss.Style          { return titleStyle }
func SelectedStyle() lipgloss.Style       { return selectedItemStyle }
func ActiveMenuStyle() lipgloss.Style     { return activeMenuItemStyle }
func UnselectedStyle() lipgloss.Style     { return unselectedItemStyle }
func DescriptionStyle() lipgloss.Style    { return itemDescriptionStyle }
func FooterStyle() lipgloss.Style         { return footerStyle }
func ErrorStyle() lipgloss.Style          { return errorStyle }
func ButtonStyle() lipgloss.Style         { return buttonStyle }
func ActiveStatusStyle() lipgloss.Style   { return activeStatusStyle }
func DisabledStatusStyle() lipgloss.Style { return disabledStatusStyle }
func DebugStyle() lipgloss.Style          { return debugStyle }
