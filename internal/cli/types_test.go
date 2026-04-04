package cli

import (
	"encoding/json"
	"testing"
)

func TestStatusInfoUnmarshal(t *testing.T) {
	data := `{"version-cli":"2.3.2","user":"root","memory":4259072,"companies":5,"apps":51,"encryption":"active (3 keys)"}`
	var s StatusInfo
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.VersionCli != "2.3.2" {
		t.Errorf("VersionCli = %q, want 2.3.2", s.VersionCli)
	}
	if s.Memory != 4259072 {
		t.Errorf("Memory = %d, want 4259072", s.Memory)
	}
	if s.Companies != 5 {
		t.Errorf("Companies = %d, want 5", s.Companies)
	}
}

func TestCompanyUnmarshal(t *testing.T) {
	data := `{"id":5,"name":"Test Co","email":"a@b.c","ic":"12345678","slug":"test-co","enabled":1,"customer":null}`
	var c Company
	if err := json.Unmarshal([]byte(data), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.ID != 5 || c.Name != "Test Co" || c.Customer != nil || c.Enabled != 1 {
		t.Errorf("unexpected: %+v", c)
	}
}

func TestJobUnmarshal(t *testing.T) {
	data := `{"id":100,"command":"test-cmd","executor":"Native","schedule_type":"hourly","exitcode":0,"pid":12345,"env":{"FOO":"bar"},"company_id":3,"app_id":55}`
	var j Job
	if err := json.Unmarshal([]byte(data), &j); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if j.ID != 100 || j.Command != "test-cmd" || j.Env["FOO"] != "bar" || j.PID != 12345 {
		t.Errorf("unexpected: %+v", j)
	}
}

func TestApplicationUnmarshal(t *testing.T) {
	data := `{"id":1,"name":"MyApp","version":"1.0","uuid":"abc-123","executable":"/usr/bin/app","enabled":1}`
	var a Application
	if err := json.Unmarshal([]byte(data), &a); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if a.Name != "MyApp" || a.UUID != "abc-123" || a.Enabled != 1 {
		t.Errorf("unexpected: %+v", a)
	}
}

func TestRunTemplateUnmarshal(t *testing.T) {
	data := `{"id":10,"name":"Daily Sync","app_id":5,"company_id":3,"interv":"d","active":1,"executor":"Native","successfull_jobs_count":42,"failed_jobs_count":2}`
	var rt RunTemplate
	if err := json.Unmarshal([]byte(data), &rt); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if rt.Name != "Daily Sync" || rt.SuccessfulJobsCount != 42 || rt.Active != 1 {
		t.Errorf("unexpected: %+v", rt)
	}
}

func TestEventSourceUnmarshal(t *testing.T) {
	data := `{"id":1,"name":"Webhook DB","adapter_type":"abraflexi-webhook","db_host":"localhost","db_database":"events","enabled":1,"poll_interval":60}`
	var es EventSource
	if err := json.Unmarshal([]byte(data), &es); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if es.Name != "Webhook DB" || es.PollInterval != 60 || es.Enabled != 1 {
		t.Errorf("unexpected: %+v", es)
	}
}

func TestEventRuleUnmarshal(t *testing.T) {
	data := `{"id":5,"event_source_id":1,"evidence":"faktura-vydana","operation":"create","runtemplate_id":10,"priority":5,"enabled":1,"env_mapping":"{}"}`
	var er EventRule
	if err := json.Unmarshal([]byte(data), &er); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if er.Evidence != "faktura-vydana" || er.Operation != "create" || er.RunTemplateID != 10 {
		t.Errorf("unexpected: %+v", er)
	}
}

func TestQueueUnmarshal(t *testing.T) {
	data := `{"id":1,"job":42,"schedule_type":"hourly","runtemplate_id":10,"runtemplate_name":"My Template","app_id":5,"app_name":"My App","company_id":3,"company_name":"My Co","after":"2026-01-01 12:00:00"}`
	var q Queue
	if err := json.Unmarshal([]byte(data), &q); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if q.Job != 42 || q.AppName != "My App" || q.CompanyName != "My Co" {
		t.Errorf("unexpected: %+v", q)
	}
}

func TestArtifactUnmarshal(t *testing.T) {
	data := `{"id":100,"job_id":200,"filename":"stdout.txt","content_type":"text/plain","artifact":"data","created_at":"2026-01-01","note":"test"}`
	var a Artifact
	if err := json.Unmarshal([]byte(data), &a); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if a.Filename != "stdout.txt" || a.JobID != 200 || a.Note != "test" {
		t.Errorf("unexpected: %+v", a)
	}
}
