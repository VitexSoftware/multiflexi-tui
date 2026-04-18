package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
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
		if err := c.List("runtemplate", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, t := range items {
			status := "Active"
			if t.Active == 0 {
				status = "Inactive"
			}
			name := t.Name
			if name == "" {
				name = "<unnamed>"
			}
			rows[i] = ui.TableRow{ID: t.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", t.ID), "name": name, "app_id": fmt.Sprintf("%d", t.AppID),
				"company": fmt.Sprintf("%d", t.CompanyID), "status": status, "executor": t.Executor,
			}, FullData: t}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		t := data.(cli.RunTemplate)
		next := ""
		if t.NextSchedule != nil {
			next = *t.NextSchedule
		}
		last := ""
		if t.LastSchedule != nil {
			last = *t.LastSchedule
		}
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", t.ID)},
			{Label: "Name", Value: t.Name},
			{Label: "App ID", Value: fmt.Sprintf("%d", t.AppID)},
			{Label: "Company ID", Value: fmt.Sprintf("%d", t.CompanyID)},
			{Label: "Active", Value: fmt.Sprintf("%d", t.Active)},
			{Label: "Interval", Value: t.Interv},
			{Label: "Cron", Value: t.Cron},
			{Label: "Executor", Value: t.Executor},
			{Label: "Last Schedule", Value: last},
			{Label: "Next Schedule", Value: next},
			{Label: "Success Jobs", Value: fmt.Sprintf("%d", t.SuccessfulJobsCount)},
			{Label: "Failed Jobs", Value: fmt.Sprintf("%d", t.FailedJobsCount)},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		t := data.(cli.RunTemplate)
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Template name", Value: t.Name},
			{Label: "Interval", Placeholder: "d/w/m/n/y", Value: t.Interv},
			{Label: "Cron", Placeholder: "*/5 * * * *", Value: t.Cron},
			{Label: "Executor", Placeholder: "Native", Value: t.Executor},
			{Label: "Active", Placeholder: "1 or 0", Value: fmt.Sprintf("%d", t.Active)},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		t := data.(cli.RunTemplate)
		args := []string{"--id", fmt.Sprintf("%d", t.ID), "--name", fields["Name"],
			"--interv", fields["Interval"], "--executor", fields["Executor"],
			"--active", fields["Active"]}
		if v := fields["Cron"]; v != "" {
			args = append(args, "--cron", v)
		}
		return args
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Template name", Required: true},
			{Label: "App ID", Placeholder: "Application ID", Required: true},
			{Label: "Company ID", Placeholder: "Company ID", Required: true},
			{Label: "Interval", Placeholder: "d/w/m/n/y"},
			{Label: "Cron", Placeholder: "*/5 * * * *"},
			{Label: "Executor", Placeholder: "Native", Value: "Native"},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		args := []string{"--name", fields["Name"], "--app_id", fields["App ID"],
			"--company_id", fields["Company ID"], "--executor", fields["Executor"]}
		if v := fields["Interval"]; v != "" {
			args = append(args, "--interv", v)
		}
		if v := fields["Cron"]; v != "" {
			args = append(args, "--cron", v)
		}
		return args
	},
	GetID:    func(data interface{}) int { return data.(cli.RunTemplate).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("RunTemplate: %s", data.(cli.RunTemplate).Name) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{
			Label:   "Schedule",
			Key:     "s",
			Command: "schedule",
			Handler: func(c cli.Client, data interface{}) tea.Cmd {
				rt := data.(cli.RunTemplate)
				form := NewActionFormView(
					fmt.Sprintf("Schedule: %s", rt.Name),
					[]ui.EditorField{
						{Label: "Schedule Time", Placeholder: "YYYY-MM-DD HH:MM:SS", Value: "now"},
						{Label: "Executor", Placeholder: "Native", Value: rt.Executor},
					},
					func(fields map[string]string) tea.Cmd {
						return func() tea.Msg {
							args := []string{"--id", fmt.Sprintf("%d", rt.ID),
								"--schedule_time", fields["Schedule Time"]}
							if v := fields["Executor"]; v != "" {
								args = append(args, "--executor", v)
							}
							_, err := c.RunRaw(append([]string{"runtemplate", "schedule", "--format=json"}, args...)...)
							if err != nil {
								return ui.StatusMsg{Text: fmt.Sprintf("Schedule failed: %v", err)}
							}
							return ui.NavigateBackMsg{}
						}
					},
				)
				return func() tea.Msg { return ui.NavigateToMsg{View: form} }
			},
		},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() {
	Register(Entry{Label: "RunTemplates", Hint: "View and manage run templates", Def: RunTemplateDef})
}
