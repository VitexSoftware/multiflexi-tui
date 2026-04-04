package entity

import (
	"fmt"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

func jobStatus(j cli.Job) string {
	if j.PID != 0 { return "Running" }
	if j.Exitcode == -1 { return "Scheduled" }
	if j.Exitcode == 0 { return "Success" }
	return "Failed"
}

var JobDef = &EntityDef{
	Name: "💼 Jobs", CLIEntity: "job", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"},
		{Header: "Command", Width: 25, Field: "command"},
		{Header: "Status", Width: 12, Field: "status"},
		{Header: "Schedule", Width: 20, Field: "schedule"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Job
		if err := c.List("job", limit, offset, &items); err != nil { return nil, err }
		rows := make([]ui.TableRow, len(items))
		for i, j := range items {
			sched := j.Schedule
			if len(sched) >= 16 { sched = sched[11:16] }
			rows[i] = ui.TableRow{ID: j.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", j.ID), "command": j.Command,
				"status": jobStatus(j), "schedule": sched,
			}, FullData: j}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		j := data.(cli.Job)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", j.ID)},
			{Label: "Command", Value: j.Command},
			{Label: "Status", Value: jobStatus(j)},
			{Label: "PID", Value: fmt.Sprintf("%d", j.PID)},
			{Label: "Exit Code", Value: fmt.Sprintf("%d", j.Exitcode)},
			{Label: "Executor", Value: j.Executor},
			{Label: "Schedule Type", Value: j.ScheduleType},
			{Label: "Schedule", Value: j.Schedule},
			{Label: "Begin", Value: j.Begin},
			{Label: "End", Value: j.End},
			{Label: "App ID", Value: fmt.Sprintf("%d", j.AppID)},
			{Label: "Company ID", Value: fmt.Sprintf("%d", j.CompanyID)},
			{Label: "RunTemplate ID", Value: fmt.Sprintf("%d", j.RunTemplateID)},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		j := data.(cli.Job)
		return []ui.EditorField{
			{Label: "Executor", Placeholder: "Executor", Value: j.Executor},
			{Label: "Schedule Type", Placeholder: "Schedule type", Value: j.ScheduleType},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		j := data.(cli.Job)
		return []string{"--id", fmt.Sprintf("%d", j.ID), "--executor", fields["Executor"], "--schedule_type", fields["Schedule Type"]}
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "RunTemplate ID", Placeholder: "RunTemplate ID", Required: true},
			{Label: "Scheduled", Placeholder: "YYYY-MM-DD HH:MM:SS or 'now'", Required: true},
			{Label: "Executor", Placeholder: "Native", Value: "Native"},
			{Label: "Schedule Type", Placeholder: "adhoc", Value: "adhoc"},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		args := []string{"--runtemplate_id", fields["RunTemplate ID"], "--scheduled", fields["Scheduled"]}
		if v := fields["Executor"]; v != "" { args = append(args, "--executor", v) }
		if v := fields["Schedule Type"]; v != "" { args = append(args, "--schedule_type", v) }
		return args
	},
	GetID: func(data interface{}) int { return data.(cli.Job).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Job %d", data.(cli.Job).ID) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() { Register(Entry{Label: "Jobs", Hint: "View and manage jobs", Def: JobDef}) }
