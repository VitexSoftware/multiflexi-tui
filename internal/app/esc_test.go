package app

import (
	"testing"

	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type stubView struct{}

func (stubView) Init() tea.Cmd                       { return nil }
func (stubView) Update(tea.Msg) (tea.Model, tea.Cmd) { return stubView{}, nil }
func (stubView) View() string                        { return "stub" }

func TestHandleEscClosesTopLevelView(t *testing.T) {
	a := &App{
		items:     []MenuItem{{Label: "Status"}, {Label: "Apps"}},
		menuFocus: false,
	}
	a.activeView = stubView{}
	a.activeMenuItem = 1

	_, _ = a.handleEsc()
	if a.activeView != nil {
		t.Fatal("expected activeView cleared")
	}
	if !a.menuFocus {
		t.Fatal("expected menu focused")
	}
}

func TestHandleEscPopsNestedView(t *testing.T) {
	a := &App{menuFocus: false}
	list := stubView{}
	a.nav.Push(ViewState{View: list, MenuIdx: 1})
	a.activeView = stubView{}

	_, _ = a.handleEsc()
	if a.activeView != list {
		t.Fatal("expected list restored")
	}
	if a.nav.Depth() != 0 {
		t.Fatalf("expected empty stack, got %d", a.nav.Depth())
	}
	if a.menuFocus {
		t.Fatal("expected content focus on restored list")
	}
}

func TestHandleEscViewerRefreshOnBack(t *testing.T) {
	list := &refreshableStub{}
	a := &App{menuFocus: false}
	a.nav.Push(ViewState{View: list, MenuIdx: 1})
	v := ui.NewViewer("result")
	v.RefreshOnBack = true
	a.activeView = v

	_, _ = a.handleEsc()
	if a.activeView != list {
		t.Fatal("expected list restored")
	}
	if !list.refreshed {
		t.Fatal("expected list.Refresh() called")
	}
}

type refreshableStub struct {
	stubView
	refreshed bool
}

func (r *refreshableStub) Refresh() tea.Cmd {
	r.refreshed = true
	return nil
}
