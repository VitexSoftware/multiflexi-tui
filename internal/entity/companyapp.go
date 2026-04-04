package entity

import (
	"fmt"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var CompanyAppDef = &EntityDef{
	Name: "🔗 Company-App Relations", CLIEntity: "companyapp", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "Company", Width: 15, Field: "company_id"},
		{Header: "App", Width: 15, Field: "app_id"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.CompanyApp
		if err := c.List("companyapp", limit, offset, &items); err != nil { return nil, err }
		rows := make([]ui.TableRow, len(items))
		for i, ca := range items {
			rows[i] = ui.TableRow{ID: ca.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", ca.ID), "company_id": fmt.Sprintf("%d", ca.CompanyID), "app_id": fmt.Sprintf("%d", ca.AppID),
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
	GetID: func(data interface{}) int { return data.(cli.CompanyApp).ID },
	GetLabel: func(data interface{}) string { ca := data.(cli.CompanyApp); return fmt.Sprintf("CompanyApp %d", ca.ID) },
	Actions: []ui.ActionDef{{Label: "Delete", Key: "d", Command: "delete"}},
}
func init() { Register(Entry{Label: "CompanyApps", Hint: "View company-app relations", Def: CompanyAppDef}) }
