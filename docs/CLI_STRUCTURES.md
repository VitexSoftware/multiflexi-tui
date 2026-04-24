# MultiFlexi CLI — Data Structures Reference

All structures in `internal/cli/types.go` map directly to `multiflexi-cli <entity> list --format=json` output.

## Entity Structures

### Application
```go
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
```

### Company
```go
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
```

### RunTemplate
```go
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
```

### Job
```go
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
```

### Credential
```go
type Credential struct {
    ID               int    `json:"id"`
    Name             string `json:"name"`
    CompanyID        int    `json:"company_id"`
    CredentialTypeID int    `json:"credential_type_id"`
}
```

### Token
```go
type Token struct {
    ID    int    `json:"id"`
    User  string `json:"user"`
    Token string `json:"token"`
}
```

### User
```go
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
```

### Artifact
```go
type Artifact struct {
    ID          int    `json:"id"`
    JobID       int    `json:"job_id"`
    Filename    string `json:"filename"`
    ContentType string `json:"content_type"`
    Artifact    string `json:"artifact"`
    CreatedAt   string `json:"created_at"`
    Note        string `json:"note"`
}
```

### CredType
```go
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
```

### CrPrototype
```go
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
```

### CompanyApp
```go
type CompanyApp struct {
    ID        int `json:"id"`
    CompanyID int `json:"company_id"`
    AppID     int `json:"app_id"`
}
```

### Queue
```go
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
```

### EventSource
```go
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
```

### EventRule
```go
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
```

### StatusInfo
```go
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
```

## CLI Integration Patterns

### List command pattern
```bash
multiflexi-cli <entity> list --format=json --limit=N --offset=N
```

### Create / update command pattern
```bash
multiflexi-cli <entity> create --name="..." --email="..."
multiflexi-cli <entity> update --id=N --name="..."
```

### Delete command pattern
```bash
multiflexi-cli <entity> delete --id=N   # or: remove --id=N
```

### CompanyApp assign / unassign pattern
```bash
# Assign application to company (creates a default RunTemplate)
multiflexi-cli companyapp assign --company_id=N --app_id=M --format=json

# Unassign application from company (removes RunTemplates and actionconfigs)
multiflexi-cli companyapp unassign --company_id=N --app_id=M --format=json

# List RunTemplates for a specific company-app pair
multiflexi-cli companyapp list --company_id=N --app_id=M --format=json
```

### Validate JSON output
```bash
multiflexi-cli <entity> list --format=json | jq .
```
