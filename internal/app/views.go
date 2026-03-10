package app

import (
	"fmt"
	"strings"

	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

// renderMenuBar renders the top menu bar with hints and horizontal scrolling
func (m Model) renderMenuBar() string {
	width := m.width
	if width == 0 {
		width = 80
	}

	title := "MultiFlexi TUI"
	titleWidth := len(title) + 4
	availableWidth := width - titleWidth

	var visibleMenuItems []string
	var currentWidth int
	style := ui.GetUnselectedItemStyle()
	if m.focus {
		style = ui.GetSelectedItemStyle()
	}

	for i := m.menuOffset; i < len(m.menuItems); i++ {
		item := m.menuItems[i]
		var renderedItem string
		if i == m.menuCursor && m.focus {
			renderedItem = style.Render(" " + item + " ")
		} else if i == m.activeMenuItem {
			renderedItem = ui.GetActiveMenuItemStyle().Render(" " + item + " ")
		} else {
			renderedItem = ui.GetUnselectedItemStyle().Render(" " + item + " ")
		}

		itemWidth := len(item) + 2 + 1
		if currentWidth+itemWidth > availableWidth && len(visibleMenuItems) > 0 {
			if currentWidth+3 <= availableWidth {
				visibleMenuItems = append(visibleMenuItems, "...")
			}
			break
		}

		visibleMenuItems = append(visibleMenuItems, renderedItem)
		currentWidth += itemWidth
	}

	if m.menuOffset > 0 {
		visibleMenuItems = append([]string{"..."}, visibleMenuItems...)
	}

	menuLine := ui.GetTitleStyle().Render(" "+title+" ") + " " + strings.Join(visibleMenuItems, " ")
	hintLine := ui.GetItemDescriptionStyle().Render(" " + m.selectedHint + " ")
	separator := strings.Repeat("═", width)

	return menuLine + "\n" + hintLine + "\n" + separator + "\n"
}

// renderHelpFooter renders the bottom status/help bar
func (m Model) renderHelpFooter() string {
	width := m.width
	if width == 0 {
		width = 80
	}

	separator := strings.Repeat("═", width)
	var helpLine string
	if m.focus {
		helpLine = ui.GetFooterStyle().Render(" ←/→: navigate menu • enter: select • tab: switch to content • q: quit ")
	} else {
		helpLine = ui.GetFooterStyle().Render(" ↑/↓: rows • ←/→: pages • enter/space: detail • e: editor • tab: menu • q: quit ")
	}

	statusLine := ""
	if m.statusMessage != "" {
		statusLine = ui.GetFooterStyle().Render(" " + m.statusMessage + " ")
		// Mark for clearing on next Update cycle (not in View!)
	}

	return separator + "\n" + statusLine + "\n" + helpLine
}

// renderSystemStatus renders the system status dashboard
func (m Model) renderSystemStatus() string {
	var content strings.Builder

	content.WriteString(ui.GetTitleStyle().Render("🖥️  MultiFlexi System Dashboard"))
	content.WriteString("\n\n")

	if m.statusInfo == nil {
		content.WriteString(ui.GetItemDescriptionStyle().Render("Loading system status..."))
		content.WriteString("\n")
		return content.String()
	}

	rows := []struct {
		icon   string
		label  string
		value  string
		status string
	}{
		{"🔧", "CLI Version", m.statusInfo.VersionCli, "info"},
		{"🗄️", "DB Migration", m.statusInfo.DbMigration, "info"},
		{"👤", "User", m.statusInfo.User, "info"},
		{"🐘", "PHP", m.statusInfo.PHP, "info"},
		{"💻", "OS", m.statusInfo.OS, "info"},
		{"🧠", "Memory", fmt.Sprintf("%d KB", m.statusInfo.Memory), "info"},
		{"🏢", "Companies", fmt.Sprintf("%d", m.statusInfo.Companies), "info"},
		{"📱", "Applications", fmt.Sprintf("%d", m.statusInfo.Apps), "info"},
		{"📄", "RunTemplates", fmt.Sprintf("%d", m.statusInfo.RunTemplates), "info"},
		{"🏷️", "Topics", fmt.Sprintf("%d", m.statusInfo.Topics), "info"},
		{"🔑", "Credentials", fmt.Sprintf("%d", m.statusInfo.Credentials), "info"},
		{"🎭", "Credential Types", fmt.Sprintf("%d", m.statusInfo.CredentialTypes), "info"},
		{"💼", "Jobs", m.statusInfo.Jobs, "info"},
		{"⚙️", "Executor", m.statusInfo.Executor, m.statusInfo.Executor},
		{"📅", "Scheduler", m.statusInfo.Scheduler, m.statusInfo.Scheduler},
		{"🔐", "Encryption", m.statusInfo.Encryption, m.statusInfo.Encryption},
		{"📊", "Zabbix", m.statusInfo.Zabbix, "info"},
		{"📈", "Telemetry", m.statusInfo.Telemetry, m.statusInfo.Telemetry},
		{"🕒", "Timestamp", m.statusInfo.Timestamp, "info"},
	}

	labelWidth := 18
	for _, row := range rows {
		var valueStyle lipgloss.Style
		switch row.status {
		case "active":
			valueStyle = ui.GetActiveStatusStyle()
		case "inactive", "disabled":
			valueStyle = ui.GetDisabledStatusStyle()
		default:
			valueStyle = ui.GetItemDescriptionStyle()
		}

		line := fmt.Sprintf("%s %-*s %s",
			row.icon, labelWidth, row.label+":", valueStyle.Render(row.value))
		content.WriteString(line)
		content.WriteString("\n")
	}

	if m.statusInfo.Database != "" {
		content.WriteString("\n")
		content.WriteString("🗄️  Database Information:")
		content.WriteString("\n")
		dbInfo := m.statusInfo.Database
		if len(dbInfo) > 80 {
			dbInfo = dbInfo[:77] + "..."
		}
		content.WriteString(ui.GetItemDescriptionStyle().Render("   " + dbInfo))
		content.WriteString("\n")
	}

	return content.String()
}
