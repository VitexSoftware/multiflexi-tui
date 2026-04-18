package cli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// Client abstracts multiflexi-cli operations for testability.
type Client interface {
	// RunRaw executes multiflexi-cli with raw args, returns stdout.
	RunRaw(args ...string) ([]byte, error)

	// List fetches a paginated list. Target must be a pointer to a slice.
	List(entity string, limit, offset int, target interface{}) error

	// Get fetches a single entity by ID. Target must be a pointer to a struct.
	Get(entity string, id int, target interface{}) error

	// Create runs a create action and returns the raw JSON response.
	Create(entity string, args ...string) ([]byte, error)

	// Update runs an update action with the given args.
	Update(entity string, args ...string) error

	// Delete runs a delete/remove action for the entity with the given ID.
	// deleteAction should be "delete" or "remove" (varies per entity).
	Delete(entity string, deleteAction string, id int) error

	// GetStatus returns system status.
	GetStatus() (*StatusInfo, error)

	// GetCommands returns the list of available CLI commands.
	GetCommands() ([]Command, error)

	// GetCommandHelp returns help text for a specific command.
	GetCommandHelp(name string) (string, error)

	// LastCmd returns the most recently executed CLI command string (for debug display).
	LastCmd() string
}

// CLIClient implements Client using exec.Command("multiflexi-cli", ...).
type CLIClient struct {
	Binary  string // path to multiflexi-cli binary; defaults to "multiflexi-cli"
	lastCmd string // last executed command string, for debug display
}

// NewCLIClient creates a CLIClient with the default binary name.
func NewCLIClient() *CLIClient {
	return &CLIClient{Binary: "multiflexi-cli"}
}

func (c *CLIClient) binary() string {
	if c.Binary != "" {
		return c.Binary
	}
	return "multiflexi-cli"
}

func (c *CLIClient) RunRaw(args ...string) ([]byte, error) {
	c.lastCmd = c.binary() + " " + strings.Join(args, " ")
	cmd := exec.Command(c.binary(), args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			c.lastCmd += "  [exit " + fmt.Sprintf("%d", exitErr.ExitCode()) + "] " + strings.TrimSpace(string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("multiflexi-cli %s: %w", strings.Join(args, " "), err)
	}
	return out, nil
}

func (c *CLIClient) LastCmd() string { return c.lastCmd }

func (c *CLIClient) List(entity string, limit, offset int, target interface{}) error {
	output, err := c.RunRaw(entity, "list",
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

func (c *CLIClient) Get(entity string, id int, target interface{}) error {
	output, err := c.RunRaw(entity, "get",
		"--format=json",
		fmt.Sprintf("--id=%d", id),
	)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(output, target); err != nil {
		return fmt.Errorf("parse %s get JSON: %w", entity, err)
	}
	return nil
}

func (c *CLIClient) Create(entity string, args ...string) ([]byte, error) {
	fullArgs := append([]string{entity, "create", "--format=json"}, args...)
	return c.RunRaw(fullArgs...)
}

func (c *CLIClient) Update(entity string, args ...string) error {
	fullArgs := append([]string{entity, "update", "--format=json"}, args...)
	_, err := c.RunRaw(fullArgs...)
	return err
}

func (c *CLIClient) Delete(entity string, deleteAction string, id int) error {
	_, err := c.RunRaw(entity, deleteAction, "--format=json", "--id", fmt.Sprintf("%d", id))
	return err
}

func (c *CLIClient) GetStatus() (*StatusInfo, error) {
	output, err := c.RunRaw("status", "--format=json")
	if err != nil {
		return nil, err
	}
	status := &StatusInfo{}
	if err := json.Unmarshal(output, status); err != nil {
		return nil, fmt.Errorf("parse status JSON: %w", err)
	}
	return status, nil
}

func (c *CLIClient) GetCommands() ([]Command, error) {
	output, err := c.RunRaw("describe")
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

func (c *CLIClient) GetCommandHelp(name string) (string, error) {
	output, err := c.RunRaw(name, "--help")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
