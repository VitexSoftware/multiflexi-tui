package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var CompanyAppDef = &EntityDef{
	Name: "🔗 Company-App Relations", CLIEntity: "companyapp", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "Company ID", Width: 12, Field: "company_id"},
		{Header: "App ID", Width: 12, Field: "app_id"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		// companyapp list requires --company_id and --app_id filters;
		// return an empty list so the TUI shows the assign/unassign actions.
		return []ui.TableRow{}, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		ca := data.(cli.CompanyApp)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", ca.ID)},
			{Label: "Company ID", Value: fmt.Sprintf("%d", ca.CompanyID)},
			{Label: "App ID", Value: fmt.Sprintf("%d", ca.AppID)},
		}
	},
	GetID:    func(data interface{}) int { return data.(cli.CompanyApp).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("CompanyApp %d", data.(cli.CompanyApp).ID) },
	Actions:  []ui.ActionDef{},
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
									"companyapp", "assign", "--format=json",
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
									"companyapp", "unassign", "--format=json",
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
