package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// inMemoryCache is the implementation of the datastore
type inMemoryCache[T any] struct {
	elements map[string]CacheElement[T]
	ttl      time.Duration
	sync.RWMutex
	opts Options
}

// NewKeyValueCache creates a new instance of the datastore with the given TTL
func NewKeyValueCache[T any](
	ctx context.Context,
	ttl time.Duration,
	opts ...Option,
) (*inMemoryCache[T], error) {
	memCache := inMemoryCache[T]{
		elements: make(map[string]CacheElement[T]),
		ttl:      ttl,
	}
	err := memCache.opts.Apply(opts...)
	if err != nil {
		return nil, fmt.Errorf("options apply: %w", err)
	}

	if memCache.opts.sweepInterval > 0 {
		go memCache.sweep(ctx, memCache.opts.sweepInterval)
	}

	return &memCache, nil
}

// Get retrieves a value from the datastore, returning whether it exists and is valid
func (ele *inMemoryCache[T]) Get(key string) (CacheElement[T], bool) {
	ele.RLock()
	defer ele.RUnlock()

	val, ok := ele.elements[key]
	if ok && ele.isExpired(val) {
		return CacheElement[T]{}, false
	}
	return val, ok
}

// Put adds a new record to the datastore
func (ele *inMemoryCache[T]) Put(key string, value T) {
	ele.Lock()
	defer ele.Unlock()

	ele.elements[key] = CacheElement[T]{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().UTC(),
		TTL:       ele.ttl,
	}
}

// GetAllKeyValues returns all valid key-value pairs, removing expired records
func (ele *inMemoryCache[T]) GetAllKeyValues() map[string]T {
	ele.Lock()
	defer ele.Unlock()

	allRecord := make(map[string]T)
	for key, val := range ele.elements {
		if ele.isExpired(val) {
			delete(ele.elements, key)
			continue
		}
		allRecord[key] = val.Value
	}
	return allRecord
}

// isExpired checks if a record has exceeded its TTL and deletes it if so
func (ele *inMemoryCache[T]) isExpired(record CacheElement[T]) bool {
	return time.Since(record.CreatedAt) > ele.ttl
}

// sweep runs as a goroutine to automatically remove expired records
func (ele *inMemoryCache[T]) sweep(ctx context.Context, checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ele.GetAllKeyValues()
		case <-ctx.Done():
			fmt.Println("sweep closed")
			return
		}
	}
}
