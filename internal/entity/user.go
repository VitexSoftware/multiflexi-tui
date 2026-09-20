package entity

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// userRoleListPayload matches multiflexi-cli user-role:list --format=json.
type userRoleListPayload struct {
	UserID          int                      `json:"user_id"`
	Login           string                   `json:"login"`
	Roles           []string                 `json:"roles"`
	RoleDetails     []map[string]interface{} `json:"role_details"`
	AvailableRoles  []string                 `json:"available_roles"`
}

func fetchUserRoles(c cli.Client, userID int) (*userRoleListPayload, []byte, error) {
	out, err := c.RunRaw("user-role:list", "--format=json", "--user_id", fmt.Sprintf("%d", userID))
	if err != nil {
		return nil, out, err
	}
	var payload userRoleListPayload
	if err := json.Unmarshal(out, &payload); err != nil {
		return nil, out, fmt.Errorf("parse user-role:list JSON: %w", err)
	}
	return &payload, out, nil
}

func formatUserRolesView(p *userRoleListPayload) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("User #%d (%s)\n\n", p.UserID, p.Login))
	if len(p.Roles) == 0 {
		b.WriteString("Assigned roles: (none)\n")
	} else {
		b.WriteString("Assigned roles:\n")
		for _, r := range p.Roles {
			b.WriteString("  • " + r + "\n")
		}
	}
	b.WriteString("\nAvailable roles: ")
	if len(p.AvailableRoles) == 0 {
		b.WriteString("(none)\n")
	} else {
		b.WriteString(strings.Join(p.AvailableRoles, ", ") + "\n")
	}
	if len(p.RoleDetails) > 0 {
		b.WriteString("\nDetails:\n")
		raw, _ := json.MarshalIndent(p.RoleDetails, "", "  ")
		b.Write(raw)
		b.WriteString("\n")
	}
	return b.String()
}

