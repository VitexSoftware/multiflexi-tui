package app

import "testing"

func TestContentSizeMsg(t *testing.T) {
	a := &App{width: 100, height: 40}
	msg := a.contentSizeMsg()
	// menuBarLines(3) + footerLines(2) = 5 → content height 35
	if msg.Height != 35 {
		t.Errorf("content height = %d, want 35", msg.Height)
	}
	if msg.Width != 100 {
		t.Errorf("content width = %d, want 100", msg.Width)
	}

	a.statusMessage = "done"
	msg = a.contentSizeMsg()
	if msg.Height != 34 {
		t.Errorf("with status, content height = %d, want 34", msg.Height)
	}
}

func TestApplyContentSizeNilSafe(t *testing.T) {
	a := &App{width: 80, height: 24}
	if a.applyContentSize(nil) != nil {
		t.Error("expected nil")
	}
}
