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
		{Header: "Name", Width: 20, Field: "name"}, {Header: "Email", Width: 30, Field: "email"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.User
		if err := c.List("user", limit, offset, &items); err != nil { return nil, err }
		rows := make([]ui.TableRow, len(items))
		for i, u := range items {
			rows[i] = ui.TableRow{ID: u.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", u.ID), "login": u.Login,
				"name": u.Firstname + " " + u.Lastname, "email": u.Email,
			}, FullData: u}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		u := data.(cli.User)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", u.ID)}, {Label: "Login", Value: u.Login},
			{Label: "First Name", Value: u.Firstname}, {Label: "Last Name", Value: u.Lastname},
			{Label: "Email", Value: u.Email}, {Label: "Enabled", Value: fmt.Sprintf("%d", u.Enabled)},
		}
	},
	GetID: func(data interface{}) int { return data.(cli.User).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("User: %s", data.(cli.User).Login) },
	Actions: []ui.ActionDef{{Label: "Delete", Key: "d", Command: "delete"}},
}
func init() { Register(Entry{Label: "Users", Hint: "View users", Def: UserDef}) }