var UserDef = &EntityDef{
	Name: "👤 Users", CLIEntity: "user", DeleteAction: "delete", Limit: 10,
	Columns: []ui.TableColumn{
		{Header: "ID", Width: 6, Field: "id"}, {Header: "Login", Width: 20, Field: "login"},
		{Header: "Name", Width: 25, Field: "name"}, {Header: "Email", Width: 30, Field: "email"},
		{Header: "Active", Width: 7, Field: "enabled"},
	},
	Fetch: func(c cli.Client, limit, offset int) ([]ui.TableRow, error) {
		var items []cli.User
		if err := c.List("user", limit, offset, &items); err != nil {
			return nil, err
		}
		rows := make([]ui.TableRow, len(items))
		for i, u := range items {
			enabled := "No"
			if u.Enabled == 1 {
				enabled = "Yes"
			}
			rows[i] = ui.TableRow{ID: u.ID, Values: map[string]string{
				"id": fmt.Sprintf("%d", u.ID), "login": u.Login,
				"name": u.Firstname + " " + u.Lastname, "email": u.Email, "enabled": enabled,
			}, FullData: u}
		}
		return rows, nil
	},
	ToDetail: func(data interface{}) []ui.DetailField {
		u := data.(cli.User)
		lastIP := ""
		if u.LastLoginIP != nil {
			lastIP = *u.LastLoginIP
		}
		lastAt := ""
		if u.LastLoginAt != nil {
			lastAt = *u.LastLoginAt
		}
		return []ui.DetailField{
			{Label: "ID", Value: fmt.Sprintf("%d", u.ID)},
			{Label: "Login", Value: u.Login},
			{Label: "First Name", Value: u.Firstname},
			{Label: "Last Name", Value: u.Lastname},
			{Label: "Email", Value: u.Email},
			{Label: "Enabled", Value: fmt.Sprintf("%d", u.Enabled)},
			{Label: "2FA Enabled", Value: fmt.Sprintf("%d", u.TwoFactorEnabled)},
			{Label: "Last Login IP", Value: lastIP},
			{Label: "Last Login At", Value: lastAt},
			{Label: "Failed Logins", Value: fmt.Sprintf("%d", u.FailedLoginAttempts)},
			{Label: "Created", Value: u.DatCreate},
		}
	},
	ToEditor: func(data interface{}) []ui.EditorField {
		u := data.(cli.User)
		return []ui.EditorField{
			{Label: "Login", Placeholder: "username", Value: u.Login},
			{Label: "First Name", Placeholder: "First name", Value: u.Firstname},
			{Label: "Last Name", Placeholder: "Last name", Value: u.Lastname},
			{Label: "Email", Placeholder: "email@example.com", Value: u.Email},
			{Label: "Password", Placeholder: "leave blank to keep current", Value: ""},
			{Label: "Enabled", Placeholder: "1 or 0", Value: fmt.Sprintf("%d", u.Enabled)},
		}
	},
	UpdateArgs: func(data interface{}, fields map[string]string) []string {
		u := data.(cli.User)
		args := []string{
			"--id", fmt.Sprintf("%d", u.ID),
			"--login", fields["Login"],
			"--firstname", fields["First Name"],
			"--lastname", fields["Last Name"],
			"--email", fields["Email"],
			"--enabled", fields["Enabled"],
		}
		if v := fields["Password"]; v != "" {
			args = append(args, "--plaintext", v)
		}
		return args
	},
	NewFields: func() []ui.EditorField {
		return []ui.EditorField{
			{Label: "Login", Placeholder: "username", Required: true},
			{Label: "Email", Placeholder: "email@example.com", Required: true},
			{Label: "Password", Placeholder: "plaintext password", Required: true},
			{Label: "First Name", Placeholder: "First name"},
			{Label: "Last Name", Placeholder: "Last name"},
		}
	},
	CreateArgs: func(fields map[string]string) []string {
		args := []string{
			"--login", fields["Login"],
			"--email", fields["Email"],
			"--plaintext", fields["Password"],
		}
		if v := fields["First Name"]; v != "" {
			args = append(args, "--firstname", v)
		}
		if v := fields["Last Name"]; v != "" {
			args = append(args, "--lastname", v)
		}
		return args
	},
	GetID:    func(data interface{}) int { return data.(cli.User).ID },
	GetLabel: func(data interface{}) string { return fmt.Sprintf("User: %s", data.(cli.User).Login) },
	Actions: []ui.ActionDef{
		{Label: "Edit", Key: "e", Command: "edit"},
		{
			Label:   "Roles",
			Key:     "r",
			Command: "roles",
			Handler: func(c cli.Client, data interface{}) tea.Cmd {
				u := data.(cli.User)
				return func() tea.Msg {
					payload, raw, err := fetchUserRoles(c, u.ID)
					viewer := ui.NewViewer(fmt.Sprintf("RBAC Roles: %s", u.Login))
					if err != nil {
						viewer.SetContent(fmt.Sprintf("RBAC Roles: %s", u.Login),
							fmt.Sprintf("Error: %v\n\n%s", err, string(raw)))
					} else {
						viewer.SetContent(fmt.Sprintf("RBAC Roles: %s", u.Login), formatUserRolesView(payload))
					}
					return ui.NavigateToMsg{View: viewer}
				}
			},
		},
		{
			Label:   "Set Roles",
			Key:     "R",
			Command: "set-roles",
			Handler: func(c cli.Client, data interface{}) tea.Cmd {
				u := data.(cli.User)
				return func() tea.Msg {
					current := ""
					hint := "admin,viewer"
					if payload, _, err := fetchUserRoles(c, u.ID); err == nil {
						current = strings.Join(payload.Roles, ",")
						if len(payload.AvailableRoles) > 0 {
							hint = strings.Join(payload.AvailableRoles, ",")
						}
					}
					form := NewActionFormView(
						fmt.Sprintf("Set RBAC Roles: %s", u.Login),
						[]ui.EditorField{
							{Label: "Roles", Placeholder: hint, Value: current, Required: false},
							{Label: "Replace", Placeholder: "true=replace all, false=add", Value: "true"},
						},
						func(fields map[string]string) tea.Cmd {
							return func() tea.Msg {
								roles := strings.TrimSpace(fields["Roles"])
								replace := strings.TrimSpace(fields["Replace"])
								if replace == "" {
									replace = "true"
								}
								args := []string{
									"user-role:set", "--format=json",
									"--user_id", fmt.Sprintf("%d", u.ID),
									"--roles", roles,
									"--replace", replace,
								}
								out, err := c.RunRaw(args...)
								viewer := ui.NewViewer(fmt.Sprintf("Set Roles: %s", u.Login))
								viewer.RefreshOnBack = true
								if err != nil {
									viewer.SetContent("Set Roles Result", fmt.Sprintf("Error: %v\n\n%s", err, string(out)))
								} else {
									viewer.SetContent("Set Roles Result", string(out))
								}
								return ui.NavigateToMsg{View: viewer}
							}
						},
					)
					return ui.NavigateToMsg{View: form}
				}
			},
		},
		{
			Label:   "Make Admin",
			Key:     "a",
			Command: "promote-admin",
			Confirm: "Add the 'admin' RBAC role to this user (keeps existing roles)?",
			Handler: func(c cli.Client, data interface{}) tea.Cmd {
				u := data.(cli.User)
				return func() tea.Msg {
					out, err := c.RunRaw(
						"user-role:set", "--format=json",
						"--user_id", fmt.Sprintf("%d", u.ID),
						"--roles", "admin",
						"--replace", "false",
					)
					viewer := ui.NewViewer(fmt.Sprintf("Promote Admin: %s", u.Login))
					viewer.RefreshOnBack = true
					if err != nil {
						viewer.SetContent("Promote Admin", fmt.Sprintf("Error: %v\n\n%s", err, string(out)))
					} else {
						viewer.SetContent("Promote Admin", string(out))
					}
					return ui.NavigateToMsg{View: viewer}
				}
			},
		},
		{Label: "Delete", Key: "d", Command: "delete"},
	},
}

func init() {
	Register(Entry{Label: "Users", Hint: "Manage users • r: roles • R: set roles • a: make admin", Def: UserDef})
}
