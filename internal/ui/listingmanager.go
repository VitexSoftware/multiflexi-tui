package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
)

// ListingManager provides enhanced data management for UI listing views
type ListingManager struct {
	fetcher      *cli.DataFetcher
	refreshTimer *time.Timer
}

// ListingConfig represents configuration for a listing view
type ListingConfig struct {
	EntityType      string        // Type of entity (e.g., "application", "company")
	RefreshInterval time.Duration // Auto-refresh interval (0 = disabled)
	CacheTTL        time.Duration // Cache time-to-live
	DefaultLimit    int           // Default page size
	MaxLimit        int           // Maximum page size
	Timeout         time.Duration // Request timeout
}

// ListingState represents the current state of a listing
type ListingState struct {
	Loading     bool
	Error       error
	Data        interface{}
	CurrentPage int
	Limit       int
	Offset      int
	TotalCount  int
	HasMore     bool
	HasPrev     bool
	LastUpdate  time.Time
	Cached      bool
	FetchTime   time.Duration
}

// DataLoadedMsg represents a successful data load
type DataLoadedMsg struct {
	EntityType string
	Data       interface{}
	State      ListingState
}

// DataErrorMsg represents a data loading error
type DataErrorMsg struct {
	EntityType string
	Error      error
}

// RefreshDataMsg represents a refresh request
type RefreshDataMsg struct {
	EntityType string
	Force      bool // Force refresh even if cached
}

// NewListingManager creates a new listing manager
func NewListingManager() *ListingManager {
	return &ListingManager{
		fetcher: cli.NewDataFetcher(),
	}
}

// LoadData loads data for a specific entity type with enhanced error handling
func (lm *ListingManager) LoadData(config ListingConfig, limit, offset int) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
		defer cancel()

		// Create fetch parameters
		params := cli.FetchParams{
			Command:    config.EntityType,
			SubCommand: "list",
			Limit:      limit,
			Offset:     offset,
			CacheTTL:   config.CacheTTL,
			Timeout:    config.Timeout,
		}

		// Determine target type and fetch data
		var result *cli.FetchResult
		var err error

		switch config.EntityType {
		case "application":
			var apps []cli.Application
			result, err = lm.fetcher.FetchData(ctx, params, &apps)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, apps, result, limit, offset)
			}

		case "company":
			var companies []cli.Company
			result, err = lm.fetcher.FetchData(ctx, params, &companies)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, companies, result, limit, offset)
			}

		case "crprototype":
			var crprototypes []cli.CrPrototype
			result, err = lm.fetcher.FetchData(ctx, params, &crprototypes)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, crprototypes, result, limit, offset)
			}

		case "credential":
			var credentials []cli.Credential
			result, err = lm.fetcher.FetchData(ctx, params, &credentials)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, credentials, result, limit, offset)
			}

		case "credentialtype":
			var credtypes []cli.CredType
			result, err = lm.fetcher.FetchData(ctx, params, &credtypes)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, credtypes, result, limit, offset)
			}

		case "token":
			var tokens []cli.Token
			result, err = lm.fetcher.FetchData(ctx, params, &tokens)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, tokens, result, limit, offset)
			}

		case "user":
			var users []cli.User
			result, err = lm.fetcher.FetchData(ctx, params, &users)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, users, result, limit, offset)
			}

		case "job":
			var jobs []cli.Job
			result, err = lm.fetcher.FetchData(ctx, params, &jobs)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, jobs, result, limit, offset)
			}

		case "runtemplate":
			var templates []cli.RunTemplate
			result, err = lm.fetcher.FetchData(ctx, params, &templates)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, templates, result, limit, offset)
			}

		case "artifact":
			var artifacts []cli.Artifact
			result, err = lm.fetcher.FetchData(ctx, params, &artifacts)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, artifacts, result, limit, offset)
			}

		case "companyapp":
			var companyapps []cli.CompanyApp
			result, err = lm.fetcher.FetchData(ctx, params, &companyapps)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, companyapps, result, limit, offset)
			}

		case "queue":
			var queue []cli.Queue
			result, err = lm.fetcher.FetchData(ctx, params, &queue)
			if err == nil {
				return lm.createDataLoadedMsg(config.EntityType, queue, result, limit, offset)
			}

		default:
			return DataErrorMsg{
				EntityType: config.EntityType,
				Error:      fmt.Errorf("unsupported entity type: %s", config.EntityType),
			}
		}

		return DataErrorMsg{
			EntityType: config.EntityType,
			Error:      err,
		}
	}
}

