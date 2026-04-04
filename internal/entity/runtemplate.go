package entity

import (
	"fmt"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var RunTemplateDef = &EntityDef{
	Name: "📋 Run Templates", CLIEntity: "runtemplate", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 5, Field: "id"}, {Header: "Name", Width: 25, Field: "name"},
		{Header: "App ID", Width: 8, Field: "app_id"}, {Header: "Company", Width: 10, Field: "company"},
		{Header: "Status", Width: 8, Field: "status"}, {Header: "Executor", Width: 12, Field: "executor"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.RunTemplate
		if err := c.List("runtemplate", limit, offset, &items); err != nil { return nil, err }
		rows := make([]ui.TableRow, len(items))
		for i, t := range items {
			status := "Active"; if t.Active == 0 { status = "Inactive" }
			name := t.Name; if name == "" { name = "<unnamed>" }
			rows[i] = ui.TableRow{ID: t.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", t.ID), "name": name, "app_id": fmt.Sprintf("%d", t.AppID),
				"company": fmt.Sprintf("%d", t.CompanyID), "status": status, "executor": t.Executor,
			}, FullData: t}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		t := data.(cli.RunTemplate)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", t.ID)}, {Label: "Name", Value: t.Name},
			{Label: "App ID", Value: fmt.Sprintf("%d", t.AppID)}, {Label: "Company ID", Value: fmt.Sprintf("%d", t.CompanyID)},
			{Label: "Active", Value: fmt.Sprintf("%d", t.Active)}, {Label: "Interval", Value: t.Interv},
			{Label: "Executor", Value: t.Executor}, {Label: "Cron", Value: t.Cron},
			{Label: "Success Jobs", Value: fmt.Sprintf("%d", t.SuccessfulJobsCount)},
			{Label: "Failed Jobs", Value: fmt.Sprintf("%d", t.FailedJobsCount)},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		t := data.(cli.RunTemplate)
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Template name", Value: t.Name},
			{Label: "Interval", Placeholder: "d/w/m/n", Value: t.Interv},
			{Label: "Executor", Placeholder: "Executor", Value: t.Executor},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		t := data.(cli.RunTemplate)
		return []string{"--id", fmt.Sprintf("%d", t.ID), "--name", fields["Name"], "--interv", fields["Interval"], "--executor", fields["Executor"]}
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Template name", Required: true},
			{Label: "App ID", Placeholder: "Application ID", Required: true},
			{Label: "Company ID", Placeholder: "Company ID", Required: true},
			{Label: "Executor", Placeholder: "Native", Value: "Native"},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		return []string{"--name", fields["Name"], "--app_id", fields["App ID"], "--company_id", fields["Company ID"], "--executor", fields["Executor"]}
	},
	GetID: func(data interface{}) int { return data.(cli.RunTemplate).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("RunTemplate: %s", data.(cli.RunTemplate).Name) },
	Actions: []ui.ActionDef{{Label: "Edit", Key: "e", Command: "edit"}, {Label: "Delete", Key: "d", Command: "delete"}},
}
func init() { Register(Entry{Label: "RunTemplates", Hint: "View and manage run templates", Def: RunTemplateDef}) }
