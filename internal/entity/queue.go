package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var QueueDef = &EntityDef{
	Name: "📬 Queue", CLIEntity: "queue", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "Job", Width: 8, Field: "job"},
		{Header: "Type", Width: 12, Field: "type"}, {Header: "After", Width: 20, Field: "after"},
		{Header: "App", Width: 20, Field: "app"}, {Header: "Company", Width: 20, Field: "company"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Queue
		if err := c.List("queue", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, q := range items {
			rows[i] = ui.TableRow{ID: q.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", q.ID), "job": fmt.Sprintf("%d", q.Job),
				"type": q.ScheduleType, "after": q.After, "app": q.AppName, "company": q.CompanyName,
			}, FullData: q}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		q := data.(cli.Queue)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", q.ID)},
			{Label: "Job", Value: fmt.Sprintf("%d", q.Job)},
			{Label: "Schedule Type", Value: q.ScheduleType},
			{Label: "RunTemplate", Value: fmt.Sprintf("%d — %s", q.RunTemplateID, q.RunTemplateName)},
			{Label: "App", Value: fmt.Sprintf("%d — %s", q.AppID, q.AppName)},
			{Label: "Company", Value: fmt.Sprintf("%d — %s", q.CompanyID, q.CompanyName)},
			{Label: "After", Value: q.After},
		}
	},
	GetID:    func(data interface{}) int { return data.(cli.Queue).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Queue %d", data.(cli.Queue).ID) },
	Actions:  []ui.ActionDef{},
	ListActions: []ui.ListActionDef{
		{
			Label:   "Fix",
			Key:     "f",
			Confirm: "Run queue fix? This repairs stuck queue entries.",
			Handler: func(c cli.Client) tea.Cmd {
				return func() tea.Msg {
					output, err := c.RunRaw("queue", "fix", "--format=json")
					if err != nil {
						return ui.StatusMsg{Text: fmt.Sprintf("Queue fix failed: %v", err)}
					}
					viewer := ui.NewViewer("Queue Fix Result")
					viewer.RefreshOnBack = true
					viewer.SetContent("Queue Fix Result", string(output))
					return ui.NavigateToMsg{View: viewer}
				}
			},
		},
		{
			Label:   "Truncate",
			Key:     "T",
			Confirm: "TRUNCATE entire queue? All pending jobs will be removed!",
			Handler: func(c cli.Client) tea.Cmd {
				return func() tea.Msg {
					_, err := c.RunRaw("queue", "truncate", "--format=json")
					if err != nil {
						return ui.StatusMsg{Text: fmt.Sprintf("Queue truncate failed: %v", err)}
					}
					return ui.RefreshCurrentMsg{Status: "Queue truncated"}
				}
			},
		},
	},
}

func init() {
	Register(Entry{Label: "Queue", Hint: "View job queue • f: fix • T: truncate", Def: QueueDef})
}
