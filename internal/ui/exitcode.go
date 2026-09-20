package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Exit-code semantic colours mirror multiflexi-web5 ExitCode.php / .mf-exit-*:
// success green, secondary blue, danger red, warning orange, info cyan.
var (
	exitRenderer = func() *lipgloss.Renderer {
		r := lipgloss.NewRenderer(os.Stdout)
		// Match web CSS hex colours; force a colourful profile so badges stay
		// readable even when stdout colour detection is conservative.
		r.SetColorProfile(termenv.TrueColor)
		return r
	}()

	exitSuccessStyle = exitRenderer.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#146c43")).
				Background(lipgloss.Color("#d1e7dd"))
	exitSecondaryStyle = exitRenderer.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#0a4275")).
				Background(lipgloss.Color("#cfe2ff"))
	exitDangerStyle = exitRenderer.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#b02a37")).
			Background(lipgloss.Color("#f8d7da"))
	exitWarningStyle = exitRenderer.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#985700")).
				Background(lipgloss.Color("#fff3cd"))
	exitInfoStyle = exitRenderer.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#087990")).
			Background(lipgloss.Color("#cff4fc"))
)

// ExitCodeStatus maps an exit code to the Bootstrap semantic state used by
// MultiFlexi\Ui\ExitCode::status().
func ExitCodeStatus(code *int) string {
	if code == nil {
		return "info" // not finished yet
	}
	switch *code {
	case 0:
		return "success"
	case -1:
		return "secondary"
	case 75, 127: // EX_TEMPFAIL / command not found
		return "warning"
	default:
		return "danger"
	}
}

func exitCodeStyle(status string) lipgloss.Style {
	switch status {
	case "success":
		return exitSuccessStyle
	case "secondary":
		return exitSecondaryStyle
	case "warning":
		return exitWarningStyle
	case "danger":
		return exitDangerStyle
	default:
		return exitInfoStyle
	}
}

// FormatExitCode renders a job exit code like the web ExitCode widget:
// pending → ⏳ (info); otherwise state emoji + coloured number.
func FormatExitCode(code *int) string {
	status := ExitCodeStatus(code)
	style := exitCodeStyle(status)
	if code == nil {
		return style.Render(" ⏳ ")
	}
	var emoji string
	switch status {
	case "success":
		emoji = "✅"
	case "secondary":
		emoji = "💤"
	case "warning":
		emoji = "⚠️"
	case "danger":
		emoji = "❌"
	default:
		emoji = "⏳"
	}
	return style.Render(fmt.Sprintf(" %s %d ", emoji, *code))
}

// FormatExitCodePlain is a colourless label for detail panes / logs.
func FormatExitCodePlain(code *int) string {
	if code == nil {
		return "⏳ (pending)"
	}
	status := ExitCodeStatus(code)
	var emoji string
	switch status {
	case "success":
		emoji = "✅"
	case "secondary":
		emoji = "💤"
	case "warning":
		emoji = "⚠️"
	case "danger":
		emoji = "❌"
	default:
		emoji = "⏳"
	}
	return fmt.Sprintf("%s %d", emoji, *code)
}

func hasANSI(s string) bool {
	return strings.Contains(s, "\x1b[")
}
