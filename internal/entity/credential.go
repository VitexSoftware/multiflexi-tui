package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var CredentialDef = &EntityDef{
	Name: "🔑 Credentials", CLIEntity: "credential", DeleteAction: "remove", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "Name", Width: 25, Field: "name"},
		{Header: "Company", Width: 12, Field: "company_id"}, {Header: "Type", Width: 12, Field: "type_id"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Credential
		if err := c.List("credential", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, cr := range items {
			rows[i] = ui.TableRow{ID: cr.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", cr.ID), "name": cr.Name,
				"company_id": fmt.Sprintf("%d", cr.CompanyID), "type_id": fmt.Sprintf("%d", cr.CredentialTypeID),
			}, FullData: cr}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		cr := data.(cli.Credential)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", cr.ID)},
			{Label: "Name", Value: cr.Name},
			{Label: "Company ID", Value: fmt.Sprintf("%d", cr.CompanyID)},
			{Label: "Credential Type ID", Value: fmt.Sprintf("%d", cr.CredentialTypeID)},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		cr := data.(cli.Credential)
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Credential name", Value: cr.Name},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		cr := data.(cli.Credential)
		return []string{"--id", fmt.Sprintf("%d", cr.ID), "--name", fields["Name"]}
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Credential name", Required: true},
			{Label: "Company ID", Placeholder: "Company ID", Required: true},
			{Label: "CredType ID", Placeholder: "Credential Type ID", Required: true},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		return []string{
			"--name", fields["Name"],
			"--company-id", fields["Company ID"],
			"--credential-type-id", fields["CredType ID"],
		}
	},
	GetID:    func(data interface{}) int { return data.(cli.Credential).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Credential: %s", data.(cli.Credential).Name) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() { Register(Entry{Label: "Credentials", Hint: "Manage credentials", Def: CredentialDef}) }
