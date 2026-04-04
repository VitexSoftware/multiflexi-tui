package app

import "testing"

func TestNavigatorPushPop(t *testing.T) {
	nav := Navigator{}
	if nav.Depth() != 0 {
		t.Errorf("expected depth 0, got %d", nav.Depth())
	}
	nav.Push(ViewState{MenuIdx: 1})
	nav.Push(ViewState{MenuIdx: 2})
	if nav.Depth() != 2 {
		t.Errorf("expected depth 2, got %d", nav.Depth())
	}
	top, ok := nav.Pop()
	if !ok || top.MenuIdx != 2 {
		t.Errorf("expected MenuIdx 2, got %d", top.MenuIdx)
	}
	top, ok = nav.Pop()
	if !ok || top.MenuIdx != 1 {
		t.Errorf("expected MenuIdx 1, got %d", top.MenuIdx)
	}
	_, ok = nav.Pop()
	if ok {
		t.Error("expected false on empty pop")
	}
}

func TestNavigatorClear(t *testing.T) {
	nav := Navigator{}
	nav.Push(ViewState{MenuIdx: 1})
	nav.Push(ViewState{MenuIdx: 2})
	nav.Clear()
	if nav.Depth() != 0 {
		t.Errorf("expected depth 0 after clear, got %d", nav.Depth())
	}
}
