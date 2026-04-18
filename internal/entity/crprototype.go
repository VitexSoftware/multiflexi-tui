package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var CrPrototypeDef = &EntityDef{
	Name: "🧬 Credential Prototypes", CLIEntity: "crprototype", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 6, Field: "id"}, {Header: "Code", Width: 20, Field: "code"},
		{Header: "Name", Width: 30, Field: "name"}, {Header: "Version", Width: 10, Field: "version"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.CrPrototype
		if err := c.List("crprototype", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, p := range items {
			rows[i] = ui.TableRow{ID: p.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", p.ID), "code": p.Code, "name": p.Name, "version": p.Version,
			}, FullData: p}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		p := data.(cli.CrPrototype)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", p.ID)},
			{Label: "UUID", Value: p.UUID},
			{Label: "Code", Value: p.Code},
			{Label: "Name", Value: p.Name},
			{Label: "Description", Value: p.Description},
			{Label: "Version", Value: p.Version},
			{Label: "URL", Value: p.URL},
			{Label: "Created At", Value: p.CreatedAt},
			{Label: "Updated At", Value: p.UpdatedAt},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		p := data.(cli.CrPrototype)
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Prototype name", Value: p.Name},
			{Label: "Code", Placeholder: "Prototype code", Value: p.Code},
			{Label: "Description", Placeholder: "Description", Value: p.Description},
			{Label: "Version", Placeholder: "1.0.0", Value: p.Version},
			{Label: "URL", Placeholder: "Homepage URL", Value: p.URL},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		p := data.(cli.CrPrototype)
		args := []string{"--id", fmt.Sprintf("%d", p.ID), "--name", fields["Name"], "--code", fields["Code"]}
		if v := fields["Description"]; v != "" {
			args = append(args, "--description", v)
		}
		if v := fields["Version"]; v != "" {
			args = append(args, "--prototype-version", v)
		}
		if v := fields["URL"]; v != "" {
			args = append(args, "--url", v)
		}
		return args
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Prototype name", Required: true},
			{Label: "Code", Placeholder: "Prototype code", Required: true},
			{Label: "Description", Placeholder: "Description"},
			{Label: "Version", Placeholder: "1.0.0"},
			{Label: "URL", Placeholder: "Homepage URL"},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		args := []string{"--name", fields["Name"], "--code", fields["Code"]}
		if v := fields["Description"]; v != "" {
			args = append(args, "--description", v)
		}
		if v := fields["Version"]; v != "" {
			args = append(args, "--prototype-version", v)
		}
		if v := fields["URL"]; v != "" {
			args = append(args, "--url", v)
		}
		return args
	},
	GetID:    func(data interface{}) int { return data.(cli.CrPrototype).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("CrPrototype: %s", data.(cli.CrPrototype).Name) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
	ListActions: []ui.ListActionDef{
		{
			Label:   "Sync All",
			Key:     "s",
			Confirm: "Sync all credential prototypes from remote?",
			Handler: func(c cli.Client) tea.Cmd {
				return func() tea.Msg {
					output, err := c.RunRaw("crprototype", "sync", "--format=json")
					if err != nil {
						return ui.StatusMsg{Text: fmt.Sprintf("Sync failed: %v", err)}
					}
					viewer := ui.NewViewer("CrPrototype Sync Result")
					viewer.SetContent("CrPrototype Sync Result", string(output))
					return ui.NavigateToMsg{View: viewer}
				}
			},
		},
	},
}

func init() {
	Register(Entry{Label: "CrPrototypes", Hint: "Manage credential prototypes • s: sync all", Def: CrPrototypeDef})
}
