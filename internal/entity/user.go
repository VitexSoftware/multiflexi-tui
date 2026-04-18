package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var UserDef = &EntityDef{
	Name: "👤 Users", CLIEntity: "user", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 6, Field: "id"}, {Header: "Login", Width: 20, Field: "login"},
		{Header: "Name", Width: 25, Field: "name"}, {Header: "Email", Width: 30, Field: "email"},
		{Header: "Active", Width: 7, Field: "enabled"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.User
		if err := c.List("user", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, u := range items {
			enabled := "No"
			if u.Enabled == 1 {
				enabled = "Yes"
			}
			rows[i] = ui.TableRow{ID: u.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", u.ID), "login": u.Login,
				"name": u.Firstname + " " + u.Lastname, "email": u.Email, "enabled": enabled,
			}, FullData: u}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		u := data.(cli.User)
		lastIP := ""
		if u.LastLoginIP != nil {
			lastIP = *u.LastLoginIP
		}
		lastAt := ""
		if u.LastLoginAt != nil {
			lastAt = *u.LastLoginAt
		}
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", u.ID)},
			{Label: "Login", Value: u.Login},
			{Label: "First Name", Value: u.Firstname},
			{Label: "Last Name", Value: u.Lastname},
			{Label: "Email", Value: u.Email},
			{Label: "Enabled", Value: fmt.Sprintf("%d", u.Enabled)},
			{Label: "2FA Enabled", Value: fmt.Sprintf("%d", u.TwoFactorEnabled)},
			{Label: "Last Login IP", Value: lastIP},
			{Label: "Last Login At", Value: lastAt},
			{Label: "Failed Logins", Value: fmt.Sprintf("%d", u.FailedLoginAttempts)},
			{Label: "Created", Value: u.DatCreate},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		u := data.(cli.User)
		return []ui.EditorField{
			{Label: "Login", Placeholder: "username", Value: u.Login},
			{Label: "First Name", Placeholder: "First name", Value: u.Firstname},
			{Label: "Last Name", Placeholder: "Last name", Value: u.Lastname},
			{Label: "Email", Placeholder: "email@example.com", Value: u.Email},
			{Label: "Enabled", Placeholder: "1 or 0", Value: fmt.Sprintf("%d", u.Enabled)},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		u := data.(cli.User)
		return []string{
			"--id", fmt.Sprintf("%d", u.ID),
			"--login", fields["Login"],
			"--firstname", fields["First Name"],
			"--lastname", fields["Last Name"],
			"--email", fields["Email"],
			"--enabled", fields["Enabled"],
		}
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Login", Placeholder: "username", Required: true},
			{Label: "Email", Placeholder: "email@example.com", Required: true},
			{Label: "Password", Placeholder: "plaintext password", Required: true},
			{Label: "First Name", Placeholder: "First name"},
			{Label: "Last Name", Placeholder: "Last name"},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		args := []string{
			"--login", fields["Login"],
			"--email", fields["Email"],
			"--plaintext", fields["Password"],
		}
		if v := fields["First Name"]; v != "" {
			args = append(args, "--firstname", v)
		}
		if v := fields["Last Name"]; v != "" {
			args = append(args, "--lastname", v)
		}
		return args
	},
	GetID:    func(data interface{}) int { return data.(cli.User).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("User: %s", data.(cli.User).Login) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() { Register(Entry{Label: "Users", Hint: "Manage users", Def: UserDef}) }
