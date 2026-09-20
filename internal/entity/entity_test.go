package entity

import (
	"fmt"
	"strings"
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
	exit := 0
	j := cli.Job{ID: 50, Command: "run-sync", Executor: "Native", ScheduleType: "daily", PID: 0, Exitcode: &exit}
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

	if JobDef.Columns[1].Field != "exitcode" {
		t.Errorf("expected Exit column second, got %+v", JobDef.Columns[1])
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

	var hasConfig bool
	for _, act := range ApplicationDef.Actions {
		if act.Key == "c" && act.Handler != nil {
			hasConfig = true
		}
	}
	if !hasConfig {
		t.Error("ApplicationDef missing Config action (c) with Handler")
	}
}

func TestConfFieldDef(t *testing.T) {
	def := NewConfFieldDef(2, "DemoApp")
	if def.CLIEntity != "conffield" {
		t.Errorf("CLIEntity = %q", def.CLIEntity)
	}
	f := cli.ConfField{
		ID: 10, AppID: 2, Keyname: "FOO", Type: "string", Description: "d",
		Hint: "h", Note: "n", Defval: "v", Required: 1, Secret: 0, Multiline: 0, Expiring: 1,
	}
	detail := def.ToDetail(f)
	if len(detail) < 10 {
		t.Fatalf("expected >=10 detail fields, got %d", len(detail))
	}
	ef := def.ToEditor(f)
	if len(ef) != 10 {
		t.Fatalf("expected 10 editor fields, got %d", len(ef))
	}
	args := def.UpdateArgs(f, map[string]string{
		"Keyword": "FOO2", "Type": "bool", "Description": "x", "Hint": "", "Note": "",
		"Default": "1", "Required": "1", "Secret": "0", "Multiline": "0", "Expiring": "0",
	})
	joined := fmt.Sprintf("%v", args)
	if !containsAll(joined, "--id", "10", "--keyname", "FOO2", "--type", "bool") {
		t.Errorf("update args incomplete: %v", args)
	}
	nf := def.NewFields()
	if len(nf) != 10 {
		t.Errorf("expected 10 new fields, got %d", len(nf))
	}
	cargs := def.CreateArgs(map[string]string{
		"Keyword": "BAR", "Type": "string", "Description": "desc", "Required": "1",
	})
	joined = fmt.Sprintf("%v", cargs)
	if !containsAll(joined, "--app_id", "2", "--keyname", "BAR", "--type", "string", "--required", "1") {
		t.Errorf("create args incomplete: %v", cargs)
	}
	if def.GetID(f) != 10 {
		t.Errorf("GetID = %d", def.GetID(f))
	}
	if def.GetLabel(f) != "Field: FOO" {
		t.Errorf("GetLabel = %q", def.GetLabel(f))
	}
}

func containsAll(hay string, parts ...string) bool {
	for _, p := range parts {
		if !strings.Contains(hay, p) {
			return false
		}
	}
	return true
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
		"company":          "remove",
		"credential":       "remove",
		"event-source":     "remove",
		"event-rule":       "remove",
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
	ca := cli.CompanyApp{ID: 7, CompanyID: 3, CompanyName: "Acme", AppID: 5, AppName: "Checker"}
	fields := CompanyAppDef.ToDetail(ca)
	if len(fields) != 7 {
		t.Fatalf("expected 7 detail fields, got %d", len(fields))
	}
	if fields[0].Label != "ID" || fields[0].Value != "7" {
		t.Errorf("first field: %+v", fields[0])
	}
	if CompanyAppDef.GetID(ca) != 7 {
		t.Errorf("GetID = %d", CompanyAppDef.GetID(ca))
	}
	if CompanyAppDef.GetLabel(ca) != "Acme → Checker" {
		t.Errorf("GetLabel = %q", CompanyAppDef.GetLabel(ca))
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
