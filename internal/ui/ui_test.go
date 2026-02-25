package ui

import (
	"testing"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
)

func TestCrPrototypesModelInit(t *testing.T) {
	model := NewCrPrototypesModel()

	if model.limit != 10 {
		t.Errorf("Expected default limit 10, got %d", model.limit)
	}
	if model.offset != 0 {
		t.Errorf("Expected default offset 0, got %d", model.offset)
	}
	if model.cursor != 0 {
		t.Errorf("Expected default cursor 0, got %d", model.cursor)
	}
	if model.loading != true {
		t.Errorf("Expected default loading true, got %t", model.loading)
	}
}

func TestCrPrototypesModelWithData(t *testing.T) {
	// Create test data
	testPrototypes := []cli.CrPrototype{
		{
			ID:      1,
			Name:    "Test Prototype 1",
			Version: "1.0.0",
		},
		{
			ID:      2,
			Name:    "Test Prototype 2",
			Version: "2.0.0",
		},
	}

	model := NewCrPrototypesModel()
	model.crprototypes = testPrototypes

	// Test that model holds the data correctly
	if len(model.crprototypes) != 2 {
		t.Errorf("Expected 2 prototypes, got %d", len(model.crprototypes))
	}
	if model.crprototypes[0].Name != "Test Prototype 1" {
		t.Errorf("Expected first prototype name 'Test Prototype 1', got '%s'", model.crprototypes[0].Name)
	}
	if model.crprototypes[1].Version != "2.0.0" {
		t.Errorf("Expected second prototype version '2.0.0', got '%s'", model.crprototypes[1].Version)
	}
}

func TestJobEditorModelInit(t *testing.T) {
	job := cli.Job{ID: 42, Command: "test-cmd", Executor: "Native", ScheduleType: "hourly"}
	model := NewJobEditorModel(job)

	if model.job.ID != 42 {
		t.Errorf("Expected job ID 42, got %d", model.job.ID)
	}
	if len(model.inputs) != 3 {
		t.Errorf("Expected 3 input fields, got %d", len(model.inputs))
	}
	if len(model.labels) != 3 {
		t.Errorf("Expected 3 labels, got %d", len(model.labels))
	}
	if model.cursor != 0 {
		t.Errorf("Expected cursor at 0, got %d", model.cursor)
	}
}

func TestCompanyEditorModelInit(t *testing.T) {
	company := cli.Company{ID: 7, Name: "Test Co", Email: "a@b.cz", IC: "123", Slug: "test-co"}
	model := NewCompanyEditorModel(company)

	if model.company.ID != 7 {
		t.Errorf("Expected company ID 7, got %d", model.company.ID)
	}
	if len(model.inputs) != 4 {
		t.Errorf("Expected 4 input fields, got %d", len(model.inputs))
	}
	if len(model.labels) != 4 {
		t.Errorf("Expected 4 labels, got %d", len(model.labels))
	}
}

func TestConfirmDialogModel(t *testing.T) {
	job := cli.Job{ID: 99, Command: "rm-all"}
	model := NewConfirmDialogModel("Job: rm-all", job)

	// Test View renders without panic
	view := model.View()
	if view == "" {
		t.Error("Confirm dialog view should not be empty")
	}
}

func TestActiveMenuItemStyle(t *testing.T) {
	style := GetActiveMenuItemStyle()
	// Verify it renders without panic
	rendered := style.Render("Test")
	if rendered == "" {
		t.Error("Active menu item style should produce output")
	}
}

func TestPaginationLogic(t *testing.T) {
	// Test pagination calculations
	limit := 10
	offset := 0

	// Test first page
	if offset != 0 {
		t.Errorf("Expected first page offset 0, got %d", offset)
	}

	// Test next page calculation
	nextOffset := offset + limit
	if nextOffset != 10 {
		t.Errorf("Expected next page offset 10, got %d", nextOffset)
	}

	// Test previous page calculation (from second page)
	offset = 10
	prevOffset := offset - limit
	if prevOffset < 0 {
		prevOffset = 0
	}
	if prevOffset != 0 {
		t.Errorf("Expected previous page offset 0, got %d", prevOffset)
	}
}
