package entity

import (
	"encoding/json"
	"fmt"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

const confFieldTypes = "string|text|integer|float|email|url|password|bool|file-path|set"

// NewConfFieldDef returns an EntityDef scoped to one application's config fields
// (same data managed by multiflexi-web5 conffield.php?app_id=N).
func NewConfFieldDef(appID int, appName string) *EntityDef {
	appIDStr := fmt.Sprintf("%d", appID)
	title := fmt.Sprintf("⚙️ Config Fields — %s", appName)

	return &EntityDef{
		Name: title, CLIEntity: "conffield", DeleteAction: "delete", Limit: 10,
		Columns: []ui.TableColumn{
			{Header: "ID", Width: 6, Field: "id"},
			{Header: "Keyword", Width: 28, Field: "keyname", Flex: true},
			{Header: "Type", Width: 10, Field: "type"},
			{Header: "Required", Width: 8, Field: "required"},
			{Header: "Secret", Width: 6, Field: "secret"},
			{Header: "Default", Width: 20, Field: "defval"},
		},
		Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
			out, err := c.RunRaw(
				"conffield:list",
				"--format=json",
				"--app_id", appIDStr,
				"--order=D",
				fmt.Sprintf("--limit=%d", limit),
				fmt.Sprintf("--offset=%d", offset),
			)
			if err != nil {
				return nil, err
			}
			var items []cli.ConfField
			if err := json.Unmarshal(out, &items); err != nil {
				return nil, fmt.Errorf("parse conffield JSON: %w", err)
			}
			rows := make([]ui.TableRow, len(items))
			for i, f := range items {
				req, sec := "No", "No"
				if f.Required == 1 {
					req = "Yes"
				}
				if f.Secret == 1 {
					sec = "Yes"
				}
				defval := f.Defval
				if f.Secret == 1 && defval != "" {
					defval = "********"
				}
				rows[i] = ui.TableRow{ID: f.ID, Values: map[string]string{
					"id": fmt.Sprintf("%d", f.ID), "keyname": f.Keyname, "type": f.Type,
					"required": req, "secret": sec, "defval": defval,
				}, FullData: f}
			}
			return rows, nil
		},
		ToDetail: func(data interface{}) []ui.DetailField {
			f := data.(cli.ConfField)
			return []ui.DetailField{
				{Label: "ID", Value: fmt.Sprintf("%d", f.ID)},
				{Label: "App ID", Value: fmt.Sprintf("%d", f.AppID)},
				{Label: "Keyword", Value: f.Keyname},
				{Label: "Type", Value: f.Type},
				{Label: "Description", Value: f.Description},
				{Label: "Hint", Value: f.Hint},
				{Label: "Note", Value: f.Note},
				{Label: "Default", Value: f.Defval},
				{Label: "Required", Value: fmt.Sprintf("%d", f.Required)},
				{Label: "Secret", Value: fmt.Sprintf("%d", f.Secret)},
				{Label: "Multiline", Value: fmt.Sprintf("%d", f.Multiline)},
				{Label: "Expiring", Value: fmt.Sprintf("%d", f.Expiring)},
			}
		},
		ToEditor: func(data interface{}) []ui.EditorField {
			f := data.(cli.ConfField)
			return []ui.EditorField{
				{Label: "Keyword", Placeholder: "ENV_VAR_NAME", Value: f.Keyname},
				{Label: "Type", Placeholder: confFieldTypes, Value: f.Type},
				{Label: "Description", Placeholder: "Description", Value: f.Description},
				{Label: "Hint", Placeholder: "Hint for users", Value: f.Hint},
				{Label: "Note", Placeholder: "Internal note", Value: f.Note},
				{Label: "Default", Placeholder: "Default value", Value: f.Defval},
				{Label: "Required", Placeholder: "1 or 0", Value: fmt.Sprintf("%d", f.Required)},
				{Label: "Secret", Placeholder: "1 or 0", Value: fmt.Sprintf("%d", f.Secret)},
				{Label: "Multiline", Placeholder: "1 or 0", Value: fmt.Sprintf("%d", f.Multiline)},
				{Label: "Expiring", Placeholder: "1 or 0", Value: fmt.Sprintf("%d", f.Expiring)},
			}
		},
		UpdateArgs: func(data interface{}, fields map[string]string) []string {
			f := data.(cli.ConfField)
			return []string{
				"--id", fmt.Sprintf("%d", f.ID),
				"--keyname", fields["Keyword"],
				"--type", fields["Type"],
				"--description", fields["Description"],
				"--hint", fields["Hint"],
				"--note", fields["Note"],
				"--defval", fields["Default"],
				"--required", fields["Required"],
				"--secret", fields["Secret"],
				"--multiline", fields["Multiline"],
				"--expiring", fields["Expiring"],
			}
		},
		NewFields: func() []ui.EditorField {
			return []ui.EditorField{
				{Label: "Keyword", Placeholder: "ENV_VAR_NAME", Required: true},
				{Label: "Type", Placeholder: confFieldTypes, Value: "string", Required: true},
				{Label: "Description", Placeholder: "Description"},
				{Label: "Hint", Placeholder: "Hint for users"},
				{Label: "Note", Placeholder: "Internal note"},
				{Label: "Default", Placeholder: "Default value"},
				{Label: "Required", Placeholder: "1 or 0", Value: "0"},
				{Label: "Secret", Placeholder: "1 or 0", Value: "0"},
				{Label: "Multiline", Placeholder: "1 or 0", Value: "0"},
				{Label: "Expiring", Placeholder: "1 or 0", Value: "0"},
			}
		},
		CreateArgs: func(fields map[string]string) []string {
			args := []string{
				"--app_id", appIDStr,
				"--keyname", fields["Keyword"],
				"--type", fields["Type"],
			}
			for _, pair := range [][2]string{
				{"--description", "Description"},
				{"--hint", "Hint"},
				{"--note", "Note"},
				{"--defval", "Default"},
				{"--required", "Required"},
				{"--secret", "Secret"},
				{"--multiline", "Multiline"},
				{"--expiring", "Expiring"},
			} {
				if v := fields[pair[1]]; v != "" {
					args = append(args, pair[0], v)
				}
			}
			return args
		},
		GetID:    func(data interface{}) int { return data.(cli.ConfField).ID },
		GetLabel: func(data interface{}) string { return fmt.Sprintf("Field: %s", data.(cli.ConfField).Keyname) },
		Actions: []ui.ActionDef{
			{Label: "Edit", Key: "e", Command: "edit"},
			{Label: "Delete", Key: "d", Command: "delete"},
		},
	}
}
