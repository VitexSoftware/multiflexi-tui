package ui

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
)

// EncryptionModel represents the encryption view
type EncryptionModel struct {
	status       string
	err          error
	initializing bool
}

// encryptionInitializedMsg is sent when encryption is initialized successfully
type encryptionInitializedMsg struct{}

// encryptionErrorMsg is sent when there's an error initializing encryption
type encryptionErrorMsg struct {
	err error
}

// NewEncryptionModel creates a new encryption model
func NewEncryptionModel(status string) EncryptionModel {
	return EncryptionModel{
		status: status,
	}
}

// Init initializes the encryption model
func (m EncryptionModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the encryption model
func (m EncryptionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case encryptionInitializedMsg:
		m.initializing = false
		m.status = "active"
		return m, nil

	case encryptionErrorMsg:
		m.initializing = false
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.initializing && m.status != "active" {
				m.initializing = true
				return m, m.initEncryptionCmd()
			}
		}
	}

	return m, nil
}

// View renders the encryption view
func (m EncryptionModel) View() string {
	var content strings.Builder

	content.WriteString(GetTitleStyle().Render("Encryption Status"))
	content.WriteString("\n\n")

	if m.err != nil {
		content.WriteString(GetErrorStyle().Render(fmt.Sprintf("Error: %v", m.err)))
		content.WriteString("\n\n")
	}

	content.WriteString(fmt.Sprintf("Current Status: %s\n\n", m.status))

	if m.initializing {
		content.WriteString("Initializing encryption...")
	} else if m.status != "active" {
		content.WriteString(GetButtonStyle().Render("Initialize Encryption"))
	}

	return content.String()
}

// initEncryptionCmd returns a command that initializes encryption
func (m EncryptionModel) initEncryptionCmd() tea.Cmd {
	return func() tea.Msg {
		err := cli.InitEncryption()
		if err != nil {
			return encryptionErrorMsg{err: err}
		}
		return encryptionInitializedMsg{}
	}
}
