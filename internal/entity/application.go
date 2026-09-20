package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var ApplicationDef = &EntityDef{
	Name: "📦 Applications", CLIEntity: "application", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 5, Field: "id"},
		{Header: "Name", Width: 30, Field: "name", Flex: true},
		{Header: "Version", Width: 15, Field: "version"},
		{Header: "Status", Width: 10, Field: "status"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Application
		if err := c.List("application", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, a := range items {
			status := "Disabled"
			if a.Enabled == 1 {
				status = "Enabled"
			}
			rows[i] = ui.TableRow{ID: a.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", a.ID), "name": a.Name, "version": a.Version, "status": status,
			}, FullData: a}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		a := data.(cli.Application)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", a.ID)}, {Label: "Name", Value: a.Name},
			{Label: "Version", Value: a.Version}, {Label: "UUID", Value: a.UUID},
			{Label: "Executable", Value: a.Executable}, {Label: "Description", Value: a.Description},
			{Label: "Homepage", Value: a.Homepage}, {Label: "Tags", Value: a.Tags},
			{Label: "Enabled", Value: fmt.Sprintf("%d", a.Enabled)},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		a := data.(cli.Application)
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Application name", Value: a.Name},
			{Label: "Description", Placeholder: "Description", Value: a.Description},
			{Label: "Executable", Placeholder: "Executable path", Value: a.Executable},
			{Label: "Homepage", Placeholder: "Homepage URL", Value: a.Homepage},
			{Label: "Tags", Placeholder: "Tags", Value: a.Tags},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		a := data.(cli.Application)
		args := []string{"--id", fmt.Sprintf("%d", a.ID), "--name", fields["Name"]}
		if v := fields["Description"]; v != "" {
			args = append(args, "--description", v)
		}
		if v := fields["Executable"]; v != "" {
			args = append(args, "--executable", v)
		}
		if v := fields["Homepage"]; v != "" {
			args = append(args, "--homepage", v)
		}
		if v := fields["Tags"]; v != "" {
			args = append(args, "--tags", v)
		}
		return args
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Application name", Required: true},
			{Label: "UUID", Placeholder: "UUID", Required: true},
			{Label: "Executable", Placeholder: "Executable path", Required: true},
			{Label: "Description", Placeholder: "Description"},
			{Label: "Homepage", Placeholder: "Homepage URL"},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		args := []string{"--name", fields["Name"], "--uuid", fields["UUID"], "--executable", fields["Executable"]}
		if v := fields["Description"]; v != "" {
			args = append(args, "--description", v)
		}
		if v := fields["Homepage"]; v != "" {
			args = append(args, "--homepage", v)
		}
		return args
	},
	GetID:    func(data interface{}) int { return data.(cli.Application).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("App: %s", data.(cli.Application).Name) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{
			Label:   "Config",
			Key:     "c",
			Command: "show-config",
			Handler: func(c cli.Client, data interface{}) tea.Cmd {
				a := data.(cli.Application)
				return func() tea.Msg {
					def := NewConfFieldDef(a.ID, a.Name)
					return ui.NavigateToMsg{View: NewListViewForEntity(c, def)}
				}
			},
		},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() {
	Register(Entry{Label: "Applications", Hint: "Browse applications • c: config fields", Def: ApplicationDef})
}
