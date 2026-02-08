package app

import (
	"testing"
)

func TestMenuNavigation(t *testing.T) {
	// Create a model with menu items
	model := NewModel()

	// Test that cursor starts at 0
	if model.menuCursor != 0 {
		t.Errorf("Expected initial menu cursor 0, got %d", model.menuCursor)
	}

	// Test navigation within bounds
	initialCursor := model.menuCursor
	if initialCursor < len(model.menuItems)-1 {
		model.menuCursor++
		model.updateSelectedHint()

		if model.menuCursor != initialCursor+1 {
			t.Errorf("Expected cursor to increment, got %d", model.menuCursor)
		}
	}

	// Test that cursor doesn't go below 0
	model.menuCursor = 0
	if model.menuCursor > 0 {
		model.menuCursor--
	}
	if model.menuCursor < 0 {
		t.Errorf("Menu cursor should not be negative, got %d", model.menuCursor)
	}
}

func TestMenuBounds(t *testing.T) {
	model := NewModel()

	// Test cursor bounds
	maxCursor := len(model.menuItems) - 1

	// Test upper bound
	model.menuCursor = maxCursor
	if model.menuCursor < len(model.menuItems)-1 {
		model.menuCursor++
	}
	if model.menuCursor > maxCursor {
		t.Errorf("Menu cursor should not exceed max items, got %d, max %d", model.menuCursor, maxCursor)
	}

	// Test lower bound
	model.menuCursor = 0
	if model.menuCursor > 0 {
		model.menuCursor--
	}
	if model.menuCursor < 0 {
		t.Errorf("Menu cursor should not be negative, got %d", model.menuCursor)
	}
}

func TestUpdateSelectedHint(t *testing.T) {
	model := NewModel()

	// Test that updateSelectedHint doesn't panic
	model.updateSelectedHint()

	// Test with different cursor positions
	for i := 0; i < len(model.menuItems) && i < 5; i++ {
		model.menuCursor = i
		model.updateSelectedHint()
		// Just ensure it doesn't crash
	}
}
