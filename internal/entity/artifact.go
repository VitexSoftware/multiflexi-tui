package entity

import (
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var ArtifactDef = &EntityDef{
	Name: "📎 Artifacts", CLIEntity: "artifact", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 8, Field: "id"}, {Header: "Job ID", Width: 10, Field: "job_id"},
		{Header: "Content Type", Width: 20, Field: "content_type"}, {Header: "Filename", Width: 35, Field: "filename"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.Artifact
		if err := c.List("artifact", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, a := range items {
			rows[i] = ui.TableRow{ID: a.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", a.ID), "job_id": fmt.Sprintf("%d", a.JobID),
				"content_type": a.ContentType, "filename": a.Filename,
			}, FullData: a}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		a := data.(cli.Artifact)
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", a.ID)},
			{Label: "Job ID", Value: fmt.Sprintf("%d", a.JobID)},
			{Label: "Filename", Value: a.Filename},
			{Label: "Content Type", Value: a.ContentType},
			{Label: "Created At", Value: a.CreatedAt},
			{Label: "Note", Value: a.Note},
		}
	},
	GetID:    func(data interface{}) int { return data.(cli.Artifact).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("Artifact: %s", data.(cli.Artifact).Filename) },
	Actions: []ui.ActionDef{
		{
			Label:   "Save to file",
			Key:     "s",
			Command: "save",
			Handler: func(c cli.Client, data interface{}) tea.Cmd {
				a := data.(cli.Artifact)
				form := NewActionFormView(
					fmt.Sprintf("Save Artifact: %s", a.Filename),
					[]ui.EditorField{
						{Label: "File Path", Placeholder: "/path/to/save", Value: "/tmp/" + a.Filename},
					},
					func(fields map[string]string) tea.Cmd {
						return func() tea.Msg {
							_, err := c.RunRaw("artifact", "save",
								"--id", fmt.Sprintf("%d", a.ID),
								"--file", fields["File Path"],
							)
							if err != nil {
								return ui.StatusMsg{Text: fmt.Sprintf("Save failed: %v", err)}
							}
							return ui.StatusMsg{Text: fmt.Sprintf("Saved to %s", fields["File Path"])}
						}
					},
				)
				return func() tea.Msg { return ui.NavigateToMsg{View: form} }
			},
		},
	},
}

func init() { Register(Entry{Label: "Artifacts", Hint: "View artifacts • s: save to file", Def: ArtifactDef}) }
