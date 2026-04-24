package entity

import (
	"testing"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

func TestCompanyDetailAndEditor(t *testing.T) {
	co := cli.Company{ID: 1, Name: "Acme", Email: "a@b.c", IC: "123", Slug: "acme", Enabled: 1, Server: 2}
	fields := CompanyDef.ToDetail(co)
	if len(fields) == 0 {
		t.Fatal("ToDetail returned empty")
	}
	if fields[0].Label != "ID" || fields[0].Value != "1" {
		t.Errorf("first field: %+v", fields[0])
	}

	ef := CompanyDef.ToEditor(co)
	if len(ef) != 4 {
		t.Fatalf("expected 4 editor fields, got %d", len(ef))
	}
	if ef[0].Value != "Acme" {
		t.Errorf("editor Name value = %q", ef[0].Value)
	}

	args := CompanyDef.UpdateArgs(co, map[string]string{"Name": "New", "Email": "x@y", "IC": "999", "Slug": "new"})
	if len(args) < 8 {
		t.Errorf("expected >= 8 update args, got %d: %v", len(args), args)
	}

	nf := CompanyDef.NewFields()
	if len(nf) != 4 {
		t.Errorf("expected 4 new fields, got %d", len(nf))
	}

	cargs := CompanyDef.CreateArgs(map[string]string{"Name": "Test", "Email": "e", "IC": "", "Slug": ""})
	found := false
	for _, a := range cargs {
		if a == "Test" {
			found = true
		}
	}
	if !found {
		t.Errorf("create args missing name: %v", cargs)
	}

	if CompanyDef.GetID(co) != 1 {
		t.Errorf("GetID = %d", CompanyDef.GetID(co))
	}
	if CompanyDef.GetLabel(co) != "Company: Acme" {
		t.Errorf("GetLabel = %q", CompanyDef.GetLabel(co))
	}
}

func TestJobDetailAndEditor(t *testing.T) {
	j := cli.Job{ID: 50, Command: "run-sync", Executor: "Native", ScheduleType: "daily", PID: 0, Exitcode: 0}
	fields := JobDef.ToDetail(j)
	if len(fields) < 10 {
		t.Errorf("expected >=10 detail fields, got %d", len(fields))
	}
	ef := JobDef.ToEditor(j)
	if len(ef) != 2 {
		t.Fatalf("expected 2 editor fields, got %d", len(ef))
	}
	if ef[0].Value != "Native" {
		t.Errorf("editor Executor = %q", ef[0].Value)
	}

	nf := JobDef.NewFields()
	if len(nf) != 4 {
		t.Errorf("expected 4 new fields for job, got %d", len(nf))
	}
}

func TestApplicationDetailAndEditor(t *testing.T) {
	a := cli.Application{ID: 1, Name: "MyApp", Version: "1.0", UUID: "abc", Executable: "/bin/app", Enabled: 1}
	fields := ApplicationDef.ToDetail(a)
	if len(fields) == 0 {
		t.Fatal("ToDetail returned empty")
	}
	ef := ApplicationDef.ToEditor(a)
	if len(ef) != 5 {
		t.Fatalf("expected 5 editor fields, got %d", len(ef))
	}
	nf := ApplicationDef.NewFields()
	if len(nf) != 5 {
		t.Errorf("expected 5 new fields, got %d", len(nf))
	}
}

func TestRunTemplateDetailAndEditor(t *testing.T) {
	rt := cli.RunTemplate{ID: 10, Name: "Daily", AppID: 5, CompanyID: 3, Interv: "d", Active: 1, Executor: "Native"}
	fields := RunTemplateDef.ToDetail(rt)
	if len(fields) < 8 {
		t.Errorf("expected >=8 detail fields, got %d", len(fields))
	}
	ef := RunTemplateDef.ToEditor(rt)
	if len(ef) < 3 {
		t.Fatalf("expected >=3 editor fields, got %d", len(ef))
	}
}

func TestAllEntitiesHaveGetIDAndLabel(t *testing.T) {
	for _, e := range All {
		if e.Def.GetID == nil {
			t.Errorf("%s: GetID is nil", e.Label)
		}
		if e.Def.GetLabel == nil {
			t.Errorf("%s: GetLabel is nil", e.Label)
		}
		if e.Def.ToDetail == nil {
			t.Errorf("%s: ToDetail is nil", e.Label)
		}
	}
}

func TestAllEntitiesDeleteAction(t *testing.T) {
	expected := map[string]string{
		"company":     "remove",
		"credential":  "remove",
		"eventsource": "remove",
		"eventrule":   "remove",
	}
	for _, e := range All {
		da := e.Def.DeleteAction
		if da != "delete" && da != "remove" {
			t.Errorf("%s: DeleteAction = %q, want 'delete' or 'remove'", e.Label, da)
		}
		if exp, ok := expected[e.Def.CLIEntity]; ok && da != exp {
			t.Errorf("%s (CLI: %s): DeleteAction = %q, want %q", e.Label, e.Def.CLIEntity, da, exp)
		}
	}
}

func TestCompanyAppListActions(t *testing.T) {
	if len(CompanyAppDef.ListActions) != 2 {
		t.Fatalf("expected 2 ListActions (assign, unassign), got %d", len(CompanyAppDef.ListActions))
	}

	assign := CompanyAppDef.ListActions[0]
	if assign.Key != "a" {
		t.Errorf("assign key = %q, want 'a'", assign.Key)
	}
	if assign.Label != "Assign" {
		t.Errorf("assign label = %q, want 'Assign'", assign.Label)
	}

	unassign := CompanyAppDef.ListActions[1]
	if unassign.Key != "u" {
		t.Errorf("unassign key = %q, want 'u'", unassign.Key)
	}
	if unassign.Label != "Unassign" {
		t.Errorf("unassign label = %q, want 'Unassign'", unassign.Label)
	}

	if len(CompanyAppDef.Actions) != 0 {
		t.Errorf("expected 0 row Actions, got %d", len(CompanyAppDef.Actions))
	}
}

func TestCompanyAppToDetail(t *testing.T) {
	ca := cli.CompanyApp{ID: 7, CompanyID: 3, AppID: 5}
	fields := CompanyAppDef.ToDetail(ca)
	if len(fields) != 3 {
		t.Fatalf("expected 3 detail fields, got %d", len(fields))
	}
	if fields[0].Label != "ID" || fields[0].Value != "7" {
		t.Errorf("first field: %+v", fields[0])
	}
	if CompanyAppDef.GetID(ca) != 7 {
		t.Errorf("GetID = %d", CompanyAppDef.GetID(ca))
	}
}

func TestConfirmDialog(t *testing.T) {
	called := false
	action := func() ui.StatusMsg {
		called = true
		return ui.StatusMsg{Text: "done"}
	}
	_ = action
	d := ui.NewConfirmDialog("Delete?", func() tea.Msg { called = true; return nil })
	view := d.View()
	if view == "" {
		t.Error("confirm dialog view is empty")
	}
	_ = called
}

func TestViewerView(t *testing.T) {
	v := ui.NewViewer("Help")
	v.SetContent("Help", "line1\nline2\nline3")
	view := v.View()
	if view == "" {
		t.Error("viewer view is empty")
	}
}
