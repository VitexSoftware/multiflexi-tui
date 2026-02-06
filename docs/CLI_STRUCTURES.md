# MultiFlexi CLI - Data Structures Reference

This document provides comprehensive reference for all CLI data structures supported by the MultiFlexi TUI.

## Entity Structures

### Enhanced Artifacts Entity

**CLI Command**: `multiflexi-cli artifact list --format=json`

```go
type Artifact struct {
    ID          int     `json:"id"`           // Unique artifact identifier
    JobID       int     `json:"job_id"`       // Associated job ID
    Filename    string  `json:"filename"`     // Original filename
    ContentType string  `json:"content_type"` // MIME type (e.g., "text/plain")
    Artifact    string  `json:"artifact"`     // Artifact content/path
    CreatedAt   string  `json:"created_at"`   // Creation timestamp
    Note        *string `json:"note,omitempty"` // Optional note (nullable)
}
```

**Field Coverage**: 7/7 fields ✅ Complete
**Recent Updates**: Added Filename, ContentType, Artifact, CreatedAt, Note fields

### Enhanced CredTypes Entity

**CLI Command**: `multiflexi-cli credentialtype list --format=json`

```go
type CredType struct {
    ID        int     `json:"id"`           // Unique credential type ID
    Name      string  `json:"name"`         // Credential type name
    UUID      string  `json:"uuid"`         // Unique identifier
    Class     string  `json:"class"`        // Credential class (e.g., "Office365")
    CompanyID *int    `json:"company_id,omitempty"` // Associated company ID (nullable)
    Logo      *string `json:"logo,omitempty"`       // Logo URL/path (nullable)
    URL       *string `json:"url,omitempty"`        // Associated URL (nullable)
    Version   int     `json:"version"`              // Version number
}
```

**Field Coverage**: 8/8 fields ✅ Complete
**Recent Updates**: Added Class, CompanyID, Logo, URL, Version fields

### New CrPrototypes Entity

**CLI Command**: `multiflexi-cli crprototype list --format=json`

```go
type CrPrototype struct {
    ID          int     `json:"id"`           // Unique prototype ID
    Name        string  `json:"name"`         // Prototype name
    Version     string  `json:"version"`      // Version string (e.g., "1.0.0")
    Description *string `json:"description,omitempty"` // Description (nullable)
    Fields      *string `json:"fields,omitempty"`      // Field definitions (nullable)
    Connections *string `json:"connections,omitempty"` // Connection info (nullable)
    CreatedAt   *string `json:"created_at,omitempty"`  // Creation timestamp (nullable)
    UpdatedAt   *string `json:"updated_at,omitempty"`  // Update timestamp (nullable)
    Author      *string `json:"author,omitempty"`      // Author info (nullable)
    Category    *string `json:"category,omitempty"`    // Category (nullable)
}
```

**Field Coverage**: 10/10 fields ✅ Complete
**Status**: New entity - previously missing from TUI

## Complete Entity Reference

### Applications
```go
type Application struct {
    ID          int    `json:"id"`
    Enabled     int    `json:"enabled"`
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

### Companies
```go
type Company struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Code    string `json:"code"`
    Email   string `json:"email"`
    Address string `json:"address"`
    Phone   string `json:"phone"`
}
```

### Jobs
```go
type Job struct {
    ID        int    `json:"id"`
    Command   string `json:"command"`
    Status    string `json:"status"`
    Schedule  string `json:"schedule"`
    StartTime string `json:"start_time"`
    EndTime   string `json:"end_time"`
}
```

### Users
```go
type User struct {
    ID       int    `json:"id"`
    Login    string `json:"login"`
    Email    string `json:"email"`
    Enabled  int    `json:"enabled"`
    LastSeen string `json:"last_seen"`
}
```

### RunTemplates
```go
type RunTemplate struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Command     string `json:"command"`
    Schedule    string `json:"schedule"`
}
```

### Credentials
```go
type Credential struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Username string `json:"username"`
    Password string `json:"password"`
    Host     string `json:"host"`
}
```

### Tokens
```go
type Token struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Token     string `json:"token"`
    ExpiresAt string `json:"expires_at"`
}
```

### CompanyApps
```go
type CompanyApp struct {
    ID          int `json:"id"`
    CompanyID   int `json:"company_id"`
    ApplicationID int `json:"application_id"`
    Enabled     int `json:"enabled"`
}
```

### Encryption
```go
type EncryptionStatus struct {
    Status string `json:"status"`
    Keys   int    `json:"keys"`
}
```

### Queue
```go
type QueueItem struct {
    ID       int    `json:"id"`
    JobID    int    `json:"job_id"`
    Status   string `json:"status"`
    Priority int    `json:"priority"`
}
```

### System Status
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

### Standard List Function Pattern
```go
func GetEntityName(limit, offset int) ([]EntityType, error) {
    cmd := exec.Command("multiflexi-cli", "entityname", "list",
        "--format=json",
        "--order=D", // Newer on top
        "--limit="+fmt.Sprintf("%d", limit),
        "--offset="+fmt.Sprintf("%d", offset),
    )
    
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("failed to get entityname: %w", err)
    }
    
    var entities []EntityType
    if err := json.Unmarshal(output, &entities); err != nil {
        return nil, fmt.Errorf("failed to parse entityname JSON output: %w", err)
    }
    
    return entities, nil
}
```

### Nullable Field Handling
- Use pointers for optional fields: `*string`, `*int`
- JSON tags with `omitempty` for nullable fields
- Proper nil checking in UI display logic

### Error Handling Best Practices
- Always check CLI command execution errors
- Validate JSON parsing with descriptive error messages
- Graceful fallback for missing or invalid data
- User-friendly error display in TUI

## UI Integration Patterns

### Standard Model Structure
```go
type EntityModel struct {
    entities []cli.EntityType
    cursor   int    // For row selection
    offset   int    // For pagination
    limit    int    // Default 10
    loading  bool
    err      error
}
```

### Pagination Controls
- `[← Previous] [Next →] [r] Refresh`
- Consistent 10-row default with configurable limits
- Newer records on top (--order=D)
- Error handling for edge cases

### View Rendering
- Fixed-width column formatting
- Cursor highlighting with `→` indicator
- Consistent header styling
- Proper field truncation for long values

## Version Compatibility

**Current CLI Version**: v2.3.2.110
**TUI Version**: Compatible with all CLI v2.3.x releases
**Last Updated**: February 2026

**Breaking Changes Handled**:
- CredType Version field type (string → int)
- Enhanced field coverage for Artifacts and CredTypes
- New CrPrototype entity integration

## Testing and Validation

### CLI Structure Testing
```bash
# Test individual entity parsing
multiflexi-cli artifact list --format=json --limit=1
multiflexi-cli credentialtype list --format=json --limit=1  
multiflexi-cli crprototype list --format=json --limit=1

# Validate JSON structure
multiflexi-cli <entity> list --format=json | jq .
```

### Integration Testing
```bash
# Build and test TUI
go build ./cmd/multiflexi-tui
./multiflexi-tui

# Navigate to each entity view
# Verify pagination controls work
# Test refresh functionality
```

This comprehensive reference ensures maintainers can easily understand and extend the CLI integration layer as MultiFlexi evolves.