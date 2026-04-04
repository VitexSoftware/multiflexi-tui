package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var CompanyDef = &EntityDef{
	Name:         "🏢 Companies",
	CLIEntity:    "company",
	DeleteAction: "remove",
	Limit:        10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 6, Field: "id"},
		{Header: "Name", Width: 30, Field: "name"},
		{Header: "IC", Width: 15, Field: "ic"},
		{Header: "Email", Width: 25, Field: "email"},
		{Header: "Status", Width: 10, Field: "status"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Company
		if err := c.List("company", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, co := range items {
			status := "Disabled"
			if co.Enabled == 1 {
				status = "Enabled"
			}
			rows[i] = ui.TableRow{
				ID: co.ID,
				Values: map[string]string{
					"id": fmt.Sprintf("%d", co.ID), "name": co.Name,
					"ic": co.IC, "email": co.Email, "status": status,
				},
				FullData: co,
			}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		co := data.(cli.Company)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", co.ID)},
			{Label: "Name", Value: co.Name},
			{Label: "IC", Value: co.IC},
			{Label: "Email", Value: co.Email},
			{Label: "Slug", Value: co.Slug},
			{Label: "Enabled", Value: fmt.Sprintf("%d", co.Enabled)},
			{Label: "Server", Value: fmt.Sprintf("%d", co.Server)},
			{Label: "Created", Value: co.DatCreate},
			{Label: "Updated", Value: co.DatUpdate},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		co := data.(cli.Company)
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Company name", Value: co.Name},
			{Label: "Email", Placeholder: "Email", Value: co.Email},
			{Label: "IC", Placeholder: "IC", Value: co.IC},
			{Label: "Slug", Placeholder: "Slug", Value: co.Slug},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		co := data.(cli.Company)
		return []string{"--id", fmt.Sprintf("%d", co.ID),
			"--name", fields["Name"], "--email", fields["Email"],
			"--ic", fields["IC"], "--slug", fields["Slug"]}
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Name", Placeholder: "Company name", Required: true},
			{Label: "Email", Placeholder: "Email"},
			{Label: "IC", Placeholder: "IC"},
			{Label: "Slug", Placeholder: "Slug"},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		args := []string{"--name", fields["Name"]}
		if v := fields["Email"]; v != "" {
			args = append(args, "--email", v)
		}
		if v := fields["IC"]; v != "" {
			args = append(args, "--ic", v)
		}
		if v := fields["Slug"]; v != "" {
			args = append(args, "--slug", v)
		}
		return args
	},
	GetID:    func(data interface{}) int { return data.(cli.Company).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Company: %s", data.(cli.Company).Name) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() {
	Register(Entry{Label: "Companies", Hint: "View and manage companies", Def: CompanyDef})
}
