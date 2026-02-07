package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"sync"
	"time"
)

// DataFetcher provides a unified interface for fetching JSON data from multiflexi-cli
type DataFetcher struct {
	cache      map[string]cachedResult
	cacheMutex sync.RWMutex
	timeout    time.Duration
}

// cachedResult represents a cached CLI result
type cachedResult struct {
	data      interface{}
	timestamp time.Time
	ttl       time.Duration
}

// FetchParams represents parameters for data fetching
type FetchParams struct {
	Command    string        // CLI command (e.g., "application")
	SubCommand string        // CLI subcommand (e.g., "list")
	Format     string        // Output format (default: "json")
	Order      string        // Sort order (default: "D")
	Limit      int           // Results limit (default: 10)
	Offset     int           // Results offset (default: 0)
	CacheTTL   time.Duration // Cache time-to-live (default: 30s)
	Timeout    time.Duration // Request timeout (default: 10s)
}

// FetchResult represents the result of a data fetch operation
type FetchResult struct {
	Data       interface{}
	Error      error
	Cached     bool
	FetchTime  time.Duration
	TotalCount int // If available from CLI
}

// NewDataFetcher creates a new data fetcher with default configuration
func NewDataFetcher() *DataFetcher {
	return &DataFetcher{
		cache:   make(map[string]cachedResult),
		timeout: 10 * time.Second,
	}
}

// FetchData fetches data from multiflexi-cli with caching and error handling
func (df *DataFetcher) FetchData(ctx context.Context, params FetchParams, target interface{}) (*FetchResult, error) {
	start := time.Now()

	// Set defaults
	if params.Format == "" {
		params.Format = "json"
	}
	if params.Order == "" {
		params.Order = "D"
	}
	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.CacheTTL == 0 {
		params.CacheTTL = 30 * time.Second
	}
	if params.Timeout == 0 {
		params.Timeout = df.timeout
	}

	cacheKey := df.buildCacheKey(params)

	// Check cache first
	if cached, found := df.getCached(cacheKey); found {
		if err := df.copyData(cached.data, target); err != nil {
			return nil, fmt.Errorf("failed to copy cached data: %w", err)
		}
		return &FetchResult{
			Data:      target,
			Cached:    true,
			FetchTime: time.Since(start),
		}, nil
	}

	// Create context with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, params.Timeout)
	defer cancel()

	// Build command arguments
	args := df.buildArgs(params)

	// Execute CLI command
	cmd := exec.CommandContext(cmdCtx, "multiflexi-cli", args...)
	output, err := cmd.Output()
	if err != nil {
		return &FetchResult{
			Error:     fmt.Errorf("CLI command failed: %w", err),
			FetchTime: time.Since(start),
		}, fmt.Errorf("failed to execute %s %s: %w", params.Command, params.SubCommand, err)
	}

	// Parse JSON response
	if err := json.Unmarshal(output, target); err != nil {
		return &FetchResult{
			Error:     fmt.Errorf("JSON parsing failed: %w", err),
			FetchTime: time.Since(start),
		}, fmt.Errorf("failed to parse JSON output for %s %s: %w", params.Command, params.SubCommand, err)
	}

	// Cache the result
	df.setCached(cacheKey, target, params.CacheTTL)

	// Calculate result count if it's a slice
	totalCount := 0
	if reflect.TypeOf(target).Kind() == reflect.Ptr {
		if slice := reflect.ValueOf(target).Elem(); slice.Kind() == reflect.Slice {
			totalCount = slice.Len()
		}
	}

	return &FetchResult{
		Data:       target,
		Cached:     false,
		FetchTime:  time.Since(start),
		TotalCount: totalCount,
	}, nil
}

// FetchBatch fetches multiple datasets concurrently
func (df *DataFetcher) FetchBatch(ctx context.Context, requests []BatchRequest) []BatchResult {
	results := make([]BatchResult, len(requests))
	var wg sync.WaitGroup

	for i, req := range requests {
		wg.Add(1)
		go func(index int, request BatchRequest) {
			defer wg.Done()

			result, err := df.FetchData(ctx, request.Params, request.Target)
			results[index] = BatchResult{
				Key:    request.Key,
				Result: result,
				Error:  err,
			}
		}(i, req)
	}

	wg.Wait()
	return results
}

// BatchRequest represents a single request in a batch
type BatchRequest struct {
	Key    string      // Identifier for this request
	Params FetchParams // Fetch parameters
	Target interface{} // Target to unmarshal into
}

// BatchResult represents the result of a batch request
type BatchResult struct {
	Key    string
	Result *FetchResult
	Error  error
}

// ClearCache removes all cached entries
func (df *DataFetcher) ClearCache() {
	df.cacheMutex.Lock()
	defer df.cacheMutex.Unlock()
	df.cache = make(map[string]cachedResult)
}

// ClearExpiredCache removes expired cache entries
func (df *DataFetcher) ClearExpiredCache() {
	df.cacheMutex.Lock()
	defer df.cacheMutex.Unlock()

	now := time.Now()
	for key, cached := range df.cache {
		if now.Sub(cached.timestamp) > cached.ttl {
			delete(df.cache, key)
		}
	}
}

