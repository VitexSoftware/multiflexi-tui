package ui

import "testing"

func TestTableWidgetNavigation(t *testing.T) {
	tw := NewTableWidget("Test", []TableColumn{
		{Header: "ID", Width: 5, Field: "id"},
	}, 10, "help")

	rows := []TableRow{
		{ID: 1, Values: map[string]string{"id": "1"}},
		{ID: 2, Values: map[string]string{"id": "2"}},
		{ID: 3, Values: map[string]string{"id": "3"}},
	}
	tw.SetData(rows)

	if tw.Cursor() != 0 {
		t.Errorf("expected cursor 0, got %d", tw.Cursor())
	}

	tw.HandleKey("down")
	if tw.Cursor() != 1 {
		t.Errorf("expected cursor 1 after down, got %d", tw.Cursor())
	}

	tw.HandleKey("up")
	if tw.Cursor() != 0 {
		t.Errorf("expected cursor 0 after up, got %d", tw.Cursor())
	}

	// Up at top stays at 0
	tw.HandleKey("up")
	if tw.Cursor() != 0 {
		t.Errorf("expected cursor 0 at boundary, got %d", tw.Cursor())
	}
}

func TestTableWidgetOpenDetail(t *testing.T) {
	tw := NewTableWidget("Test", []TableColumn{{Header: "ID", Width: 5, Field: "id"}}, 10, "")
	tw.SetData([]TableRow{{ID: 1, Values: map[string]string{"id": "1"}}})

	_, _, _, openDetail, _, _ := tw.HandleKey("enter")
	if !openDetail {
		t.Error("expected openDetail on enter")
	}
}

func TestTableWidgetView(t *testing.T) {
	tw := NewTableWidget("Test", []TableColumn{{Header: "ID", Width: 5, Field: "id"}}, 10, "help text")
	view := tw.View()
	if view == "" {
		t.Error("table view should not be empty")
	}
}
