package entity

import (
	"fmt"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var CrPrototypeDef = &EntityDef{
	Name: "🧬 Credential Prototypes", CLIEntity: "crprototype", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "Name", Width: 30, Field: "name"},
		{Header: "Description", Width: 40, Field: "description"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.CrPrototype
		if err := c.List("crprototype", limit, offset, &items); err != nil { return nil, err }
		rows := make([]ui.TableRow, len(items))
		for i, p := range items {
			desc := p.Description; if len(desc) > 38 { desc = desc[:35] + "..." }
			rows[i] = ui.TableRow{ID: p.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", p.ID), "name": p.Name, "description": desc,
			}, FullData: p}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		p := data.(cli.CrPrototype)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", p.ID)}, {Label: "Name", Value: p.Name},
			{Label: "UUID", Value: p.UUID}, {Label: "Code", Value: p.Code},
			{Label: "Description", Value: p.Description}, {Label: "Version", Value: p.Version},
		}
	},
	GetID: func(data interface{}) int { return data.(cli.CrPrototype).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("CrPrototype: %s", data.(cli.CrPrototype).Name) },
	Actions: []ui.ActionDef{},
}
func init() { Register(Entry{Label: "CrPrototypes", Hint: "View credential prototypes", Def: CrPrototypeDef}) }
