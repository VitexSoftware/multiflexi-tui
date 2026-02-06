



# MultiFlexi CLI - TUI Integration Guide

## Current CLI Coverage Status (February 2026)

### ✅ **100% MultiFlexi CLI Parity Achieved**

The TUI now provides complete coverage of all `multiflexi-cli --format=json` functionality with comprehensive entity management:

#### **Complete Entity Coverage (15 entities)**

| **Entity** | **CLI Command** | **TUI Integration** | **Field Coverage** | **Status** |
|------------|----------------|-------------------|-------------------|------------|
| **Applications** | `application list` | ✅ Full UI | All fields | Complete |
| **Companies** | `company list` | ✅ Full UI | All fields | Complete |
| **Jobs** | `job list` | ✅ Full UI | All fields | Complete |
| **Users** | `user list` | ✅ Full UI | All fields | Complete |
| **RunTemplates** | `runtemplate list` | ✅ Full UI | All fields | Complete |
| **Credentials** | `credential list` | ✅ Full UI | All fields | Complete |
| **Tokens** | `token list` | ✅ Full UI | All fields | Complete |
| **Artifacts** | `artifact list` | ✅ Full UI | **7/7 fields** (Enhanced) | **Updated** |
| **CredTypes** | `credentialtype list` | ✅ Full UI | **8/8 fields** (Enhanced) | **Updated** |
| **CrPrototypes** | `crprototype list` | ✅ Full UI | **10/10 fields** (New) | **Added** |
| **CompanyApps** | `companyapp list` | ✅ Full UI | All fields | Complete |
| **Encryption** | `encryption status` | ✅ Full UI | All fields | Complete |
| **Queue** | `queue list` | ✅ Full UI | All fields | Complete |
| **Prune** | `prune operations` | ✅ Full UI | All operations | Complete |
| **Status** | `status --format=json` | ✅ Full UI | All fields | Complete |

#### **Recent Enhancements (v2.3.2.110 CLI compatibility)**

**🆕 CrPrototypes Entity (New)**
- Complete implementation of credential prototype management
- 10 fields: ID, Name, Version, Description, Fields, Connections, etc.
- Full CRUD operation support through CLI integration

**📈 Enhanced Artifacts Entity (3→7 fields)**
```go
type Artifact struct {
    ID          int     `json:"id"`
    JobID       int     `json:"job_id"`  
    Filename    string  `json:"filename"`      // NEW
    ContentType string  `json:"content_type"`  // NEW
    Artifact    string  `json:"artifact"`      // NEW
    CreatedAt   string  `json:"created_at"`    // NEW
    Note        *string `json:"note,omitempty"`// NEW
}
```

**📈 Enhanced CredTypes Entity (3→8 fields)**  
```go
type CredType struct {
    ID        int     `json:"id"`
    Name      string  `json:"name"`
    UUID      string  `json:"uuid"`
    Class     string  `json:"class"`      // NEW
    CompanyID *int    `json:"company_id,omitempty"` // NEW
    Logo      *string `json:"logo,omitempty"`       // NEW
    URL       *string `json:"url,omitempty"`        // NEW
    Version   int     `json:"version"`              // NEW
}
```

### **Navigation Integration**

All entities are accessible through the main menu system:
```
Status | RunTemplates | Jobs | Applications | Companies | Credentials | 
Tokens | Users | Artifacts | CredTypes | CrPrototypes | CompanyApps | 
Encryption | Queue | Prune | Commands | Help | Quit
```

## Implementation Architecture

### **CLI Layer (`internal/cli/cli.go`)**
- All entity structs with complete field mapping
- JSON unmarshaling with proper nullable field handling  
- Pagination support (limit/offset) for all list operations
- Error handling and CLI command execution

### **UI Layer (`internal/ui/`)**
- Individual model files for each entity type
- Consistent pagination and navigation patterns
- Loading states and error handling
- Cursor management for row selection

