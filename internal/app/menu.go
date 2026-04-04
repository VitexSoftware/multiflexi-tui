package app

import tea "github.com/charmbracelet/bubbletea"

// MenuItem defines a single menu entry with its action.
type MenuItem struct {
	Label  string
	Hint   string
	Action func(a *App) (tea.Model, tea.Cmd)
}
