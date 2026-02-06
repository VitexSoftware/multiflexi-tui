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