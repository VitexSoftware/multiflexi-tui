package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var TokenDef = &EntityDef{
	Name: "🎟️ Tokens", CLIEntity: "token", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "User", Width: 20, Field: "user"},
		{Header: "Token", Width: 45, Field: "token"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Token
		if err := c.List("token", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, t := range items {
			rows[i] = ui.TableRow{ID: t.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", t.ID), "user": t.User, "token": t.Token,
			}, FullData: t}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		t := data.(cli.Token)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", t.ID)},
			{Label: "User", Value: t.User},
			{Label: "Token", Value: t.Token},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		t := data.(cli.Token)
		return []ui.EditorField{
			{Label: "User ID", Placeholder: "User ID", Value: t.User},
			{Label: "Token", Placeholder: "Token value", Value: t.Token},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		t := data.(cli.Token)
		return []string{"--id", fmt.Sprintf("%d", t.ID), "--user", fields["User ID"], "--token", fields["Token"]}
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "User ID", Placeholder: "User ID", Required: true},
			{Label: "Token", Placeholder: "Token value (leave blank to generate)", Value: ""},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		args := []string{"--user", fields["User ID"]}
		if v := fields["Token"]; v != "" {
			args = append(args, "--token", v)
		}
		return args
	},
	GetID:    func(data interface{}) int { return data.(cli.Token).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Token %d", data.(cli.Token).ID) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{
			Label:   "Generate",
			Key:     "g",
			Command: "generate",
			Handler: func(c cli.Client, data interface{}) tea.Cmd {
				t := data.(cli.Token)
				return func() tea.Msg {
					output, err := c.RunRaw("token", "generate", "--format=json",
						"--user", t.User)
					viewer := ui.NewViewer(fmt.Sprintf("Generate Token for User %s", t.User))
					if err != nil {
						viewer.SetContent("Generate Token", fmt.Sprintf("Error: %v\n\n%s", err, string(output)))
					} else {
						viewer.SetContent("Generated Token", string(output))
					}
					return ui.NavigateToMsg{View: viewer}
				}
			},
		},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() { Register(Entry{Label: "Tokens", Hint: "Manage API tokens", Def: TokenDef}) }
