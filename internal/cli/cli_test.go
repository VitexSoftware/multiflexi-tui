package cli

import (
	"encoding/json"
	"testing"
)

func TestArtifactStructure(t *testing.T) {
	// Test JSON unmarshaling for enhanced Artifact structure
	jsonData := `{
		"id": 52779,
		"job_id": 143770,
		"filename": "stdout.txt",
		"content_type": "text/plain",
		"artifact": "some content",
		"created_at": "2026-02-06T10:00:00Z",
		"note": "Test artifact"
	}`

	var artifact Artifact
	err := json.Unmarshal([]byte(jsonData), &artifact)
	if err != nil {
		t.Fatalf("Failed to unmarshal Artifact: %v", err)
	}

	// Verify all fields are populated correctly
	if artifact.ID != 52779 {
		t.Errorf("Expected ID 52779, got %d", artifact.ID)
	}
	if artifact.Job_ID != 143770 {
		t.Errorf("Expected Job_ID 143770, got %d", artifact.Job_ID)
	}
	if artifact.Filename != "stdout.txt" {
		t.Errorf("Expected Filename 'stdout.txt', got '%s'", artifact.Filename)
	}
	if artifact.ContentType != "text/plain" {
		t.Errorf("Expected ContentType 'text/plain', got '%s'", artifact.ContentType)
	}
	if artifact.Artifact != "some content" {
		t.Errorf("Expected Artifact 'some content', got '%s'", artifact.Artifact)
	}
	if artifact.CreatedAt != "2026-02-06T10:00:00Z" {
		t.Errorf("Expected CreatedAt '2026-02-06T10:00:00Z', got '%s'", artifact.CreatedAt)
	}
	if artifact.Note != "Test artifact" {
		t.Errorf("Expected Note 'Test artifact', got '%s'", artifact.Note)
	}
}

func TestCredTypeStructure(t *testing.T) {
	// Test JSON unmarshaling for enhanced CredType structure
	jsonData := `{
		"id": 10,
		"name": "Office365 App",
		"uuid": "30e25903-7db9-4629-8151-2952305b6987",
		"class": "Office365",
		"company_id": 5,
		"logo": "https://example.com/logo.png",
		"url": "https://office365.com",
		"version": 1
	}`

	var credType CredType
	err := json.Unmarshal([]byte(jsonData), &credType)
	if err != nil {
		t.Fatalf("Failed to unmarshal CredType: %v", err)
	}

	// Verify all fields are populated correctly
	if credType.ID != 10 {
		t.Errorf("Expected ID 10, got %d", credType.ID)
	}
	if credType.Name != "Office365 App" {
		t.Errorf("Expected Name 'Office365 App', got '%s'", credType.Name)
	}
	if credType.Class != "Office365" {
		t.Errorf("Expected Class 'Office365', got '%s'", credType.Class)
	}
	if credType.Version != 1 {
		t.Errorf("Expected Version 1, got %d", credType.Version)
	}
}

func TestCrPrototypeStructure(t *testing.T) {
	// Test JSON unmarshaling for new CrPrototype structure
	jsonData := `{
		"id": 9,
		"name": "Test Credential Prototype",
		"version": "1.0.0",
		"description": "A test credential prototype"
	}`

	var crPrototype CrPrototype
	err := json.Unmarshal([]byte(jsonData), &crPrototype)
	if err != nil {
		t.Fatalf("Failed to unmarshal CrPrototype: %v", err)
	}

	// Verify all fields are populated correctly
	if crPrototype.ID != 9 {
		t.Errorf("Expected ID 9, got %d", crPrototype.ID)
	}
	if crPrototype.Name != "Test Credential Prototype" {
		t.Errorf("Expected Name 'Test Credential Prototype', got '%s'", crPrototype.Name)
	}
	if crPrototype.Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", crPrototype.Version)
	}
}

func TestJobStructure(t *testing.T) {
	jsonData := `{
		"id": 100,
		"app_id": 55,
		"command": "test-command",
		"executor": "Native",
		"schedule_type": "hourly",
		"exitcode": 0,
		"pid": 12345,
		"env": {"FOO": "bar"},
		"company_id": 3
	}`

	var job Job
	err := json.Unmarshal([]byte(jsonData), &job)
	if err != nil {
		t.Fatalf("Failed to unmarshal Job: %v", err)
	}

	if job.ID != 100 {
		t.Errorf("Expected ID 100, got %d", job.ID)
	}
	if job.Command != "test-command" {
		t.Errorf("Expected Command 'test-command', got '%s'", job.Command)
	}
	if job.Env["FOO"] != "bar" {
		t.Errorf("Expected Env[FOO]='bar', got '%s'", job.Env["FOO"])
	}
	if job.ScheduleType != "hourly" {
		t.Errorf("Expected ScheduleType 'hourly', got '%s'", job.ScheduleType)
	}
}

func TestCompanyStructure(t *testing.T) {
	jsonData := `{
		"id": 5,
		"name": "Test Company",
		"email": "test@example.com",
		"ic": "12345678",
		"slug": "test-company",
		"enabled": 1,
		"customer": null
	}`

	var company Company
	err := json.Unmarshal([]byte(jsonData), &company)
	if err != nil {
		t.Fatalf("Failed to unmarshal Company: %v", err)
	}

	if company.ID != 5 {
		t.Errorf("Expected ID 5, got %d", company.ID)
	}
	if company.Name != "Test Company" {
		t.Errorf("Expected Name 'Test Company', got '%s'", company.Name)
	}
	if company.Customer != nil {
		t.Errorf("Expected Customer nil, got %v", company.Customer)
	}
}

func TestQueueStructure(t *testing.T) {
	jsonData := `{
		"id": 1,
		"job": 42,
		"schedule_type": "hourly",
		"runtemplate_id": 10,
		"runtemplate_name": "My Template",
		"app_id": 5,
		"app_name": "My App",
		"company_id": 3,
		"company_name": "My Company",
		"after": "2026-01-01 12:00:00"
	}`

	var queue Queue
	err := json.Unmarshal([]byte(jsonData), &queue)
	if err != nil {
		t.Fatalf("Failed to unmarshal Queue: %v", err)
	}

	if queue.Job != 42 {
		t.Errorf("Expected Job 42, got %d", queue.Job)
	}
	if queue.RunTemplateName != "My Template" {
		t.Errorf("Expected RunTemplateName 'My Template', got '%s'", queue.RunTemplateName)
	}
	if queue.AppName != "My App" {
		t.Errorf("Expected AppName 'My App', got '%s'", queue.AppName)
	}
}

func TestStatusInfoStructure(t *testing.T) {
	// Test JSON unmarshaling for StatusInfo structure
	jsonData := `{
		"version-cli": "2.3.2.110",
		"user": "root",
		"memory": 4259072,
		"companies": 5,
		"apps": 51,
		"encryption": "active (3 keys)"
	}`

	var status StatusInfo
	err := json.Unmarshal([]byte(jsonData), &status)
	if err != nil {
		t.Fatalf("Failed to unmarshal StatusInfo: %v", err)
	}

	// Verify key fields
	if status.VersionCli != "2.3.2.110" {
		t.Errorf("Expected VersionCli '2.3.2.110', got '%s'", status.VersionCli)
	}
	if status.User != "root" {
		t.Errorf("Expected User 'root', got '%s'", status.User)
	}
	if status.Memory != 4259072 {
		t.Errorf("Expected Memory 4259072, got %d", status.Memory)
	}
}
