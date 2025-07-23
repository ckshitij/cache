package datastore

import (
	"time"
)

// MapKey is a type constraint for keys that can be used in maps (i.e., comparable types).
type MapKey interface {
	comparable
}

// CacheItem wraps a cached value with metadata like TTL and creation time.
type CacheItem[T any] struct {
	Value     T
	CreatedAt time.Time
	TTL       time.Duration
}

// DefaultTTL defines the fallback TTL if none is provided (24 hours).
const DefaultTTL = 24 * time.Hour

// NewCacheItem creates a new CacheItem with an optional TTL.
// If TTL is not provided, DefaultTTL is used.
func NewCacheItem[T any](value T, ttl ...time.Duration) CacheItem[T] {
	actualTTL := DefaultTTL
	if len(ttl) > 0 {
		actualTTL = ttl[0]
	}
	return CacheItem[T]{
		Value:     value,
		CreatedAt: time.Now().UTC(),
		TTL:       actualTTL,
	}
}
