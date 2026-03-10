package app

import (
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// ViewKind identifies which section is currently active (for menu highlight).
type ViewKind int

const (
	ViewHome ViewKind = iota
	ViewRunTemplates
	ViewJobs
	ViewApplications
	ViewCompanies
	ViewCredentials
	ViewTokens
	ViewUsers
	ViewArtifacts
	ViewCredTypes
	ViewCrPrototypes
	ViewCompanyApps
	ViewEncryption
	ViewQueue
	ViewPrune
	ViewCommands
	ViewHelp
	ViewDetail
	ViewEditor
	ViewConfirmDelete
)

// Model represents the main application model.
type Model struct {
	// Current active sub-view (implements tea.Model)
	activeView tea.Model
	viewKind   ViewKind
	// Previous view for back-navigation from detail/editor
	prevView tea.Model
	prevKind ViewKind

	// Shared state
	detailView    ui.DetailViewModel
	confirmDialog ui.ConfirmDialogModel
	viewer        ui.ViewerModel
	menu          ui.MenuModel

	// Layout
	width  int
	height int

	// System status (loaded once, refreshable)
	statusInfo *cli.StatusInfo

	// Status bar message (transient)
	statusMessage     string
	statusMsgConsumed bool // cleared on next Update cycle

	// Menu bar state
	menuItems      []string
	menuCursor     int
	menuOffset     int
	activeMenuItem int
	selectedHint   string
	focus          bool // true = menu focused, false = content focused
}

// NewModel creates and returns a new application model.
func NewModel() *Model {
	menuItems := []string{
		"Status", "RunTemplates", "Jobs", "Applications", "Companies",
		"Credentials", "Tokens", "Users", "Artifacts", "CredTypes",
		"CrPrototypes", "CompanyApps", "Encryption", "Queue", "Prune",
		"Commands", "Help", "Quit",
	}

	return &Model{
		activeView:   nil, // HomeView renders inline
		viewKind:     ViewHome,
		detailView:   ui.NewDetailViewModel(),
		viewer:       ui.NewViewerModel("", ""),
		menu:         ui.NewMenuModel(nil),
		menuItems:    menuItems,
		menuCursor:   0,
		menuOffset:   0,
		selectedHint: "View system dashboard with status information",
		focus:        true,
	}
}
