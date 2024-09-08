package autoreload

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ckshitij/cache/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock DataFunc that returns some mock data for testing
func mockLoadFuncSuccess() (map[string]cache.CacheElement[string], error) {
	data := map[string]cache.CacheElement[string]{
		"key1": {Value: "value1"},
		"key2": {Value: "value2"},
	}
	return data, nil
}

func mockLoadFuncFail() (map[string]cache.CacheElement[string], error) {
	return nil, errors.New("failed to load data")
}

func TestNewAutoReload_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test cache creation with a successful load function
	autoCache, err := NewAutoReload[string](ctx, "testCache", mockLoadFuncSuccess)
	require.NoError(t, err, "Error should be nil on creating cache")
	assert.NotNil(t, autoCache, "Cache should not be nil after initialization")

	// Test that data is loaded into cache immediately
	data, ok := autoCache.Get("key1")
	assert.True(t, ok, "Cache key 'key1' should exist")
	assert.Equal(t, "value1", data.Value, "Cache value for 'key1' should be 'value1'")

	// Test Get for non-existent key
	_, ok = autoCache.Get("invalidKey")
	assert.False(t, ok, "Cache should return false for non-existent keys")

	// Test refresh duration
	assert.Equal(t, time.Second*10, autoCache.GetRefreshDuration(), "Refresh duration should be 1 minute")

	// Test that the initial data is available in cache
	data, ok = autoCache.Get("key2")
	assert.True(t, ok, "Cache key 'key2' should exist")
	assert.Equal(t, "value2", data.Value, "Cache value for 'key2' should be 'value2'")
}

func TestNewAutoReload_Failure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test cache creation with a failing load function
	_, err := NewAutoReload[string](ctx, "failingCache", mockLoadFuncFail)
	require.Error(t, err, "Expected an error on creating cache with a failing load function")
	assert.Contains(t, err.Error(), "cache initialization failed", "Error message should indicate cache initialization failure")
}

func TestAutoReload_ReloadSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test cache creation with a successful load function
	autoCache, err := NewAutoReload[string](ctx, "testCache", mockLoadFuncSuccess)
	require.NoError(t, err, "Error should be nil on creating cache")

	// Simulate time passing to trigger the auto-reload
	time.Sleep(11 * time.Second)

	// After auto-reload, data should still be valid
	data, ok := autoCache.Get("key1")
	assert.True(t, ok, "Cache key 'key1' should exist after reload")
	assert.Equal(t, "value1", data.Value, "Cache value for 'key1' should be 'value1'")
}
