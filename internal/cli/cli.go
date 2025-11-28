package cli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// Command represents a single command from multiflexi-cli
type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CommandInfo represents the command information structure from multiflexi-cli describe
type CommandInfo struct {
	Description string      `json:"description"`
	Arguments   interface{} `json:"arguments,omitempty"`
	Options     interface{} `json:"options,omitempty"`
}

// GetCommands runs "multiflexi-cli describe" and parses the JSON output
func GetCommands() ([]Command, error) {
	cmd := exec.Command("multiflexi-cli", "describe")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli describe: %w", err)
	}

	// Parse the JSON into a map where keys are command names and values are CommandInfo
	var cmdMap map[string]CommandInfo
	if err := json.Unmarshal(output, &cmdMap); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	// Convert the map to a slice of Command structs
	var commands []Command
	for name, info := range cmdMap {
		// Skip internal commands that start with underscore
		if strings.HasPrefix(name, "_") {
			continue
		}
		commands = append(commands, Command{
			Name:        name,
			Description: info.Description,
		})
	}

	// Sort commands alphabetically by name for consistent display
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands, nil
}

// GetStatus runs "multiflexi-cli status --format=json" and returns formatted status
func GetStatus() (string, error) {
	cmd := exec.Command("multiflexi-cli", "status", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		return "Status unavailable", fmt.Errorf("failed to run multiflexi-cli status: %w", err)
	}

	outputStr := strings.TrimSpace(string(output))

	// Try to parse as JSON and format key information
	var statusData map[string]interface{}
	if err := json.Unmarshal([]byte(outputStr), &statusData); err != nil {
		// If JSON parsing fails, return raw output (truncated)
		if len(outputStr) > 200 {
			return outputStr[:200] + "...", nil
		}
		return outputStr, nil
	}

	// Format comprehensive status information
	var statusParts []string

	if version, ok := statusData["version-cli"].(string); ok {
		statusParts = append(statusParts, fmt.Sprintf("CLI: %s", version))
	}

	if user, ok := statusData["user"].(string); ok {
		statusParts = append(statusParts, fmt.Sprintf("User: %s", user))
	}

	if php, ok := statusData["php"].(string); ok {
		statusParts = append(statusParts, fmt.Sprintf("PHP: %s", php))
	}

	if companies, ok := statusData["companies"].(float64); ok {
		statusParts = append(statusParts, fmt.Sprintf("Companies: %.0f", companies))
	}

	if apps, ok := statusData["apps"].(float64); ok {
		statusParts = append(statusParts, fmt.Sprintf("Apps: %.0f", apps))
	}

	if executor, ok := statusData["executor"].(string); ok {
		statusParts = append(statusParts, fmt.Sprintf("Executor: %s", executor))
	}

	if scheduler, ok := statusData["scheduler"].(string); ok {
		statusParts = append(statusParts, fmt.Sprintf("Scheduler: %s", scheduler))
	}

	if len(statusParts) == 0 {
		return "Status loaded", nil
	}

	return strings.Join(statusParts, " | "), nil
}

// StatusInfo represents comprehensive system status
type StatusInfo struct {
	Version    string
	User       string
	PHP        string
	OS         string
	Companies  int
	Apps       int
	Templates  int
	Executor   string
	Scheduler  string
	Zabbix     string
	Telemetry  string
	Encryption string
	Database   string
}

// GetStatusInfo returns comprehensive status information
func GetStatusInfo() (*StatusInfo, error) {
	cmd := exec.Command("multiflexi-cli", "status", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli status: %w", err)
	}

	outputStr := strings.TrimSpace(string(output))

	var statusData map[string]interface{}
	if err := json.Unmarshal([]byte(outputStr), &statusData); err != nil {
		return nil, fmt.Errorf("failed to parse status JSON: %w", err)
	}

	status := &StatusInfo{}

	if version, ok := statusData["version-cli"].(string); ok {
		status.Version = version
	}

	if user, ok := statusData["user"].(string); ok {
		status.User = user
	}

	if php, ok := statusData["php"].(string); ok {
		status.PHP = php
	}

	if os, ok := statusData["os"].(string); ok {
		status.OS = os
	}

	if companies, ok := statusData["companies"].(float64); ok {
		status.Companies = int(companies)
	}

	if apps, ok := statusData["apps"].(float64); ok {
		status.Apps = int(apps)
	}

	if templates, ok := statusData["runtemplates"].(float64); ok {
		status.Templates = int(templates)
	}

	if executor, ok := statusData["executor"].(string); ok {
		status.Executor = executor
	}

	if scheduler, ok := statusData["scheduler"].(string); ok {
		status.Scheduler = scheduler
	}

	if zabbix, ok := statusData["zabbix"].(string); ok {
		status.Zabbix = zabbix
	}

	if telemetry, ok := statusData["telemetry"].(string); ok {
		status.Telemetry = telemetry
	}

	if encryption, ok := statusData["encryption"].(string); ok {
		status.Encryption = encryption
	}

	if database, ok := statusData["database"].(string); ok {
		// Truncate database info for display
		if len(database) > 50 {
			status.Database = database[:50] + "..."
		} else {
			status.Database = database
		}
	}

	return status, nil
}

