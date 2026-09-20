package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func tokenUserLabel(t cli.Token) string {
	if t.User != "" {
		return t.User
	}
	if t.UserID > 0 {
		return fmt.Sprintf("#%d", t.UserID)
	}
	return ""
}

func tokenUntilLabel(t cli.Token) string {
	if t.Until == nil || *t.Until == "" {
		return "never"
	}
	return *t.Until
}

var TokenDef = &EntityDef{
	Name: "🎟️ Tokens", CLIEntity: "token", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"},
		{Header: "User", Width: 16, Field: "user"},
		{Header: "Token", Width: 28, Field: "token", Flex: true},
		{Header: "Until", Width: 20, Field: "until"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Token
		if err := c.List("token", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, t := range items {
			tok := t.Token
			if len(tok) > 12 {
				tok = tok[:6] + "…" + tok[len(tok)-4:]
			}
			rows[i] = ui.TableRow{ID: t.ID, Values: map[string]string{
				"id":    fmt.Sprintf("%d", t.ID),
				"user":  tokenUserLabel(t),
				"token": tok,
				"until": tokenUntilLabel(t),
			}, FullData: t}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		t := data.(cli.Token)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", t.ID)},
			{Label: "User ID", Value: fmt.Sprintf("%d", t.UserID)},
			{Label: "User", Value: tokenUserLabel(t)},
			{Label: "Token", Value: t.Token},
			{Label: "Start", Value: t.Start},
			{Label: "Until", Value: tokenUntilLabel(t)},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		t := data.(cli.Token)
		return []ui.EditorField{
			{Label: "User ID", Placeholder: "User ID", Value: fmt.Sprintf("%d", t.UserID)},
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
				userArg := fmt.Sprintf("%d", t.UserID)
				if t.UserID == 0 && t.User != "" {
					userArg = t.User
				}
				return func() tea.Msg {
					args := []string{"token:generate", "--format=json"}
					if t.UserID > 0 {
						args = append(args, "--user", userArg)
					} else {
						args = append(args, "--login", userArg)
					}
					output, err := c.RunRaw(args...)
					viewer := ui.NewViewer(fmt.Sprintf("Generate Token for User %s", userArg))
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
