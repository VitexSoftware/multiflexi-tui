package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var EventSourceDef = &EntityDef{
	Name: "📡 Event Sources", CLIEntity: "eventsource", DeleteAction: "remove", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 5, Field: "id"}, {Header: "Name", Width: 25, Field: "name"},
		{Header: "Adapter", Width: 30, Field: "adapter"}, {Header: "Enabled", Width: 8, Field: "enabled"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.EventSource
		if err := c.List("eventsource", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, es := range items {
			en := "No"
			if es.Enabled == 1 {
				en = "Yes"
			}
			rows[i] = ui.TableRow{ID: es.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", es.ID), "name": es.Name, "adapter": es.AdapterType, "enabled": en,
			}, FullData: es}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		es := data.(cli.EventSource)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", es.ID)},
			{Label: "Name", Value: es.Name},
			{Label: "Adapter Type", Value: es.AdapterType},
			{Label: "DB Connection", Value: es.DbConnection},
			{Label: "DB Host", Value: fmt.Sprintf("%s:%s", es.DbHost, es.DbPort)},
			{Label: "DB Database", Value: es.DbDatabase},
			{Label: "DB Username", Value: es.DbUsername},
			{Label: "Poll Interval", Value: fmt.Sprintf("%d s", es.PollInterval)},
			{Label: "Enabled", Value: fmt.Sprintf("%d", es.Enabled)},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		es := data.(cli.EventSource)
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Source name", Value: es.Name},
			{Label: "Adapter Type", Placeholder: "abraflexi-webhook-acceptor", Value: es.AdapterType},
			{Label: "DB Connection", Placeholder: "mysql|pgsql|sqlite", Value: es.DbConnection},
			{Label: "DB Host", Placeholder: "localhost", Value: es.DbHost},
			{Label: "DB Port", Placeholder: "3306", Value: es.DbPort},
			{Label: "DB Database", Placeholder: "database name", Value: es.DbDatabase},
			{Label: "DB Username", Placeholder: "username", Value: es.DbUsername},
			{Label: "DB Password", Placeholder: "password", Value: es.DbPassword},
			{Label: "Poll Interval", Placeholder: "60", Value: fmt.Sprintf("%d", es.PollInterval)},
			{Label: "Enabled", Placeholder: "1 or 0", Value: fmt.Sprintf("%d", es.Enabled)},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		es := data.(cli.EventSource)
		return []string{
			"--id", fmt.Sprintf("%d", es.ID),
			"--name", fields["Name"],
			"--adapter_type", fields["Adapter Type"],
			"--db_connection", fields["DB Connection"],
			"--db_host", fields["DB Host"],
			"--db_port", fields["DB Port"],
			"--db_database", fields["DB Database"],
			"--db_username", fields["DB Username"],
			"--db_password", fields["DB Password"],
			"--poll_interval", fields["Poll Interval"],
			"--enabled", fields["Enabled"],
		}
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Source name", Required: true},
			{Label: "Adapter Type", Placeholder: "abraflexi-webhook-acceptor", Required: true},
			{Label: "DB Connection", Placeholder: "mysql", Value: "mysql"},
			{Label: "DB Host", Placeholder: "localhost", Value: "localhost"},
			{Label: "DB Port", Placeholder: "3306", Value: "3306"},
			{Label: "DB Database", Placeholder: "database name"},
			{Label: "DB Username", Placeholder: "username"},
			{Label: "DB Password", Placeholder: "password"},
			{Label: "Poll Interval", Placeholder: "60", Value: "60"},
			{Label: "Enabled", Placeholder: "1 or 0", Value: "1"},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		args := []string{
			"--name", fields["Name"],
			"--adapter_type", fields["Adapter Type"],
			"--enabled", fields["Enabled"],
			"--poll_interval", fields["Poll Interval"],
		}
		if v := fields["DB Connection"]; v != "" {
			args = append(args, "--db_connection", v)
		}
		if v := fields["DB Host"]; v != "" {
			args = append(args, "--db_host", v)
		}
		if v := fields["DB Port"]; v != "" {
			args = append(args, "--db_port", v)
		}
		if v := fields["DB Database"]; v != "" {
			args = append(args, "--db_database", v)
		}
		if v := fields["DB Username"]; v != "" {
			args = append(args, "--db_username", v)
		}
		if v := fields["DB Password"]; v != "" {
			args = append(args, "--db_password", v)
		}
		return args
	},
	GetID:    func(data interface{}) int { return data.(cli.EventSource).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("EventSource: %s", data.(cli.EventSource).Name) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{
			Label:   "Test",
			Key:     "t",
			Command: "test",
			Handler: func(c cli.Client, data interface{}) tea.Cmd {
				es := data.(cli.EventSource)
				return func() tea.Msg {
					output, err := c.RunRaw("eventsource", "test", "--format=json",
						"--id", fmt.Sprintf("%d", es.ID))
					viewer := ui.NewViewer(fmt.Sprintf("Test: %s", es.Name))
					if err != nil {
						viewer.SetContent(fmt.Sprintf("Test: %s", es.Name), fmt.Sprintf("Error: %v\n\n%s", err, string(output)))
					} else {
						viewer.SetContent(fmt.Sprintf("Test: %s", es.Name), string(output))
					}
					return ui.NavigateToMsg{View: viewer}
				}
			},
		},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() {
	Register(Entry{Label: "EventSources", Hint: "Manage event sources", Def: EventSourceDef})
}
