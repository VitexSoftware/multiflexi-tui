package app

import (
	"testing"
)

func TestMenuOffsetLogic(t *testing.T) {
	// Create a model with many menu items
	model := NewModel()
	model.width = 60 // Simulate a narrow terminal

	// Test that offset starts at 0
	if model.menuOffset != 0 {
		t.Errorf("Expected initial menu offset 0, got %d", model.menuOffset)
	}

	// Test navigation to a position that should trigger scrolling
	model.menuCursor = 10 // Move to "CrPrototypes"
	model.updateMenuOffset()

	// The offset should have been adjusted to keep the cursor visible
	if model.menuOffset < 0 {
		t.Errorf("Menu offset should not be negative, got %d", model.menuOffset)
	}

	if model.menuOffset > model.menuCursor {
		t.Errorf("Menu offset (%d) should not be greater than cursor (%d)", model.menuOffset, model.menuCursor)
	}
}

func TestMenuOffsetBounds(t *testing.T) {
	model := NewModel()
	model.width = 40 // Very narrow terminal

	// Test with cursor at various positions
	testPositions := []int{0, 5, 10, 15, len(model.menuItems) - 1}

	for _, pos := range testPositions {
		if pos < len(model.menuItems) {
			model.menuCursor = pos
			model.updateMenuOffset()

			// Check bounds
			if model.menuOffset < 0 {
				t.Errorf("Menu offset should not be negative at cursor %d, got %d", pos, model.menuOffset)
			}

			if model.menuOffset >= len(model.menuItems) {
				t.Errorf("Menu offset should not exceed items count at cursor %d, got %d", pos, model.menuOffset)
			}

			// Cursor should always be >= offset
			if model.menuCursor < model.menuOffset {
				t.Errorf("Cursor (%d) should be >= offset (%d)", model.menuCursor, model.menuOffset)
			}
		}
	}
}
