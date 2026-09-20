package ui

import (
	"strings"
	"testing"
)

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

func TestTableWidgetSetContentHeight(t *testing.T) {
	tw := NewTableWidget("Test", []TableColumn{{Header: "ID", Width: 5, Field: "id"}}, 10, "")
	if tw.Limit() != 10 {
		t.Fatalf("initial limit = %d, want 10", tw.Limit())
	}
	if !tw.SetContentHeight(25) {
		t.Fatal("expected limit change")
	}
	if tw.Limit() != 20 {
		t.Errorf("limit = %d, want 20", tw.Limit())
	}
	if tw.SetContentHeight(25) {
		t.Error("same height should not report a change")
	}
	if !tw.SetContentHeight(4) {
		t.Fatal("expected clamp change")
	}
	if tw.Limit() != minTableRows {
		t.Errorf("limit = %d, want %d", tw.Limit(), minTableRows)
	}
}

func TestTableWidgetFlexWidth(t *testing.T) {
	tw := NewTableWidget("Apps", []TableColumn{
		{Header: "ID", Width: 5, Field: "id"},
		{Header: "Name", Width: 30, Field: "name", Flex: true},
		{Header: "Status", Width: 10, Field: "status"},
	}, 10, "")
	tw.SetContentWidth(80)
	widths := tw.effectiveWidths()
	// used = 1 + 5 + 1 + 30 + 1 + 10 = 48; extra = 80-48 = 32 → Name = 62
	if widths[1] != 62 {
		t.Errorf("flex Name width = %d, want 62 (got %v)", widths[1], widths)
	}
	if widths[0] != 5 || widths[2] != 10 {
		t.Errorf("fixed columns changed: %v", widths)
	}

	long := "AbraFlexi Revolut statement downloader for SPOJE"
	tw.SetData([]TableRow{{ID: 1, Values: map[string]string{
		"id": "96", "name": long, "status": "Disabled",
	}}})
	view := tw.View()
	if strings.Contains(view, "...") {
		t.Errorf("name should not be truncated; view=%q", view)
	}
	if !strings.Contains(view, long) {
		t.Errorf("expected full name in view, got %q", view)
	}
}

func TestTableWidgetAutoFlex(t *testing.T) {
	// No explicit Flex — columns with Width >= 20 share leftover space
	tw := NewTableWidget("T", []TableColumn{
		{Header: "ID", Width: 5, Field: "id"},
		{Header: "Name", Width: 30, Field: "name"},
		{Header: "Note", Width: 20, Field: "note"},
	}, 5, "")
	tw.SetContentWidth(100)
	widths := tw.effectiveWidths()
	// used = 1+5+1+30+1+20 = 58; extra = 42 shared by Name+Note → +21 each
	if widths[1] != 51 || widths[2] != 41 {
		t.Errorf("auto-flex widths = %v, want Name=51 Note=41", widths)
	}
}