// Application represents an application from multiflexi-cli
type Application struct {
	ID          int    `json:"id"`
	Enabled     int    `json:"enabled"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Executable  string `json:"executable"`
	Version     string `json:"version"`
	Code        string `json:"code"`
	Topics      string `json:"topics"`
}

// Company represents a company from multiflexi-cli
type Company struct {
	ID      int    `json:"id"`
	Enabled int    `json:"enabled"`
	Name    string `json:"name"`
	IC      string `json:"ic"`
	Email   string `json:"email"`
	Slug    string `json:"slug"`
}

// RunTemplate represents a run template from multiflexi-cli
type RunTemplate struct {
	ID           int    `json:"id"`
	AppID        int    `json:"app_id"`
	CompanyID    int    `json:"company_id"`
	Name         string `json:"name"`
	Interv       string `json:"interv"`
	Active       int    `json:"active"`
	Executor     string `json:"executor"`
	Cron         string `json:"cron"`
	LastSchedule string `json:"last_schedule"`
	NextSchedule string `json:"next_schedule"`
}

// GetApplications fetches applications from multiflexi-cli with pagination
func GetApplications(limit, offset int) ([]Application, error) {
	cmd := exec.Command("multiflexi-cli", "application", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli application list: %w", err)
	}

	var apps []Application
	if err := json.Unmarshal(output, &apps); err != nil {
		return nil, fmt.Errorf("failed to parse application JSON output: %w", err)
	}

	return apps, nil
}

// GetCompanies fetches companies from multiflexi-cli with pagination
func GetCompanies(limit, offset int) ([]Company, error) {
	cmd := exec.Command("multiflexi-cli", "company", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli company list: %w", err)
	}

	var companies []Company
	if err := json.Unmarshal(output, &companies); err != nil {
		return nil, fmt.Errorf("failed to parse company JSON output: %w", err)
	}

	return companies, nil
}

// GetRunTemplates fetches run templates from multiflexi-cli with pagination
func GetRunTemplates(limit, offset int) ([]RunTemplate, error) {
	cmd := exec.Command("multiflexi-cli", "runtemplate", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli runtemplate list: %w", err)
	}

	var templates []RunTemplate
	if err := json.Unmarshal(output, &templates); err != nil {
		return nil, fmt.Errorf("failed to parse runtemplate JSON output: %w", err)
	}

	return templates, nil
}

// Credential represents a credential from multiflexi-cli
type Credential struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	CompanyID        int    `json:"company_id"`
	CredentialTypeID int    `json:"credential_type_id"`
}

// GetCredentials fetches credentials from multiflexi-cli with pagination
func GetCredentials(limit, offset int) ([]Credential, error) {
	cmd := exec.Command("multiflexi-cli", "credential", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli credential list: %w", err)
	}

	var credentials []Credential
	if err := json.Unmarshal(output, &credentials); err != nil {
		return nil, fmt.Errorf("failed to parse credential JSON output: %w", err)
	}

	return credentials, nil
}

// Token represents a token from multiflexi-cli
type Token struct {
	ID    int    `json:"id"`
	User  string `json:"user"`
	Token string `json:"token"`
}

// GetTokens fetches tokens from multiflexi-cli with pagination
func GetTokens(limit, offset int) ([]Token, error) {
	cmd := exec.Command("multiflexi-cli", "token", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli token list: %w", err)
	}

	var tokens []Token
	if err := json.Unmarshal(output, &tokens); err != nil {
		return nil, fmt.Errorf("failed to parse token JSON output: %w", err)
	}

	return tokens, nil
}

// User represents a user from multiflexi-cli
type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

// GetUsers fetches users from multiflexi-cli with pagination
func GetUsers(limit, offset int) ([]User, error) {
	cmd := exec.Command("multiflexi-cli", "user", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli user list: %w", err)
	}

	var users []User
	if err := json.Unmarshal(output, &users); err != nil {
		return nil, fmt.Errorf("failed to parse user JSON output: %w", err)
	}

	return users, nil
}

// Artifact represents an artifact from multiflexi-cli
type Artifact struct {
	ID     int    `json:"id"`
	Job_ID int    `json:"job_id"`
	File   string `json:"file"`
}

// GetArtifacts fetches artifacts from multiflexi-cli with pagination
func GetArtifacts(limit, offset int) ([]Artifact, error) {
	cmd := exec.Command("multiflexi-cli", "artifact", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli artifact list: %w", err)
	}

	var artifacts []Artifact
	if err := json.Unmarshal(output, &artifacts); err != nil {
		return nil, fmt.Errorf("failed to parse artifact JSON output: %w", err)
	}

	return artifacts, nil
}

// CredType represents a credential type from multiflexi-cli
type CredType struct {
	ID   int    `json:"id"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// GetCredTypes fetches credential types from multiflexi-cli with pagination
func GetCredTypes(limit, offset int) ([]CredType, error) {
	cmd := exec.Command("multiflexi-cli", "credtype", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli credtype list: %w", err)
	}

	var credTypes []CredType
	if err := json.Unmarshal(output, &credTypes); err != nil {
		return nil, fmt.Errorf("failed to parse credtype JSON output: %w", err)
	}

	return credTypes, nil
}

// CompanyApp represents a company-application relation from multiflexi-cli
type CompanyApp struct {
	ID        int `json:"id"`
	CompanyID int `json:"company_id"`
	AppID     int `json:"app_id"`
}

// GetCompanyApps fetches company-application relations from multiflexi-cli with pagination
func GetCompanyApps(limit, offset int) ([]CompanyApp, error) {
	cmd := exec.Command("multiflexi-cli", "companyapp", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli companyapp list: %w", err)
	}

	var companyApps []CompanyApp
	if err := json.Unmarshal(output, &companyApps); err != nil {
		return nil, fmt.Errorf("failed to parse companyapp JSON output: %w", err)
	}

	return companyApps, nil
}

// InitEncryption initializes encryption
func InitEncryption() error {
	cmd := exec.Command("multiflexi-cli", "encryption", "init")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run multiflexi-cli encryption init: %w", err)
	}
	return nil
}

