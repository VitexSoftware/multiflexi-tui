package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var TaskDef = &EntityDef{
	Name: "⏱️ Tasks", CLIEntity: "task", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"},
		{Header: "Template", Width: 10, Field: "runtemplate_id"},
		{Header: "State", Width: 14, Field: "state"},
		{Header: "Window Start", Width: 20, Field: "window_start"},
		{Header: "Deadline", Width: 20, Field: "deadline"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Task
		if err := c.List("task", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, t := range items {
			rows[i] = ui.TableRow{ID: t.ID, Values: map[string]string{
				"id":             fmt.Sprintf("%d", t.ID),
				"runtemplate_id": fmt.Sprintf("%d", t.RunTemplateID),
				"state":          t.State,
				"window_start":   t.WindowStart,
				"deadline":       t.Deadline,
			}, FullData: t}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		t := data.(cli.Task)
		fields := []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", t.ID)},
			{Label: "RunTemplate ID", Value: fmt.Sprintf("%d", t.RunTemplateID)},
			{Label: "State", Value: t.State},
			{Label: "Window Start", Value: t.WindowStart},
			{Label: "Window End", Value: t.WindowEnd},
			{Label: "Deadline", Value: t.Deadline},
			{Label: "Attempts", Value: fmt.Sprintf("%d", t.Attempts)},
			{Label: "Created At", Value: t.CreatedAt},
		}
		if t.FulfilledByJobID != nil {
			fields = append(fields, ui.DetailField{Label: "Fulfilled By Job", Value: fmt.Sprintf("%d", *t.FulfilledByJobID)})
		}
		if t.FulfilledAt != nil {
			fields = append(fields, ui.DetailField{Label: "Fulfilled At", Value: *t.FulfilledAt})
		}
		return fields
	},
	GetID:    func(data interface{}) int { return data.(cli.Task).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Task %d", data.(cli.Task).ID) },
	Actions:  []ui.ActionDef{},
}

func init() {
	Register(Entry{Label: "Tasks", Hint: "View scheduled task windows", Def: TaskDef})
}
