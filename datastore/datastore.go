package datastore

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Datastore is the implementation of the datastore
type Datastore[K MapKey, T any] struct {
	elements map[K]CacheItem[T]
	ttl      time.Duration
	mutex    sync.RWMutex
	opts     Options
}

// NewDatastore creates a new instance of the datastore with the given TTL
func NewDatastore[K MapKey, T any](
	ctx context.Context,
	ttl time.Duration,
	opts ...Option,
) (*Datastore[K, T], error) {
	memCache := Datastore[K, T]{
		elements: make(map[K]CacheItem[T]),
		ttl:      ttl,
	}
	err := memCache.opts.Apply(opts...)
	if err != nil {
		return nil, fmt.Errorf("options apply: %w", err)
	}

	if memCache.opts.SweepInterval > 0 {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic in sweep goroutine: %v\n", r)
				}
			}()
			memCache.sweep(ctx, memCache.opts.SweepInterval)
		}()
	}

	return &memCache, nil
}

// Get retrieves a value from the datastore, returning whether it exists and is valid
func (ele *Datastore[K, T]) Get(key K) (CacheItem[T], bool) {
	ele.mutex.RLock()
	val, ok := ele.elements[key]
	ele.mutex.RUnlock()

	if ok && ele.isExpired(val) {
		// Upgrade to write lock and remove the expired key
		ele.mutex.Lock()
		delete(ele.elements, key)
		ele.mutex.Unlock()
		return CacheItem[T]{}, false
	}
	return val, ok
}

// Put adds a new record to the datastore
func (ele *Datastore[K, T]) Put(key K, value T) {
	ele.mutex.Lock()
	defer ele.mutex.Unlock()

	ele.elements[key] = NewCacheItem(value, ele.ttl)
}

// GetAllKeyValues returns all valid key-value pairs, removing expired records
func (ele *Datastore[K, T]) GetAllKeyValues() map[K]T {
	ele.mutex.Lock()
	defer ele.mutex.Unlock()

	allRecord := make(map[K]T)
	for key, val := range ele.elements {
		if ele.isExpired(val) {
			delete(ele.elements, key)
			continue
		}
		allRecord[key] = val.Value
	}
	return allRecord
}

// isExpired checks if a record has exceeded its TTL
func (ele *Datastore[K, T]) isExpired(record CacheItem[T]) bool {
	return time.Since(record.CreatedAt) > ele.ttl
}

// sweep runs as a goroutine to automatically remove expired records
func (ele *Datastore[K, T]) sweep(ctx context.Context, checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ele.GetAllKeyValues() // Will clean up expired elements
		case <-ctx.Done():
			log.Println("sweep closed")
			return
		}
	}
}