// LoadBatchData loads multiple entity types concurrently
func (lm *ListingManager) LoadBatchData(requests []BatchLoadRequest) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		// Convert to fetcher batch requests
		batchReqs := make([]cli.BatchRequest, len(requests))
		for i, req := range requests {
			params := cli.FetchParams{
				Command:    req.EntityType,
				SubCommand: "list",
				Limit:      req.Limit,
				Offset:     req.Offset,
				CacheTTL:   req.CacheTTL,
				Timeout:    req.Timeout,
			}

			// Create target based on entity type
			var target interface{}
			switch req.EntityType {
			case "application":
				target = &[]cli.Application{}
			case "company":
				target = &[]cli.Company{}
			case "crprototype":
				target = &[]cli.CrPrototype{}
			default:
				continue // Skip unsupported types
			}

			batchReqs[i] = cli.BatchRequest{
				Key:    req.EntityType,
				Params: params,
				Target: target,
			}
		}

		// Execute batch fetch
		results := lm.fetcher.FetchBatch(ctx, batchReqs)

		// Convert results to messages
		var messages []tea.Msg
		for _, result := range results {
			if result.Error != nil {
				messages = append(messages, DataErrorMsg{
					EntityType: result.Key,
					Error:      result.Error,
				})
			} else {
				// Create appropriate data loaded message
				msg := lm.createDataLoadedMsg(result.Key, result.Result.Data, result.Result, 10, 0)
				messages = append(messages, msg)
			}
		}

		return BatchDataLoadedMsg{Messages: messages}
	}
}

// BatchLoadRequest represents a request for batch loading
type BatchLoadRequest struct {
	EntityType string
	Limit      int
	Offset     int
	CacheTTL   time.Duration
	Timeout    time.Duration
}

// BatchDataLoadedMsg represents multiple data loads completing
type BatchDataLoadedMsg struct {
	Messages []tea.Msg
}

// RefreshData refreshes data with optional force flag
func (lm *ListingManager) RefreshData(config ListingConfig, limit, offset int, force bool) tea.Cmd {
	return func() tea.Msg {
		if force {
			// Clear relevant cache entries
			lm.fetcher.ClearExpiredCache()
		}

		// Load data normally (cache will be bypassed if cleared)
		cmd := lm.LoadData(config, limit, offset)
		return cmd()
	}
}

// SetupAutoRefresh sets up automatic data refresh
func (lm *ListingManager) SetupAutoRefresh(config ListingConfig, limit, offset int) tea.Cmd {
	if config.RefreshInterval == 0 {
		return nil // Auto-refresh disabled
	}

	return tea.Tick(config.RefreshInterval, func(t time.Time) tea.Msg {
		return RefreshDataMsg{
			EntityType: config.EntityType,
			Force:      false,
		}
	})
}

// GetCacheStats returns cache statistics
func (lm *ListingManager) GetCacheStats() cli.CacheStats {
	return lm.fetcher.GetCacheStats()
}

// ClearCache clears all cached data
func (lm *ListingManager) ClearCache() {
	lm.fetcher.ClearCache()
}

// Helper methods

func (lm *ListingManager) createDataLoadedMsg(entityType string, data interface{}, result *cli.FetchResult, limit, offset int) DataLoadedMsg {
	state := ListingState{
		Loading:     false,
		Data:        data,
		CurrentPage: (offset / limit) + 1,
		Limit:       limit,
		Offset:      offset,
		TotalCount:  result.TotalCount,
		HasMore:     result.TotalCount == limit, // Assume more if we got full page
		HasPrev:     offset > 0,
		LastUpdate:  time.Now(),
		Cached:      result.Cached,
		FetchTime:   result.FetchTime,
	}

	return DataLoadedMsg{
		EntityType: entityType,
		Data:       data,
		State:      state,
	}
}

// DefaultConfigs provides default configurations for common entity types
var DefaultConfigs = map[string]ListingConfig{
	"application": {
		EntityType:      "application",
		RefreshInterval: 5 * time.Minute,
		CacheTTL:        2 * time.Minute,
		DefaultLimit:    10,
		MaxLimit:        100,
		Timeout:         10 * time.Second,
	},
	"company": {
		EntityType:      "company",
		RefreshInterval: 10 * time.Minute,
		CacheTTL:        5 * time.Minute,
		DefaultLimit:    10,
		MaxLimit:        50,
		Timeout:         10 * time.Second,
	},
	"crprototype": {
		EntityType:      "crprototype",
		RefreshInterval: 2 * time.Minute,
		CacheTTL:        1 * time.Minute,
		DefaultLimit:    10,
		MaxLimit:        50,
		Timeout:         10 * time.Second,
	},
	"job": {
		EntityType:      "job",
		RefreshInterval: 30 * time.Second,
		CacheTTL:        15 * time.Second,
		DefaultLimit:    10,
		MaxLimit:        100,
		Timeout:         15 * time.Second,
	},
}

// Global instance
var DefaultListingManager = NewListingManager()
