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
	Description string `json:"description"`
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

// GetCommandHelp runs "multiflexi-cli <command> --help" and returns the output
func GetCommandHelp(commandName string) (string, error) {
	cmd := exec.Command("multiflexi-cli", commandName, "--help")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run multiflexi-cli %s --help: %w", commandName, err)
	}

	return strings.TrimSpace(string(output)), nil
}