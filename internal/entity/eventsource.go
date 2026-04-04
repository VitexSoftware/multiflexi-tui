package entity

import (
	"fmt"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var EventSourceDef = &EntityDef{
	Name: "📡 Event Sources", CLIEntity: "eventsource", DeleteAction: "remove", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 5, Field: "id"}, {Header: "Name", Width: 25, Field: "name"},
		{Header: "Adapter", Width: 25, Field: "adapter"}, {Header: "Enabled", Width: 8, Field: "enabled"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.EventSource
		if err := c.List("eventsource", limit, offset, &items); err != nil { return nil, err }
		rows := make([]ui.TableRow, len(items))
		for i, es := range items {
			en := "No"; if es.Enabled == 1 { en = "Yes" }
			rows[i] = ui.TableRow{ID: es.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", es.ID), "name": es.Name, "adapter": es.AdapterType, "enabled": en,
			}, FullData: es}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		es := data.(cli.EventSource)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", es.ID)}, {Label: "Name", Value: es.Name},
			{Label: "Adapter", Value: es.AdapterType}, {Label: "DB Host", Value: es.DbHost},
			{Label: "DB Database", Value: es.DbDatabase}, {Label: "Enabled", Value: fmt.Sprintf("%d", es.Enabled)},
		}
	},
	GetID: func(data interface{}) int { return data.(cli.EventSource).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("EventSource: %s", data.(cli.EventSource).Name) },
	Actions: []ui.ActionDef{{Label: "Delete", Key: "d", Command: "delete"}},
}
func init() { Register(Entry{Label: "EventSources", Hint: "View event sources", Def: EventSourceDef}) }