// Prune prunes logs and jobs
func Prune(logs, jobs bool, keep int) error {
	args := []string{"prune"}
	if logs {
		args = append(args, "--logs")
	}
	if jobs {
		args = append(args, "--jobs")
	}
	args = append(args, "--keep", fmt.Sprintf("%d", keep))

	cmd := exec.Command("multiflexi-cli", args...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run multiflexi-cli prune: %w", err)
	}
	return nil
}

// Queue represents a queue item from multiflexi-cli
type Queue struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// GetQueue fetches queue items from multiflexi-cli with pagination
func GetQueue(limit, offset int) ([]Queue, error) {
	cmd := exec.Command("multiflexi-cli", "queue", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli queue list: %w", err)
	}

	var queue []Queue
	if err := json.Unmarshal(output, &queue); err != nil {
		return nil, fmt.Errorf("failed to parse queue JSON output: %w", err)
	}

	return queue, nil
}

// TruncateQueue truncates the queue
func TruncateQueue() error {
	cmd := exec.Command("multiflexi-cli", "queue", "truncate")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run multiflexi-cli queue truncate: %w", err)
	}
	return nil
}

type Job struct {
	ID           int    `json:"id"`
	AppID        int    `json:"app"`
	Command      string `json:"command"`
	Begin        string `json:"begin"`
	End          string `json:"end"`
	Exitcode     int    `json:"exitcode"`
	Executor     string `json:"executor"`
	PID          int    `json:"pid"`
	Schedule     string `json:"schedule"`
	ScheduleType string `json:"schedule_type"`
}

// GetJobs fetches jobs from multiflexi-cli with pagination
func GetJobs(limit, offset int) ([]Job, error) {
	cmd := exec.Command("multiflexi-cli", "job", "list",
		"--format=json",
		"--order=D", // Newer on top (descending order)
		"--limit="+fmt.Sprintf("%d", limit),
		"--offset="+fmt.Sprintf("%d", offset),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run multiflexi-cli job list: %w", err)
	}

	var jobs []Job
	if err := json.Unmarshal(output, &jobs); err != nil {
		return nil, fmt.Errorf("failed to parse job JSON output: %w", err)
	}

	return jobs, nil
}

