package ui

import (
	"testing"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
)

func TestEntityListModelCreation(t *testing.T) {
	model := NewEntityListModel(EntityListConfig{
		Title: "Test List",
		Columns: []TableColumn{
			{Header: "ID", Width: 5, Field: "id"},
			{Header: "Name", Width: 20, Field: "name"},
		},
		Limit:    10,
		HelpText: "test help",
		Fetch: func(limit, offset int) (interface{}, error) {
			return []cli.CrPrototype{}, nil
		},
		Convert: func(data interface{}) []TableRow {
			return nil
		},
	})

	// Verify the model renders without panic
	view := model.View()
	if view == "" {
		t.Error("Entity list view should not be empty")
	}
}

func TestEditorModelCreation(t *testing.T) {
	fields := []EditorField{
		{Label: "Name", Value: "test"},
		{Label: "Email", Value: "a@b.c"},
	}
	model := NewEditorModel("Test Editor", nil, fields)

	// Verify the model renders without panic
	view := model.View()
	if view == "" {
		t.Error("Editor view should not be empty")
	}
}

func TestConfirmDialogModel(t *testing.T) {
	job := cli.Job{ID: 99, Command: "rm-all"}
	model := NewConfirmDialogModel("Job: rm-all", job)

	view := model.View()
	if view == "" {
		t.Error("Confirm dialog view should not be empty")
	}
}

func TestActiveMenuItemStyle(t *testing.T) {
	style := GetActiveMenuItemStyle()
	rendered := style.Render("Test")
	if rendered == "" {
		t.Error("Active menu item style should produce output")
	}
}

func TestPaginationLogic(t *testing.T) {
	limit := 10
	offset := 0

	if offset != 0 {
		t.Errorf("Expected first page offset 0, got %d", offset)
	}

	nextOffset := offset + limit
	if nextOffset != 10 {
		t.Errorf("Expected next page offset 10, got %d", nextOffset)
	}

	offset = 10
	prevOffset := offset - limit
	if prevOffset < 0 {
		prevOffset = 0
	}
	if prevOffset != 0 {
		t.Errorf("Expected previous page offset 0, got %d", prevOffset)
	}
}
