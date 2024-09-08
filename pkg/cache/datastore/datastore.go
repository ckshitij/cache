package datastore

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ckshitij/cache/pkg/cache"
)

// Datastore is the implementation of the datastore
type Datastore[T any] struct {
	elements map[string]cache.CacheElement[T]
	ttl      time.Duration
	mutex    sync.RWMutex
	opts     cache.Options
}

// NewDatastore creates a new instance of the datastore with the given TTL
func NewDatastore[T any](
	ctx context.Context,
	ttl time.Duration,
	opts ...cache.Option,
) (*Datastore[T], error) {
	memCache := Datastore[T]{
		elements: make(map[string]cache.CacheElement[T]),
		ttl:      ttl,
	}
	err := memCache.opts.Apply(opts...)
	if err != nil {
		return nil, fmt.Errorf("options apply: %w", err)
	}

	if memCache.opts.SweepInterval > 0 {
		go memCache.sweep(ctx, memCache.opts.SweepInterval)
	}

	return &memCache, nil
}

// Get retrieves a value from the datastore, returning whether it exists and is valid
func (ele *Datastore[T]) Get(key string) (cache.CacheElement[T], bool) {
	ele.mutex.RLock()
	defer ele.mutex.RUnlock()

	val, ok := ele.elements[key]
	if ok && ele.isExpired(val) {
		return cache.CacheElement[T]{}, false
	}
	return val, ok
}

// Put adds a new record to the datastore
func (ele *Datastore[T]) Put(key string, value T) {
	ele.mutex.Lock()
	defer ele.mutex.Unlock()

	ele.elements[key] = cache.NewCacheElement(key, value, ele.ttl)
}

// GetAllKeyValues returns all valid key-value pairs, removing expired records
func (ele *Datastore[T]) GetAllKeyValues() map[string]T {
	ele.mutex.Lock()
	defer ele.mutex.Unlock()

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
func (ele *Datastore[T]) isExpired(record cache.CacheElement[T]) bool {
	return time.Since(record.CreatedAt) > ele.ttl
}

// sweep runs as a goroutine to automatically remove expired records
func (ele *Datastore[T]) sweep(ctx context.Context, checkInterval time.Duration) {
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
