package cli

// StatusInfo represents comprehensive system status.
type StatusInfo struct {
	VersionCli      string `json:"version-cli"`
	DbMigration     string `json:"db-migration"`
	User            string `json:"user"`
	PHP             string `json:"php"`
	OS              string `json:"os"`
	Memory          int    `json:"memory"`
	Companies       int    `json:"companies"`
	Apps            int    `json:"apps"`
	RunTemplates    int    `json:"runtemplates"`
	Topics          int    `json:"topics"`
	Credentials     int    `json:"credentials"`
	CredentialTypes int    `json:"credential_types"`
	Jobs            string `json:"jobs"`
	Database        string `json:"database"`
	Encryption      string `json:"encryption"`
	Zabbix          string `json:"zabbix"`
	Telemetry       string `json:"telemetry"`
	Executor        string `json:"executor"`
	Scheduler       string `json:"scheduler"`
	Timestamp       string `json:"timestamp"`
}

// Command represents a CLI command from multiflexi-cli describe.
type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CommandInfo is the raw structure from multiflexi-cli describe.
type CommandInfo struct {
	Description string      `json:"description"`
	Arguments   interface{} `json:"arguments,omitempty"`
	Options     interface{} `json:"options,omitempty"`
}

// Application represents an application.
type Application struct {
	ID           int     `json:"id"`
	Enabled      int     `json:"enabled"`
	Image        *string `json:"image"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Executable   string  `json:"executable"`
	DatCreate    string  `json:"DatCreate"`
	DatUpdate    string  `json:"DatUpdate"`
	Setup        string  `json:"setup"`
	CmdParams    string  `json:"cmdparams"`
	Deploy       string  `json:"deploy"`
	Homepage     string  `json:"homepage"`
	Requirements string  `json:"requirements"`
	OciImage     string  `json:"ociimage"`
	Version      string  `json:"version"`
	Code         string  `json:"code"`
	UUID         string  `json:"uuid"`
	Topics       string  `json:"topics"`
	ResultFile   string  `json:"resultfile"`
	Artifacts    string  `json:"artifacts"`
}

// Company represents a company.
type Company struct {
	ID                int     `json:"id"`
	Enabled           int     `json:"enabled"`
	Settings          *string `json:"settings"`
	Logo              *string `json:"logo"`
	Server            int     `json:"server"`
	Name              string  `json:"name"`
	IC                string  `json:"ic"`
	Setup             int     `json:"setup"`
	DatCreate         string  `json:"DatCreate"`
	DatUpdate         string  `json:"DatUpdate"`
	Customer          *int    `json:"customer"`
	Email             string  `json:"email"`
	Slug              string  `json:"slug"`
	ZabbixHost        *string `json:"zabbix_host"`
	RetentionUntil    *string `json:"retention_until"`
	MarkedForDeletion int     `json:"marked_for_deletion"`
}

// RunTemplate represents a run template.
type RunTemplate struct {
	ID                  int     `json:"id"`
	AppID               int     `json:"app_id"`
	CompanyID           int     `json:"company_id"`
	Interv              string  `json:"interv"`
	Prepared            *int    `json:"prepared"`
	Success             *string `json:"success"`
	Fail                *string `json:"fail"`
	Name                string  `json:"name"`
	Delay               int     `json:"delay"`
	Executor            string  `json:"executor"`
	Active              int     `json:"active"`
	Cron                string  `json:"cron"`
	LastSchedule        *string `json:"last_schedule"`
	NextSchedule        *string `json:"next_schedule"`
	Note                *string `json:"note"`
	DatCreate           string  `json:"DatCreate"`
	DatSave             string  `json:"DatSave"`
	SuccessfulJobsCount int     `json:"successfull_jobs_count"`
	FailedJobsCount     int     `json:"failed_jobs_count"`
}

// Job represents a job.
type Job struct {
	ID                int               `json:"id"`
	AppID             int               `json:"app_id"`
	Begin             string            `json:"begin"`
	End               string            `json:"end"`
	CompanyID         int               `json:"company_id"`
	Exitcode          int               `json:"exitcode"`
	Stdout            string            `json:"stdout"`
	Stderr            string            `json:"stderr"`
	LaunchedBy        int               `json:"launched_by"`
	Env               map[string]string `json:"env"`
	Command           string            `json:"command"`
	Schedule          string            `json:"schedule"`
	Executor          string            `json:"executor"`
	RunTemplateID     int               `json:"runtemplate_id"`
	AppVersion        string            `json:"app_version"`
	ScheduleType      string            `json:"schedule_type"`
	PID               int               `json:"pid"`
	RetentionUntil    *string           `json:"retention_until"`
	MarkedForDeletion int               `json:"marked_for_deletion"`
}

// Credential represents a credential.
type Credential struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	CompanyID        int    `json:"company_id"`
	CredentialTypeID int    `json:"credential_type_id"`
}

// Token represents a token.
type Token struct {
	ID    int    `json:"id"`
	User  string `json:"user"`
	Token string `json:"token"`
}

// User represents a user.
type User struct {
	ID                  int     `json:"id"`
	Enabled             int     `json:"enabled"`
	Settings            *string `json:"settings"`
	Email               string  `json:"email"`
	Firstname           string  `json:"firstname"`
	Lastname            string  `json:"lastname"`
	Password            string  `json:"password"`
	PasswordChangedAt   *string `json:"password_changed_at"`
	PasswordExpiresAt   *string `json:"password_expires_at"`
	FailedLoginAttempts int     `json:"failed_login_attempts"`
	LockedUntil         *string `json:"locked_until"`
	TwoFactorEnabled    int     `json:"two_factor_enabled"`
	LastLoginIP         *string `json:"last_login_ip"`
	LastLoginAt         *string `json:"last_login_at"`
	SecuritySettings    *string `json:"security_settings"`
	Login               string  `json:"login"`
	DatCreate           string  `json:"DatCreate"`
	DatSave             string  `json:"DatSave"`
	LastModifierID      *int    `json:"last_modifier_id"`
	DeletedAt           *string `json:"deleted_at"`
	DeletionReason      *string `json:"deletion_reason"`
	AnonymizedAt        *string `json:"anonymized_at"`
	LastActivityAt      *string `json:"last_activity_at"`
	InactiveSince       *string `json:"inactive_since"`
	RetentionUntil      *string `json:"retention_until"`
}

// Artifact represents an artifact.
type Artifact struct {
	ID          int    `json:"id"`
	JobID       int    `json:"job_id"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Artifact    string `json:"artifact"`
	CreatedAt   string `json:"created_at"`
	Note        string `json:"note"`
}

