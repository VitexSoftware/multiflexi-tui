package entity

import "testing"

func TestRegistryPopulated(t *testing.T) {
	// All entity init() functions should have run
	if len(All) == 0 {
		t.Fatal("entity registry is empty; init() functions did not run")
	}
	// Expect at least 15 entities (13 original + 2 new)
	if len(All) < 14 {
		t.Errorf("expected at least 15 entities, got %d", len(All))
	}

	// Check each entry has required fields
	for _, e := range All {
		if e.Label == "" {
			t.Error("entity entry has empty label")
		}
		if e.Def == nil {
			t.Errorf("entity %s has nil def", e.Label)
			continue
		}
		if e.Def.CLIEntity == "" {
			t.Errorf("entity %s has empty CLIEntity", e.Label)
		}
		if e.Def.Fetch == nil {
			t.Errorf("entity %s has nil Fetch", e.Label)
		}
		if len(e.Def.Columns) == 0 {
			t.Errorf("entity %s has no columns", e.Label)
		}
	}
}
