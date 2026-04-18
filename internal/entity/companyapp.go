package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var CompanyAppDef = &EntityDef{
	Name: "🔗 Company-App Relations", CLIEntity: "companyapp", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "Company ID", Width: 12, Field: "company_id"},
		{Header: "App ID", Width: 12, Field: "app_id"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.CompanyApp
		if err := c.List("companyapp", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, ca := range items {
			rows[i] = ui.TableRow{ID: ca.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", ca.ID), "company_id": fmt.Sprintf("%d", ca.CompanyID),
				"app_id": fmt.Sprintf("%d", ca.AppID),
			}, FullData: ca}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		ca := data.(cli.CompanyApp)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", ca.ID)},
			{Label: "Company ID", Value: fmt.Sprintf("%d", ca.CompanyID)},
			{Label: "App ID", Value: fmt.Sprintf("%d", ca.AppID)},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		ca := data.(cli.CompanyApp)
		return []ui.EditorField{
			{Label: "Company ID", Placeholder: "Company ID", Value: fmt.Sprintf("%d", ca.CompanyID)},
			{Label: "App ID", Placeholder: "Application ID", Value: fmt.Sprintf("%d", ca.AppID)},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		ca := data.(cli.CompanyApp)
		return []string{"--id", fmt.Sprintf("%d", ca.ID),
			"--company_id", fields["Company ID"], "--app_id", fields["App ID"]}
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Company ID", Placeholder: "Company ID", Required: true},
			{Label: "App ID", Placeholder: "Application ID", Required: true},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		return []string{"--company_id", fields["Company ID"], "--app_id", fields["App ID"]}
	},
	GetID:    func(data interface{}) int { return data.(cli.CompanyApp).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("CompanyApp %d", data.(cli.CompanyApp).ID) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() {
	Register(Entry{Label: "CompanyApps", Hint: "Manage company-app relations", Def: CompanyAppDef})
}