// CredType represents a credential type.
type CredType struct {
	ID        int    `json:"id"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Class     string `json:"class"`
	CompanyID int    `json:"company_id"`
	Logo      string `json:"logo"`
	URL       string `json:"url"`
	Version   int    `json:"version"`
}

// CrPrototype represents a credential prototype.
type CrPrototype struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	URL         string `json:"url"`
	Version     string `json:"version"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CompanyApp represents a company-application relation.
type CompanyApp struct {
	ID        int `json:"id"`
	CompanyID int `json:"company_id"`
	AppID     int `json:"app_id"`
}

// Queue represents a queue item.
type Queue struct {
	ID              int    `json:"id"`
	Job             int    `json:"job"`
	ScheduleType    string `json:"schedule_type"`
	RunTemplateID   int    `json:"runtemplate_id"`
	RunTemplateName string `json:"runtemplate_name"`
	AppID           int    `json:"app_id"`
	AppName         string `json:"app_name"`
	CompanyID       int    `json:"company_id"`
	CompanyName     string `json:"company_name"`
	After           string `json:"after"`
}

// EventSource represents an event source (webhook adapter DB connection).
type EventSource struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	AdapterType  string `json:"adapter_type"`
	DbConnection string `json:"db_connection"`
	DbHost       string `json:"db_host"`
	DbPort       string `json:"db_port"`
	DbDatabase   string `json:"db_database"`
	DbUsername   string `json:"db_username"`
	DbPassword   string `json:"db_password"`
	PollInterval int    `json:"poll_interval"`
	Enabled      int    `json:"enabled"`
}

// EventRule represents an event-to-RunTemplate mapping.
type EventRule struct {
	ID            int    `json:"id"`
	EventSourceID int    `json:"event_source_id"`
	Evidence      string `json:"evidence"`
	Operation     string `json:"operation"`
	RunTemplateID int    `json:"runtemplate_id"`
	Priority      int    `json:"priority"`
	Enabled       int    `json:"enabled"`
	EnvMapping    string `json:"env_mapping"`
}
