package ui

import (
	"context"
	"testing"
	"time"

	"github.com/VitexSoftware/multiflexi-tui/internal/cli"
)

func TestListingManagerInit(t *testing.T) {
	manager := NewListingManager()
	if manager == nil {
		t.Error("NewListingManager should not return nil")
	}
	if manager.fetcher == nil {
		t.Error("ListingManager should have a fetcher")
	}
}

func TestListingManagerCacheOps(t *testing.T) {
	manager := NewListingManager()

	// Test cache stats
	stats := manager.GetCacheStats()
	if stats.TotalEntries != 0 {
		t.Errorf("Expected 0 cache entries initially, got %d", stats.TotalEntries)
	}

	// Test cache clear
	manager.ClearCache()
	stats = manager.GetCacheStats()
	if stats.TotalEntries != 0 {
		t.Errorf("Expected 0 cache entries after clear, got %d", stats.TotalEntries)
	}
}

func TestDefaultConfigs(t *testing.T) {
	// Test that default configs exist for key entity types
	requiredConfigs := []string{"application", "company", "crprototype", "job"}

	for _, entityType := range requiredConfigs {
		config, exists := DefaultConfigs[entityType]
		if !exists {
			t.Errorf("Missing default config for entity type: %s", entityType)
			continue
		}

		if config.EntityType != entityType {
			t.Errorf("Config entity type mismatch: expected %s, got %s", entityType, config.EntityType)
		}

		if config.DefaultLimit <= 0 {
			t.Errorf("Config should have positive default limit for %s", entityType)
		}

		if config.Timeout <= 0 {
			t.Errorf("Config should have positive timeout for %s", entityType)
		}
	}
}

func TestEnhancedCrPrototypesModel(t *testing.T) {
	model := NewCrPrototypesModel()

	// Test that enhanced model is properly initialized
	if model.listingManager == nil {
		t.Error("CrPrototypesModel should have a listing manager")
	}

	if model.config.EntityType != "crprototype" {
		t.Errorf("Expected entity type 'crprototype', got '%s'", model.config.EntityType)
	}

	if model.limit != model.config.DefaultLimit {
		t.Errorf("Expected limit %d, got %d", model.config.DefaultLimit, model.limit)
	}
}

func TestDataFetcherIntegration(t *testing.T) {
	// This test would require a running multiflexi-cli instance
	// Skip if not available
	t.Skip("Integration test - requires multiflexi-cli")

	fetcher := cli.NewDataFetcher()
	ctx := context.Background()

	// Test cache functionality
	params := cli.FetchParams{
		Command:    "application",
		SubCommand: "list",
		Limit:      1,
		Offset:     0,
		CacheTTL:   1 * time.Second,
	}

	var apps []cli.Application
	result1, err := fetcher.FetchData(ctx, params, &apps)
	if err != nil {
		t.Fatalf("First fetch failed: %v", err)
	}

	if result1.Cached {
		t.Error("First fetch should not be cached")
	}

	// Fetch again immediately - should be cached
	var apps2 []cli.Application
	result2, err := fetcher.FetchData(ctx, params, &apps2)
	if err != nil {
		t.Fatalf("Second fetch failed: %v", err)
	}

	if !result2.Cached {
		t.Error("Second fetch should be cached")
	}

	if result2.FetchTime >= result1.FetchTime {
		t.Error("Cached fetch should be faster than first fetch")
	}
}
