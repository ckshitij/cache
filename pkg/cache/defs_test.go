package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testKey   = "testKey"
	testValue = "testValue"
)

func TestNewCacheElement_WithDefaultTTL(t *testing.T) {
	// Test cache element creation without providing TTL (default TTL case)
	cacheElement := NewCacheElement(testKey, testValue)

	// Assertions
	assert.Equal(t, testKey, cacheElement.Key, "Key should be set correctly")
	assert.Equal(t, testValue, cacheElement.Value, "Value should be set correctly")
	assert.WithinDuration(t, time.Now().UTC(), cacheElement.CreatedAt, time.Second, "CreatedAt should be set to the current time")
	assert.Equal(t, time.Hour, cacheElement.TTL, "Default TTL should be 24 hours")
}

func TestNewCacheElement_WithCustomTTL(t *testing.T) {
	// Test cache element creation with a custom TTL
	customTTL := 2 * time.Hour

	cacheElement := NewCacheElement(testKey, testValue, customTTL)

	// Assertions
	assert.Equal(t, testKey, cacheElement.Key, "Key should be set correctly")
	assert.Equal(t, testValue, cacheElement.Value, "Value should be set correctly")
	assert.WithinDuration(t, time.Now().UTC(), cacheElement.CreatedAt, time.Second, "CreatedAt should be set to the current time")
	assert.Equal(t, customTTL, cacheElement.TTL, "Custom TTL should be set correctly")
}

func TestNewCacheElement_InvalidTTLHandling(t *testing.T) {
	// Test edge case where TTL slice is passed but is empty
	cacheElement := NewCacheElement(testKey, testValue, []time.Duration{}...)

	// Assertions
	assert.Equal(t, time.Hour, cacheElement.TTL, "Default TTL should be an hours when empty TTL slice is passed")
}