### **Application Layer (`internal/app/`)**
- ViewState management for all entities
- Menu integration and navigation handling
- Model initialization and state transitions
- Context-aware hint system

---

# MultiFlexi CLI - Listing Pagination Guide

## Overview
MultiFlexi CLI provides comprehensive pagination and filtering capabilities for all listing commands. You can control how many records to display, skip records, sort results, and customize output fields.

## Code Maintenance Guide: Updating When CLI Output Changes

### When MultiFlexi CLI Schema Changes
When `multiflexi-cli describe` output is updated or new listing columns are added, follow these steps to update the TUI code:

#### 1. Identify Changes in CLI Output
```bash
# Test the updated CLI command to see new structure
multiflexi-cli <entity> list --format=json --limit=1

# Compare with current struct definitions in internal/cli/cli.go
# Look for new fields, changed field names, or different data types
```

#### 2. Update Data Structures in `internal/cli/cli.go`
- **Add new fields** to existing structs with proper JSON tags
- **Update field names** if CLI field names changed  
- **Change data types** if CLI output types changed
- **Add new entity structs** for new listing commands

```go
// Example: Adding new fields to existing struct
type Application struct {
    ID          int    `json:"id"`
    Enabled     int    `json:"enabled"`
    Name        string `json:"name"`
    Description string `json:"description"`
    // NEW FIELD - add with proper JSON tag
    Category    string `json:"category"`
    // NEW FIELD - nullable fields
    LastRun     *string `json:"last_run,omitempty"`
}

// Example: New entity struct
type NewEntity struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Status string `json:"status"`
}
```

#### 3. Add/Update CLI Functions
```go
// Example: New entity CLI function
func GetNewEntities(limit, offset int) ([]NewEntity, error) {
    cmd := exec.Command("multiflexi-cli", "newentity", "list",
        "--format=json",
        "--order=D", // Newer on top
        "--limit="+fmt.Sprintf("%d", limit),
        "--offset="+fmt.Sprintf("%d", offset),
    )
    // ... rest of implementation
}
```

#### 4. Update UI Models in `internal/ui/`
- **Add new model file** for new entities: `internal/ui/newentity.go`
- **Update existing models** to handle new fields
- **Add cursor functionality** for focus management
- **Update View functions** to display new columns

```go
// Example: New UI model structure
type NewEntityModel struct {
    entities []cli.NewEntity
    cursor   int        // For row selection
    offset   int        // For pagination
    limit    int        // Default 10
    loading  bool
    err      error
}

// Update View function with new columns
func (m NewEntityModel) View() string {
    // Add new columns to header
    content.WriteString(fmt.Sprintf("%-5s %-25s %-15s %-10s",
        "ID", "Name", "Category", "Status"))
    
    // Add new fields to row display
    for i, entity := range m.entities {
        highlight := ""
        if i == m.cursor {
            highlight = " → " // Focus indicator
        }
        line := fmt.Sprintf("%s%-5d %-25s %-15s %-10s",
            highlight, entity.ID, entity.Name, entity.Category, entity.Status)
    }
}
```

#### 5. Update Application Architecture in `internal/app/app.go`
- **Add new ViewState** constant for new entities
- **Add model field** to main Model struct  
- **Update menu items** to include new entity
- **Add view handling** in View(), Update(), and handleMenuSelection()

```go
// Add ViewState
const (
    HomeView ViewState = iota
    RunTemplatesView
    JobsView
    ApplicationsView
    CompaniesView
    NewEntityView        // NEW
    MenuView
    HelpView
)

// Add to Model struct
type Model struct {
    // ... existing fields
    newEntities ui.NewEntityModel  // NEW
}

// Update menu items
menuItems := []string{"Status", "RunTemplates", "Jobs", 
    "Applications", "Companies", "NewEntities", "Commands", "Help", "Quit"}

// Add to handleMenuSelection
case 5: // NewEntities
    m.state = NewEntityView
    m.newEntities = ui.NewNewEntityModel()
    return m, m.newEntities.Init()
```

