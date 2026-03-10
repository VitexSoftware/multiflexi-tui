package cli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// RunCLI executes multiflexi-cli with the given arguments and returns raw output.
func RunCLI(args ...string) ([]byte, error) {
	cmd := exec.Command("multiflexi-cli", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("multiflexi-cli %s: %w", strings.Join(args, " "), err)
	}
	return output, nil
}

// fetchList is a generic helper that runs "multiflexi-cli <entity> list --format=json ..."
// and unmarshals into the provided target slice pointer.
func fetchList(entity string, limit, offset int, target interface{}) error {
	output, err := RunCLI(entity, "list",
		"--format=json",
		"--order=D",
		fmt.Sprintf("--limit=%d", limit),
		fmt.Sprintf("--offset=%d", offset),
	)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(output, target); err != nil {
		return fmt.Errorf("parse %s JSON: %w", entity, err)
	}
	return nil
}

// --- Status ---

func GetStatusInfo() (*StatusInfo, error) {
	output, err := RunCLI("status", "--format=json")
	if err != nil {
		return nil, err
	}
	status := &StatusInfo{}
	if err := json.Unmarshal(output, status); err != nil {
		return nil, fmt.Errorf("parse status JSON: %w", err)
	}
	return status, nil
}

// --- Commands ---

func GetCommands() ([]Command, error) {
	output, err := RunCLI("describe")
	if err != nil {
		return nil, err
	}
	var cmdMap map[string]CommandInfo
	if err := json.Unmarshal(output, &cmdMap); err != nil {
		return nil, fmt.Errorf("parse describe JSON: %w", err)
	}
	var commands []Command
	for name, info := range cmdMap {
		if strings.HasPrefix(name, "_") {
			continue
		}
		commands = append(commands, Command{Name: name, Description: info.Description})
	}
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})
	return commands, nil
}

func GetCommandHelp(commandName string) (string, error) {
	output, err := RunCLI(commandName, "--help")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// --- List operations ---

func GetApplications(limit, offset int) ([]Application, error) {
	var apps []Application
	return apps, fetchList("application", limit, offset, &apps)
}

func GetCompanies(limit, offset int) ([]Company, error) {
	var items []Company
	return items, fetchList("company", limit, offset, &items)
}

func GetRunTemplates(limit, offset int) ([]RunTemplate, error) {
	var items []RunTemplate
	return items, fetchList("runtemplate", limit, offset, &items)
}

func GetJobs(limit, offset int) ([]Job, error) {
	var items []Job
	return items, fetchList("job", limit, offset, &items)
}

func GetCredentials(limit, offset int) ([]Credential, error) {
	var items []Credential
	return items, fetchList("credential", limit, offset, &items)
}

func GetTokens(limit, offset int) ([]Token, error) {
	var items []Token
	return items, fetchList("token", limit, offset, &items)
}

func GetUsers(limit, offset int) ([]User, error) {
	var items []User
	return items, fetchList("user", limit, offset, &items)
}

func GetArtifacts(limit, offset int) ([]Artifact, error) {
	var items []Artifact
	return items, fetchList("artifact", limit, offset, &items)
}

func GetCredTypes(limit, offset int) ([]CredType, error) {
	var items []CredType
	return items, fetchList("credtype", limit, offset, &items)
}

func GetCrPrototypes(limit, offset int) ([]CrPrototype, error) {
	var items []CrPrototype
	return items, fetchList("crprototype", limit, offset, &items)
}

func GetCompanyApps(limit, offset int) ([]CompanyApp, error) {
	var items []CompanyApp
	return items, fetchList("companyapp", limit, offset, &items)
}

func GetQueue(limit, offset int) ([]Queue, error) {
	var items []Queue
	return items, fetchList("queue", limit, offset, &items)
}

// --- Update operations ---

func UpdateApplication(app Application) error {
	_, err := RunCLI("application", "update", fmt.Sprintf("%d", app.ID), "--name", app.Name)
	return err
}

func UpdateRunTemplate(template RunTemplate) error {
	_, err := RunCLI("runtemplate", "update", fmt.Sprintf("%d", template.ID), "--name", template.Name)
	return err
}

func UpdateJob(job Job) error {
	_, err := RunCLI("job", "update",
		"--id", fmt.Sprintf("%d", job.ID),
		"--executor", job.Executor,
		"--schedule_type", job.ScheduleType,
	)
	return err
}

func UpdateCompany(company Company) error {
	_, err := RunCLI("company", "update",
		"--id", fmt.Sprintf("%d", company.ID),
		"--name", company.Name,
		"--email", company.Email,
		"--ic", company.IC,
		"--slug", company.Slug,
	)
	return err
}

// --- Delete operations ---

func DeleteJob(id int) error {
	_, err := RunCLI("job", "delete", "--id", fmt.Sprintf("%d", id))
	return err
}

func DeleteApplication(id int) error {
	_, err := RunCLI("application", "delete", "--id", fmt.Sprintf("%d", id))
	return err
}

func DeleteCompany(id int) error {
	_, err := RunCLI("company", "remove", "--id", fmt.Sprintf("%d", id))
	return err
}

func DeleteRunTemplate(id int) error {
	_, err := RunCLI("runtemplate", "delete", "--id", fmt.Sprintf("%d", id))
	return err
}

// --- Special operations ---

func InitEncryption() error {
	_, err := RunCLI("encryption", "init")
	return err
}

func TruncateQueue() error {
	_, err := RunCLI("queue", "truncate")
	return err
}

func Prune(logs, jobs bool, keep int) error {
	args := []string{"prune"}
	if logs {
		args = append(args, "--logs")
	}
	if jobs {
		args = append(args, "--jobs")
	}
	args = append(args, "--keep", fmt.Sprintf("%d", keep))
	_, err := RunCLI(args...)
	return err
}
