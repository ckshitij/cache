package datastore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testKey   = "testKey"
	testValue = "testValue"
)

func TestNewCacheItem_WithDefaultTTL(t *testing.T) {
	// Test cache element creation without providing TTL (default TTL case)
	cacheElement := NewCacheItem(testValue)

	// Assertions
	assert.Equal(t, testValue, cacheElement.Value, "Value should be set correctly")
	assert.WithinDuration(t, time.Now().UTC(), cacheElement.CreatedAt, time.Second, "CreatedAt should be set to the current time")
	assert.Equal(t, 24*time.Hour, cacheElement.TTL, "Default TTL should be 24 hours")
}

func TestNewCacheItem_WithCustomTTL(t *testing.T) {
	// Test cache element creation with a custom TTL
	customTTL := 2 * time.Hour

	cacheElement := NewCacheItem(testValue, customTTL)

	// Assertions
	assert.Equal(t, testValue, cacheElement.Value, "Value should be set correctly")
	assert.WithinDuration(t, time.Now().UTC(), cacheElement.CreatedAt, time.Second, "CreatedAt should be set to the current time")
	assert.Equal(t, customTTL, cacheElement.TTL, "Custom TTL should be set correctly")
}

func TestNewCacheItem_InvalidTTLHandling(t *testing.T) {
	// Test edge case where TTL slice is passed but is empty
	cacheElement := NewCacheItem(testValue, []time.Duration{time.Hour}...)

	// Assertions
	assert.Equal(t, time.Hour, cacheElement.TTL, "Default TTL should be an hours when empty TTL slice is passed")
}
