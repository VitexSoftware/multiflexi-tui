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
	JobsView
	ApplicationsView
	CompaniesView
	CredentialsView
	TokensView
	UsersView
	ArtifactsView
	CredTypesView
	CompanyAppsView
	EncryptionView
	QueueView
	PruneView
	MenuView
	HelpView
	DetailView
	RunTemplateEditorView
	ApplicationEditorView
	RunTemplateSchedulerView
)

// StatusLoadedMsg is sent when status is loaded
type StatusLoadedMsg struct {
	status *cli.StatusInfo
}

// Model represents the main application model
type Model struct {
	state                ViewState
	previousState        ViewState
	jobs                 ui.JobsModel
	runTemplates         ui.RunTemplatesModel
	runTemplateEditor    ui.RunTemplateEditorModel
	runTemplateScheduler ui.RunTemplateSchedulerModel
	applicationEditor    ui.ApplicationEditorModel
	detailView           ui.DetailViewModel
	applications         ui.ApplicationsModel
	companies            ui.CompaniesModel
	credentials          ui.CredentialsModel
	tokens               ui.TokensModel
	users                ui.UsersModel
	artifacts            ui.ArtifactsModel
	credTypes            ui.CredTypesModel
	companyApps          ui.CompanyAppsModel
	encryption           ui.EncryptionModel
	queue                ui.QueueModel
	prune                ui.PruneModel
	menu                 ui.MenuModel
	viewer               ui.ViewerModel
	width                int
	height               int
	statusInfo           *cli.StatusInfo
	statusMessage        string
	menuItems            []string
	menuCursor           int
	selectedHint         string
	focus                bool // true for menu, false for content
}

// NewModel creates and returns a new application model
func NewModel() *Model {
	menuItems := []string{"Status", "RunTemplates", "Jobs", "Applications", "Companies", "Credentials", "Tokens", "Users", "Artifacts", "CredTypes", "CompanyApps", "Encryption", "Queue", "Prune", "Commands", "Help", "Quit"}
	jobs := ui.NewJobsModel()
	runTemplates := ui.NewRunTemplatesModel()
	runTemplateEditor := ui.NewRunTemplateEditorModel(cli.RunTemplate{})
	runTemplateScheduler := ui.NewRunTemplateSchedulerModel(cli.RunTemplate{})
	applicationEditor := ui.NewApplicationEditorModel(cli.Application{})
	detailView := ui.NewDetailViewModel()
	applications := ui.NewApplicationsModel()
	companies := ui.NewCompaniesModel()
	credentials := ui.NewCredentialsModel()
	tokens := ui.NewTokensModel()
	users := ui.NewUsersModel()
	artifacts := ui.NewArtifactsModel()
	credTypes := ui.NewCredTypesModel()
	companyApps := ui.NewCompanyAppsModel()
	encryption := ui.NewEncryptionModel("")
	queue := ui.NewQueueModel()
	prune := ui.NewPruneModel()
	menu := ui.NewMenuModel(nil)
	viewer := ui.NewViewerModel("", "")
	return &Model{
		state:                HomeView,
		jobs:                 jobs,
		runTemplates:         runTemplates,
		runTemplateEditor:    runTemplateEditor,
		runTemplateScheduler: runTemplateScheduler,
		applicationEditor:    applicationEditor,
		detailView:           detailView,
		applications:         applications,
		companies:            companies,
		credentials:          credentials,
		tokens:               tokens,
		users:                users,
		artifacts:            artifacts,
		credTypes:            credTypes,
		companyApps:          companyApps,
		encryption:           encryption,
		queue:                queue,
		prune:                prune,
		menu:                 menu,
		viewer:               viewer,
		menuItems:            menuItems,
		menuCursor:           0,
		selectedHint:         "View system dashboard with status information",
		focus:                true, // focus menu by default
	}
}
