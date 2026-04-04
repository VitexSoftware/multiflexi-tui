package main

import (
	"fmt"
	"os"

	"github.com/VitexSoftware/multiflexi-tui/internal/app"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/entity"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	client := cli.NewCLIClient()

	// Build menu items: Status (home) + all registered entities + Help + Quit
	items := []app.MenuItem{
		{
			Label: "Status",
			Hint:  "View system dashboard with status information",
			Action: func(a *app.App) (tea.Model, tea.Cmd) {
				return nil, nil // nil view = show status dashboard
			},
		},
	}

	// Add all registered entities from the registry
	for _, e := range entity.All {
		entry := e // capture loop variable
		items = append(items, app.MenuItem{
			Label: entry.Label,
			Hint:  entry.Hint,
			Action: func(a *app.App) (tea.Model, tea.Cmd) {
				view := entity.NewListViewForEntity(a.Client, entry.Def)
				return view, nil
			},
		})
	}

	// Help
	items = append(items, app.MenuItem{
		Label: "Help",
		Hint:  "View help and documentation",
		Action: func(a *app.App) (tea.Model, tea.Cmd) {
			viewer := ui.NewViewer("Help")
			return viewer, func() tea.Msg {
				content, err := a.Client.GetCommandHelp("help")
				if err != nil {
					return ui.HelpErrorMsg{Command: "help", Err: err}
				}
				return ui.HelpLoadedMsg{Command: "help", Content: content}
			}
		},
	})

	// Quit
	items = append(items, app.MenuItem{
		Label: "Quit",
		Hint:  "Exit the application",
		Action: func(a *app.App) (tea.Model, tea.Cmd) {
			return nil, tea.Quit
		},
	})

	if err := app.Run(client, items); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
