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
		if err := c.List("credential", limit, offset, &items); err != nil { return nil, err }
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
			{Label: "ID", Value: fmt.Sprintf("%d", cr.ID)}, {Label: "Name", Value: cr.Name},
			{Label: "Company ID", Value: fmt.Sprintf("%d", cr.CompanyID)},
			{Label: "Credential Type ID", Value: fmt.Sprintf("%d", cr.CredentialTypeID)},
		}
	},
	GetID: func(data interface{}) int { return data.(cli.Credential).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Credential: %s", data.(cli.Credential).Name) },
	Actions: []ui.ActionDef{{Label: "Delete", Key: "d", Command: "delete"}},
}
func init() { Register(Entry{Label: "Credentials", Hint: "View credentials", Def: CredentialDef}) }
