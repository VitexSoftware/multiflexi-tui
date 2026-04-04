package entity

import (
	"fmt"
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

var ArtifactDef = &EntityDef{
	Name: "📎 Artifacts", CLIEntity: "artifact", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "Job ID", Width: 10, Field: "job_id"},
		{Header: "Filename", Width: 40, Field: "filename"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Artifact
		if err := c.List("artifact", limit, offset, &items); err != nil { return nil, err }
		rows := make([]ui.TableRow, len(items))
		for i, a := range items {
			rows[i] = ui.TableRow{ID: a.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", a.ID), "job_id": fmt.Sprintf("%d", a.JobID), "filename": a.Filename,
			}, FullData: a}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		a := data.(cli.Artifact)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", a.ID)}, {Label: "Job ID", Value: fmt.Sprintf("%d", a.JobID)},
			{Label: "Filename", Value: a.Filename}, {Label: "Content Type", Value: a.ContentType},
			{Label: "Created At", Value: a.CreatedAt}, {Label: "Note", Value: a.Note},
		}
	},
	GetID: func(data interface{}) int { return data.(cli.Artifact).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Artifact: %s", data.(cli.Artifact).Filename) },
	Actions: []ui.ActionDef{},
}
func init() { Register(Entry{Label: "Artifacts", Hint: "View artifacts", Def: ArtifactDef}) }