func getMockJobs(limit, offset int) []Job {
	allJobs := []Job{
		{ID: 100439, AppID: 55, Command: "abraflexi-matcher-out", Begin: "2025-11-10 14:01:02", End: "2025-11-10 14:01:02", Exitcode: 0, Executor: "Native", PID: 85845, Schedule: "2025-11-10 15:01:01", ScheduleType: "hourly"},
		{ID: 100438, AppID: 54, Command: "abraflexi-pull-bank", Begin: "2025-11-10 13:45:32", End: "2025-11-10 13:45:35", Exitcode: 0, Executor: "Native", PID: 84923, Schedule: "2025-11-10 13:45:00", ScheduleType: "hourly"},
		{ID: 100437, AppID: 53, Command: "backup-database", Begin: "2025-11-10 12:30:15", End: "2025-11-10 12:30:47", Exitcode: 0, Executor: "Native", PID: 82156, Schedule: "2025-11-10 12:30:00", ScheduleType: "daily"},
		{ID: 100436, AppID: 52, Command: "sync-contacts", Begin: "2025-11-10 11:15:22", End: "2025-11-10 11:15:28", Exitcode: 0, Executor: "Native", PID: 79334, Schedule: "2025-11-10 11:15:00", ScheduleType: "manual"},
		{ID: 100435, AppID: 51, Command: "generate-reports", Begin: "2025-11-10 10:00:45", End: "2025-11-10 10:02:12", Exitcode: 0, Executor: "Native", PID: 76542, Schedule: "2025-11-10 10:00:00", ScheduleType: "weekly"},
		{ID: 100434, AppID: 50, Command: "cleanup-temp-files", Begin: "2025-11-10 09:30:33", End: "2025-11-10 09:30:41", Exitcode: 1, Executor: "Native", PID: 74823, Schedule: "2025-11-10 09:30:00", ScheduleType: "daily"},
		{ID: 100433, AppID: 49, Command: "update-inventory", Begin: "2025-11-10 08:45:12", End: "2025-11-10 08:47:23", Exitcode: 0, Executor: "Native", PID: 72145, Schedule: "2025-11-10 08:45:00", ScheduleType: "hourly"},
		{ID: 100432, AppID: 48, Command: "process-emails", Begin: "2025-11-10 08:00:05", End: "2025-11-10 08:00:18", Exitcode: 0, Executor: "Native", PID: 69823, Schedule: "2025-11-10 08:00:00", ScheduleType: "hourly"},
		{ID: 100431, AppID: 47, Command: "data-validation", Begin: "2025-11-10 07:30:44", End: "2025-11-10 07:31:02", Exitcode: 2, Executor: "Native", PID: 67234, Schedule: "2025-11-10 07:30:00", ScheduleType: "manual"},
		{ID: 100430, AppID: 46, Command: "archive-old-logs", Begin: "2025-11-10 06:15:18", End: "2025-11-10 06:16:33", Exitcode: 0, Executor: "Native", PID: 65012, Schedule: "2025-11-10 06:15:00", ScheduleType: "weekly"},
		{ID: 100429, AppID: 45, Command: "check-system-health", Begin: "2025-11-10 05:00:27", End: "2025-11-10 05:00:31", Exitcode: 0, Executor: "Native", PID: 62789, Schedule: "2025-11-10 05:00:00", ScheduleType: "daily"},
		{ID: 100428, AppID: 44, Command: "process-queue", Begin: "2025-11-10 04:45:11", End: "2025-11-10 04:45:14", Exitcode: 0, Executor: "Native", PID: 60456, Schedule: "2025-11-10 04:45:00", ScheduleType: "hourly"},
		{ID: 100427, AppID: 43, Command: "sync-external-api", Begin: "2025-11-10 03:30:55", End: "2025-11-10 03:32:18", Exitcode: 1, Executor: "Native", PID: 58123, Schedule: "2025-11-10 03:30:00", ScheduleType: "manual"},
		{ID: 100426, AppID: 42, Command: "optimize-database", Begin: "2025-11-10 02:15:33", End: "2025-11-10 02:18:47", Exitcode: 0, Executor: "Native", PID: 55890, Schedule: "2025-11-10 02:15:00", ScheduleType: "weekly"},
		{ID: 100425, AppID: 41, Command: "monitor-services", Begin: "2025-11-10 01:00:22", End: "2025-11-10 01:00:26", Exitcode: 0, Executor: "Native", PID: 53567, Schedule: "2025-11-10 01:00:00", ScheduleType: "hourly"},
	}

	// Apply offset and limit
	start := offset
	if start >= len(allJobs) {
		return []Job{}
	}

	end := start + limit
	if end > len(allJobs) {
		end = len(allJobs)
	}

	return allJobs[start:end]
} // GetCommandHelp runs "multiflexi-cli <command> --help" and returns the output
func GetCommandHelp(commandName string) (string, error) {
	cmd := exec.Command("multiflexi-cli", commandName, "--help")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run multiflexi-cli %s --help: %w", commandName, err)
	}

	return strings.TrimSpace(string(output)), nil
}