#### 6. Testing and Validation Checklist
- [ ] **CLI Command Test**: Verify new CLI command works and returns expected JSON
- [ ] **Struct Parsing**: Test that Go structs properly unmarshal CLI JSON
- [ ] **UI Display**: Check that new columns display correctly in TUI
- [ ] **Pagination**: Verify pagination works with new entity listings
- [ ] **Focus Navigation**: Test up/down arrow navigation on new tables
- [ ] **Menu Integration**: Confirm new entity is accessible from main menu
- [ ] **Error Handling**: Test behavior when CLI command fails

#### 7. Common Update Patterns

**Pattern 1: Adding Single Column**
1. Add field to struct in `cli.go`
2. Update View function header and row display
3. Test with real CLI data

**Pattern 2: New Entity Type**  
1. Add struct and Get function to `cli.go`
2. Create new UI model file
3. Add ViewState and integrate into app.go
4. Update menu items

**Pattern 3: Changed Field Names**
1. Update JSON tags in existing struct
2. Test that unmarshaling still works
3. Update any hardcoded field references

#### 8. Debugging Tips
```bash
# Test CLI output format
multiflexi-cli <entity> list --format=json | jq .

# Test specific pagination
multiflexi-cli <entity> list --format=json --limit=2 --offset=0

# Validate JSON structure
echo '{"test": "data"}' | jq . # Should be valid JSON
```

#### 9. Code Quality Standards
- **Consistent Naming**: Use same patterns as existing models
- **Error Handling**: Always handle CLI command failures gracefully  
- **Plain Text Display**: Default to plain text columns, avoid styling
- **Pagination**: Support standard 10-row pages with newer-on-top ordering
- **Focus Management**: Include cursor for row selection
- **Help Text**: Update pagination buttons format `[← Previous] [Next →] [r] Refresh`

This systematic approach ensures the TUI stays synchronized with MultiFlexi CLI evolution while maintaining code quality and user experience consistency.

## Available Pagination Options

### 1. Limit Results (`--limit=<number>`)
Restrict the number of records returned:
```bash
# Show only 5 users
multiflexi-cli user list --limit=5

# Display 10 applications
multiflexi-cli application list --limit=10
```

### 2. Offset Records (`--offset=<number>`)
Skip a specified number of records (useful for pagination):
```bash
# Skip first 10 records, show the rest
multiflexi-cli user list --offset=10

# Skip 20 records, then show 5
multiflexi-cli company list --offset=20 --limit=5
```

### 3. Sort Order (`--order=<A|D>`)
Control the sort order of results:
- `A` = Ascending (default)
- `D` = Descending

```bash
# Sort users in descending order
multiflexi-cli user list --order=D

# Get latest 10 jobs (descending order)
multiflexi-cli job list --order=D --limit=10
```

### 4. Select Fields (`--fields=<field1,field2,field3>`)
Display only specific fields in the output:
```bash
# Show only ID and name for applications
multiflexi-cli application list --fields=id,name

# Display specific user fields
multiflexi-cli user list --fields=id,login,email --limit=5
```

## Combining Options for Advanced Pagination

### Page-by-Page Navigation
```bash
# Page 1: First 10 records
multiflexi-cli user list --limit=10 --offset=0

# Page 2: Next 10 records  
multiflexi-cli user list --limit=10 --offset=10

# Page 3: Next 10 records
multiflexi-cli user list --limit=10 --offset=20
```

### Efficient Data Browsing
```bash
# Get latest 5 jobs with key fields only
multiflexi-cli job list --order=D --limit=5 --fields=id,name,status

# Browse companies starting from 6th record
multiflexi-cli company list --offset=5 --limit=10 --order=A

# Get specific application data
multiflexi-cli application list --fields=name,version --limit=20
```