// GetCacheStats returns cache statistics
func (df *DataFetcher) GetCacheStats() CacheStats {
	df.cacheMutex.RLock()
	defer df.cacheMutex.RUnlock()

	total := len(df.cache)
	expired := 0
	now := time.Now()

	for _, cached := range df.cache {
		if now.Sub(cached.timestamp) > cached.ttl {
			expired++
		}
	}

	return CacheStats{
		TotalEntries:   total,
		ExpiredEntries: expired,
		ActiveEntries:  total - expired,
	}
}

// CacheStats represents cache statistics
type CacheStats struct {
	TotalEntries   int
	ExpiredEntries int
	ActiveEntries  int
}

// Helper methods

func (df *DataFetcher) buildCacheKey(params FetchParams) string {
	return fmt.Sprintf("%s:%s:fmt=%s:ord=%s:lim=%d:off=%d",
		params.Command, params.SubCommand, params.Format, params.Order, params.Limit, params.Offset)
}

func (df *DataFetcher) buildArgs(params FetchParams) []string {
	args := []string{params.Command, params.SubCommand}

	if params.Format != "" {
		args = append(args, "--format="+params.Format)
	}
	if params.Order != "" {
		args = append(args, "--order="+params.Order)
	}
	if params.Limit > 0 {
		args = append(args, fmt.Sprintf("--limit=%d", params.Limit))
	}
	if params.Offset > 0 {
		args = append(args, fmt.Sprintf("--offset=%d", params.Offset))
	}

	return args
}

func (df *DataFetcher) getCached(key string) (cachedResult, bool) {
	df.cacheMutex.RLock()
	defer df.cacheMutex.RUnlock()

	cached, exists := df.cache[key]
	if !exists {
		return cachedResult{}, false
	}

	// Check if expired
	if time.Since(cached.timestamp) > cached.ttl {
		return cachedResult{}, false
	}

	return cached, true
}

func (df *DataFetcher) setCached(key string, data interface{}, ttl time.Duration) {
	df.cacheMutex.Lock()
	defer df.cacheMutex.Unlock()

	// Deep copy the data to avoid reference issues
	copiedData := df.deepCopy(data)

	df.cache[key] = cachedResult{
		data:      copiedData,
		timestamp: time.Now(),
		ttl:       ttl,
	}
}

func (df *DataFetcher) deepCopy(src interface{}) interface{} {
	// Simple deep copy using JSON marshal/unmarshal
	// This works well for our JSON-serializable structs
	data, _ := json.Marshal(src)

	srcType := reflect.TypeOf(src)
	dst := reflect.New(srcType.Elem()).Interface()
	json.Unmarshal(data, dst)

	return dst
}

func (df *DataFetcher) copyData(src, dst interface{}) error {
	// Copy data from src to dst using JSON marshal/unmarshal
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dst)
}

// Convenience methods for common operations

// FetchApplications fetches applications using the improved data fetcher
func (df *DataFetcher) FetchApplications(ctx context.Context, limit, offset int) ([]Application, error) {
	var apps []Application

	params := FetchParams{
		Command:    "application",
		SubCommand: "list",
		Limit:      limit,
		Offset:     offset,
	}

	result, err := df.FetchData(ctx, params, &apps)
	if err != nil {
		return nil, err
	}

	_ = result // Can be used for additional metadata
	return apps, nil
}

// FetchCompanies fetches companies using the improved data fetcher
func (df *DataFetcher) FetchCompanies(ctx context.Context, limit, offset int) ([]Company, error) {
	var companies []Company

	params := FetchParams{
		Command:    "company",
		SubCommand: "list",
		Limit:      limit,
		Offset:     offset,
	}

	result, err := df.FetchData(ctx, params, &companies)
	if err != nil {
		return nil, err
	}

	_ = result
	return companies, nil
}

// FetchCrPrototypes fetches credential prototypes using the improved data fetcher
func (df *DataFetcher) FetchCrPrototypes(ctx context.Context, limit, offset int) ([]CrPrototype, error) {
	var crprototypes []CrPrototype

	params := FetchParams{
		Command:    "crprototype",
		SubCommand: "list",
		Limit:      limit,
		Offset:     offset,
	}

	result, err := df.FetchData(ctx, params, &crprototypes)
	if err != nil {
		return nil, err
	}

	_ = result
	return crprototypes, nil
}

// Global instance for backward compatibility
var DefaultFetcher = NewDataFetcher()

// Enhanced versions of existing functions using the new data fetcher
func GetApplicationsEnhanced(limit, offset int) ([]Application, error) {
	return DefaultFetcher.FetchApplications(context.Background(), limit, offset)
}

func GetCompaniesEnhanced(limit, offset int) ([]Company, error) {
	return DefaultFetcher.FetchCompanies(context.Background(), limit, offset)
}

func GetCrPrototypesEnhanced(limit, offset int) ([]CrPrototype, error) {
	return DefaultFetcher.FetchCrPrototypes(context.Background(), limit, offset)
}
