package entity

import (
	"fmt"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var TokenDef = &EntityDef{
	Name: "🎟️ Tokens", CLIEntity: "token", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "User", Width: 25, Field: "user"},
		{Header: "Token", Width: 30, Field: "token"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Token
		if err := c.List("token", limit, offset, &items); err != nil { return nil, err }
		rows := make([]ui.TableRow, len(items))
		for i, t := range items {
			rows[i] = ui.TableRow{ID: t.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", t.ID), "user": t.User, "token": t.Token,
			}, FullData: t}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		t := data.(cli.Token)
		return []ui.DetailField{{Label: "ID", Value: fmt.Sprintf("%d", t.ID)}, {Label: "User", Value: t.User}, {Label: "Token", Value: t.Token}}
	},
	GetID: func(data interface{}) int { return data.(cli.Token).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Token %d", data.(cli.Token).ID) },
	Actions: []ui.ActionDef{{Label: "Delete", Key: "d", Command: "delete"}},
}
func init() { Register(Entry{Label: "Tokens", Hint: "View tokens", Def: TokenDef}) }
