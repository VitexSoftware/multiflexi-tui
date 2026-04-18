package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var CredTypeDef = &EntityDef{
	Name: "🏷️ Credential Types", CLIEntity: "credtype", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 6, Field: "id"}, {Header: "Name", Width: 30, Field: "name"},
		{Header: "Class", Width: 35, Field: "class"}, {Header: "Version", Width: 8, Field: "version"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.CredType
		if err := c.List("credtype", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, t := range items {
			rows[i] = ui.TableRow{ID: t.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", t.ID), "name": t.Name,
				"class": t.Class, "version": fmt.Sprintf("%d", t.Version),
			}, FullData: t}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		t := data.(cli.CredType)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", t.ID)},
			{Label: "UUID", Value: t.UUID},
			{Label: "Name", Value: t.Name},
			{Label: "Class", Value: t.Class},
			{Label: "Company ID", Value: fmt.Sprintf("%d", t.CompanyID)},
			{Label: "URL", Value: t.URL},
			{Label: "Version", Value: fmt.Sprintf("%d", t.Version)},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		t := data.(cli.CredType)
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Credential type name", Value: t.Name},
			{Label: "Class", Placeholder: "PHP class name", Value: t.Class},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		t := data.(cli.CredType)
		return []string{"--id", fmt.Sprintf("%d", t.ID), "--name", fields["Name"], "--class", fields["Class"]}
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Credential type name", Required: true},
			{Label: "Company ID", Placeholder: "Company ID", Required: true},
			{Label: "Class", Placeholder: "PHP class name", Required: true},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		return []string{
			"--name", fields["Name"],
			"--company-id", fields["Company ID"],
			"--class", fields["Class"],
		}
	},
	GetID:    func(data interface{}) int { return data.(cli.CredType).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("CredType: %s", data.(cli.CredType).Name) },
	Actions:  []ui.ActionDef{{Label: "Edit", Key: "e", Command: "edit"}},
}

func init() { Register(Entry{Label: "CredTypes", Hint: "Manage credential types", Def: CredTypeDef}) }