## Output Formatting

### JSON Output for APIs/Scripts
Add `--format=json` for machine-readable output:
```bash
# JSON pagination for scripts
multiflexi-cli user list --limit=5 --offset=10 --format=json

# Structured data with custom fields
multiflexi-cli application list --fields=id,name --format=json
```

### Human-Readable Output (Default)
```bash
# Default text output (no --format needed)
multiflexi-cli user list --limit=5 --order=D
```

## Commands Supporting Pagination

All listing operations support pagination options:

- `multiflexi-cli user list`
- `multiflexi-cli application list`
- `multiflexi-cli company list`
- `multiflexi-cli job list`
- `multiflexi-cli token list`
- `multiflexi-cli runtemplate list`
- `multiflexi-cli credential list`
- `multiflexi-cli credentialtype list`
- `multiflexi-cli artifact list`
- `multiflexi-cli queue list`
- `multiflexi-cli userdataerasure list`

## System Status Command

The `multiflexi-cli status` command provides comprehensive system information:

```bash
# Get system status in JSON format
multiflexi-cli status --format=json
```

**Current Output Format (v2.3.2.110):**
```json
{
    "version-cli": "2.3.2.110",
    "db-migration": "GdprDataExport (20251220151500)",
    "user": "root",
    "php": "8.4.17",
    "os": "Linux",
    "memory": 4259072,
    "companies": 5,
    "apps": 51,
    "runtemplates": 151,
    "topics": 67,
    "credentials": 12,
    "credential_types": 10,
    "jobs": "total: 107, monthly: 107, weekly: 107, daily: 10, hourly: 0, minute avg: 0.01",
    "database": "mysql Localhost via UNIX socket Uptime: 462910  Threads: 5  Questions: 1907901  Slow queries: 10  Opens: 2794  Open tables: 1114  Queries per second avg: 4.121 11.8.3-MariaDB-0+deb13u1 from Debian",
    "encryption": "active (3 keys)",
    "zabbix": "krax.vitexsoftware.brevnov.czf => zabbix-dev.serverovna.brevnov.czf",
    "telemetry": "disabled",
    "executor": "active",
    "scheduler": "inactive",
    "timestamp": "2026-02-02T22:52:21+00:00"
}
```

**Status Data Structure for TUI:**
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

## Practical Examples

### Dashboard Summary
```bash
# Quick overview: 5 latest jobs
multiflexi-cli job list --order=D --limit=5 --fields=id,name,status

# Recent applications
multiflexi-cli application list --limit=10 --fields=name,version
```

### Data Analysis
```bash
# Export all users to JSON (in batches)
multiflexi-cli user list --format=json --limit=100 --offset=0
multiflexi-cli user list --format=json --limit=100 --offset=100

# Find specific records
multiflexi-cli company list --fields=id,name,email --order=A
```

### Performance Optimization
```bash
# Large datasets: use pagination to avoid memory issues
multiflexi-cli job list --limit=50 --offset=0    # First batch
multiflexi-cli job list --limit=50 --offset=50   # Second batch
multiflexi-cli job list --limit=50 --offset=100  # Third batch
```

## Best Practices

1. **Use `--limit`** for large datasets to improve performance
2. **Combine `--offset` and `--limit`** for efficient pagination
3. **Use `--fields`** to reduce output size and improve readability  
4. **Use `--format=json`** for automated processing and API integration
5. **Use `--order=D`** to get most recent records first

## Tips for TUI Implementation

- **Page Navigation**: Implement next/previous page buttons using offset calculations
- **Dynamic Limits**: Allow users to change page size (10, 25, 50, 100 records)
- **Field Selection**: Provide checkboxes for field selection
- **Sort Toggle**: Click column headers to toggle A/D ordering
- **Search Integration**: Combine pagination with filtering options

---

This guide provides comprehensive information about pagination features that can be implemented in the multiflexi-tui interface for an optimal user experience.