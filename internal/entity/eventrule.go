package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var EventRuleDef = &EntityDef{
	Name: "📌 Event Rules", CLIEntity: "eventrule", DeleteAction: "remove", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 5, Field: "id"}, {Header: "Source", Width: 8, Field: "source"},
		{Header: "Evidence", Width: 20, Field: "evidence"}, {Header: "Operation", Width: 10, Field: "op"},
		{Header: "Template", Width: 8, Field: "template"}, {Header: "Enabled", Width: 8, Field: "enabled"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.EventRule
		if err := c.List("eventrule", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, er := range items {
			en := "No"
			if er.Enabled == 1 {
				en = "Yes"
			}
			rows[i] = ui.TableRow{ID: er.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", er.ID), "source": fmt.Sprintf("%d", er.EventSourceID),
				"evidence": er.Evidence, "op": er.Operation,
				"template": fmt.Sprintf("%d", er.RunTemplateID), "enabled": en,
			}, FullData: er}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		er := data.(cli.EventRule)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", er.ID)},
			{Label: "Event Source ID", Value: fmt.Sprintf("%d", er.EventSourceID)},
			{Label: "Evidence", Value: er.Evidence},
			{Label: "Operation", Value: er.Operation},
			{Label: "RunTemplate ID", Value: fmt.Sprintf("%d", er.RunTemplateID)},
			{Label: "Priority", Value: fmt.Sprintf("%d", er.Priority)},
			{Label: "Enabled", Value: fmt.Sprintf("%d", er.Enabled)},
			{Label: "Env Mapping", Value: er.EnvMapping},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		er := data.(cli.EventRule)
		return []ui.EditorField{
			{Label: "Evidence", Placeholder: "e.g. faktura-vydana", Value: er.Evidence},
			{Label: "Operation", Placeholder: "any|create|update|delete", Value: er.Operation},
			{Label: "RunTemplate ID", Placeholder: "RunTemplate ID", Value: fmt.Sprintf("%d", er.RunTemplateID)},
			{Label: "Priority", Placeholder: "0", Value: fmt.Sprintf("%d", er.Priority)},
			{Label: "Enabled", Placeholder: "1 or 0", Value: fmt.Sprintf("%d", er.Enabled)},
			{Label: "Env Mapping", Placeholder: `{"KEY":"value"}`, Value: er.EnvMapping},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		er := data.(cli.EventRule)
		args := []string{
			"--id", fmt.Sprintf("%d", er.ID),
			"--evidence", fields["Evidence"],
			"--operation", fields["Operation"],
			"--runtemplate_id", fields["RunTemplate ID"],
			"--priority", fields["Priority"],
			"--enabled", fields["Enabled"],
		}
		if v := fields["Env Mapping"]; v != "" {
			args = append(args, "--env_mapping", v)
		}
		return args
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Event Source ID", Placeholder: "Event Source ID", Required: true},
			{Label: "Evidence", Placeholder: "e.g. faktura-vydana", Required: true},
			{Label: "RunTemplate ID", Placeholder: "RunTemplate ID", Required: true},
			{Label: "Operation", Placeholder: "any", Value: "any"},
			{Label: "Priority", Placeholder: "0", Value: "0"},
			{Label: "Enabled", Placeholder: "1 or 0", Value: "1"},
			{Label: "Env Mapping", Placeholder: `{"KEY":"value"}`},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		args := []string{
			"--event_source_id", fields["Event Source ID"],
			"--evidence", fields["Evidence"],
			"--runtemplate_id", fields["RunTemplate ID"],
			"--operation", fields["Operation"],
			"--priority", fields["Priority"],
			"--enabled", fields["Enabled"],
		}
		if v := fields["Env Mapping"]; v != "" {
			args = append(args, "--env_mapping", v)
		}
		return args
	},
	GetID:    func(data interface{}) int { return data.(cli.EventRule).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("EventRule %d", data.(cli.EventRule).ID) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() { Register(Entry{Label: "EventRules", Hint: "Manage event rules", Def: EventRuleDef}) }
