package entity

import (
	"fmt"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var QueueDef = &EntityDef{
	Name: "📬 Queue", CLIEntity: "queue", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "Job", Width: 8, Field: "job"},
		{Header: "Type", Width: 12, Field: "type"}, {Header: "App", Width: 20, Field: "app"},
		{Header: "Company", Width: 20, Field: "company"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Queue
		if err := c.List("queue", limit, offset, &items); err != nil { return nil, err }
		rows := make([]ui.TableRow, len(items))
		for i, q := range items {
			rows[i] = ui.TableRow{ID: q.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", q.ID), "job": fmt.Sprintf("%d", q.Job),
				"type": q.ScheduleType, "app": q.AppName, "company": q.CompanyName,
			}, FullData: q}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		q := data.(cli.Queue)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", q.ID)}, {Label: "Job", Value: fmt.Sprintf("%d", q.Job)},
			{Label: "Schedule Type", Value: q.ScheduleType}, {Label: "App", Value: q.AppName},
			{Label: "Company", Value: q.CompanyName}, {Label: "After", Value: q.After},
		}
	},
	GetID: func(data interface{}) int { return data.(cli.Queue).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Queue %d", data.(cli.Queue).ID) },
	Actions: []ui.ActionDef{},
}
func init() { Register(Entry{Label: "Queue", Hint: "View job queue", Def: QueueDef}) }
