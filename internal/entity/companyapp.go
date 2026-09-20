package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var CompanyAppDef = &EntityDef{
	Name: "🔗 Company-App Relations", CLIEntity: "company-app", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 6, Field: "id"},
		{Header: "Company", Width: 25, Field: "company"},
		{Header: "App", Width: 30, Field: "app"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.CompanyApp
		if err := c.List("company-app", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, ca := range items {
			rows[i] = ui.TableRow{ID: ca.ID, Values: map[string]string{
				"id":      fmt.Sprintf("%d", ca.ID),
				"company": ca.CompanyName,
				"app":     ca.AppName,
			}, FullData: ca}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		ca := data.(cli.CompanyApp)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", ca.ID)},
			{Label: "Company ID", Value: fmt.Sprintf("%d", ca.CompanyID)},
			{Label: "Company", Value: ca.CompanyName},
			{Label: "Company Slug", Value: ca.CompanySlug},
			{Label: "App ID", Value: fmt.Sprintf("%d", ca.AppID)},
			{Label: "App", Value: ca.AppName},
			{Label: "App UUID", Value: ca.AppUUID},
		}
	},
	GetID:    func(data interface{}) int { return data.(cli.CompanyApp).ID },
	GetLabel: func(data interface{}) string {
		ca := data.(cli.CompanyApp)
		return fmt.Sprintf("%s → %s", ca.CompanyName, ca.AppName)
	},
	Actions: []ui.ActionDef{},
	ListActions: []ui.ListActionDef{
		{
			Label: "Assign",
			Key:   "a",
			Handler: func(c cli.Client) tea.Cmd {
				return func() tea.Msg {
					form := NewActionFormView(
						"Assign Application to Company",
						[]ui.EditorField{
							{Label: "Company ID", Placeholder: "Company ID (number)", Required: true},
							{Label: "App ID", Placeholder: "Application ID (number)", Required: true},
						},
						func(fields map[string]string) tea.Cmd {
							return func() tea.Msg {
								out, err := c.RunRaw(
									"company-app:assign", "--format=json",
									"--company_id", fields["Company ID"],
									"--app_id", fields["App ID"],
								)
								viewer := ui.NewViewer("Assign Result")
								viewer.RefreshOnBack = true
								if err != nil {
									viewer.SetContent("Assign Result", fmt.Sprintf("Error: %v\n\n%s", err, string(out)))
								} else {
									viewer.SetContent("Assign Result", string(out))
								}
								return ui.NavigateToMsg{View: viewer}
							}
						},
					)
					return ui.NavigateToMsg{View: form}
				}
			},
		},
		{
			Label: "Unassign",
			Key:   "u",
			Handler: func(c cli.Client) tea.Cmd {
				return func() tea.Msg {
					form := NewActionFormView(
						"Unassign Application from Company",
						[]ui.EditorField{
							{Label: "Company ID", Placeholder: "Company ID (number)", Required: true},
							{Label: "App ID", Placeholder: "Application ID (number)", Required: true},
						},
						func(fields map[string]string) tea.Cmd {
							return func() tea.Msg {
								out, err := c.RunRaw(
									"company-app:unassign", "--format=json",
									"--company_id", fields["Company ID"],
									"--app_id", fields["App ID"],
								)
								viewer := ui.NewViewer("Unassign Result")
								viewer.RefreshOnBack = true
								if err != nil {
									viewer.SetContent("Unassign Result", fmt.Sprintf("Error: %v\n\n%s", err, string(out)))
								} else {
									viewer.SetContent("Unassign Result", string(out))
								}
								return ui.NavigateToMsg{View: viewer}
							}
						},
					)
					return ui.NavigateToMsg{View: form}
				}
			},
		},
	},
}

func init() {
	Register(Entry{Label: "CompanyApps", Hint: "Manage company-app assignments • a: assign • u: unassign", Def: CompanyAppDef})
}
