package ui

import "github.com/VitexSoftware/multiflexi-tui/internal/cli"

// --- Navigation messages ---

// BackMsg requests navigation back to the previous view
type BackMsg struct{}

// ShowHelpMsg requests showing help for a CLI command
type ShowHelpMsg struct {
	Command string
}

// ShowMenuMsg requests showing the command menu
type ShowMenuMsg struct{}

// StatusMessage displays a transient message in the footer
type StatusMessage struct {
	Text string
}

// --- Detail view messages ---

// OpenDetailMsg requests opening a detail view for any entity
type OpenDetailMsg struct {
	Data interface{}
}

// EditItemMsg requests opening an editor for an item
type EditItemMsg struct {
	Data interface{}
}

// ScheduleItemMsg requests opening a scheduler for an item
type ScheduleItemMsg struct {
	Data interface{}
}

// DeleteItemMsg requests deletion (will show confirmation first)
type DeleteItemMsg struct {
	Data  interface{}
	Label string
}

// --- Save messages (one per editable entity) ---

type SaveApplicationMsg struct{ App cli.Application }
type SaveJobMsg struct{ Job cli.Job }
type SaveCompanyMsg struct{ Company cli.Company }
type SaveRunTemplateMsg struct{ Template cli.RunTemplate }

// --- Confirm dialog messages ---

type ConfirmDeleteYesMsg struct{ Data interface{} }
type ConfirmDeleteNoMsg struct{}
