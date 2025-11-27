package app

import (
	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	"github.com/VitexSoftware/multiflexi-tui/internal/ui"
)

// ViewState represents the current view being displayed
type ViewState int

const (
	HomeView ViewState = iota // Default view showing system status
	RunTemplatesView
	RunTemplateDetailView
	JobsView
	ApplicationsView
	CompaniesView
	MenuView
	HelpView
)

// StatusLoadedMsg is sent when status is loaded
type StatusLoadedMsg struct {
	status *cli.StatusInfo
}

// Model represents the main application model
type Model struct {
	state             ViewState
	jobs              ui.JobsModel
	runTemplates      ui.RunTemplatesModel
	runTemplateDetail *ui.DetailWidget
	applications      ui.ApplicationsModel
	companies         ui.CompaniesModel
	menu              ui.MenuModel
	viewer            ui.ViewerModel
	width             int
	height            int
	statusInfo        *cli.StatusInfo
	menuItems         []string
	menuCursor        int
	selectedHint      string
}

// NewModel creates and returns a new application model
func NewModel() *Model {
	menuItems := []string{"Status", "RunTemplates", "Jobs", "Applications", "Companies", "Commands", "Help", "Quit"}
	jobs := ui.NewJobsModel()
	runTemplates := ui.NewRunTemplatesModel()
	runTemplateDetail := &ui.DetailWidget{} // You may want to initialize with config
	applications := ui.NewApplicationsModel()
	companies := ui.NewCompaniesModel()
	menu := ui.NewMenuModel(nil)
	viewer := ui.NewViewerModel("", "")
	return &Model{
		state:             HomeView,
		jobs:              jobs,
		runTemplates:      runTemplates,
		runTemplateDetail: runTemplateDetail,
		applications:      applications,
		companies:         companies,
		menu:              menu,
		viewer:            viewer,
		menuItems:         menuItems,
		menuCursor:        0,
		selectedHint:      "View system dashboard with status information",
	}
}
