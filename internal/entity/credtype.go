package entity

import (
	"fmt"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var CredTypeDef = &EntityDef{
	Name: "🏷️ Credential Types", CLIEntity: "credtype", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "UUID", Width: 38, Field: "uuid"},
		{Header: "Name", Width: 30, Field: "name"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.CredType
		if err := c.List("credtype", limit, offset, &items); err != nil { return nil, err }
		rows := make([]ui.TableRow, len(items))
		for i, t := range items {
			rows[i] = ui.TableRow{ID: t.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", t.ID), "uuid": t.UUID, "name": t.Name,
			}, FullData: t}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		t := data.(cli.CredType)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", t.ID)}, {Label: "UUID", Value: t.UUID},
			{Label: "Name", Value: t.Name}, {Label: "Class", Value: t.Class},
			{Label: "URL", Value: t.URL}, {Label: "Version", Value: fmt.Sprintf("%d", t.Version)},
		}
	},
	GetID: func(data interface{}) int { return data.(cli.CredType).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("CredType: %s", data.(cli.CredType).Name) },
	Actions: []ui.ActionDef{},
}
func init() { Register(Entry{Label: "CredTypes", Hint: "View credential types", Def: CredTypeDef}) }
