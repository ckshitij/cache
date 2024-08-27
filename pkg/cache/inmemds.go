package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// inMemoryCache is the implementation of the datastore
type InMemoryCache[T any] struct {
	elements sync.Map
	ttl      time.Duration
	opts     Options
}

// NewKeyValueCache creates a new instance of the datastore with the given TTL
func NewKeyValueCache[T any](ctx context.Context, ttl time.Duration, opts ...Option) (*InMemoryCache[T], error) {
	memCache := InMemoryCache[T]{
		ttl: ttl,
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
func (c *InMemoryCache[T]) Get(key string) (T, bool) {
	value, ok := c.elements.Load(key)
	if !ok {
		var zeroValue T
		return zeroValue, false
	}

	element := value.(CacheElement[T])
	if element.CreatedAt.Add(element.TTL).Before(time.Now()) {
		var zeroValue T
		c.elements.Delete(key)
		return zeroValue, false
	}

	return element.Value, true
}

// Put adds a new record to the datastore
func (c *InMemoryCache[T]) Put(key string, value T) {
	element := CacheElement[T]{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		TTL:       c.ttl,
	}
	c.elements.Store(key, element)
}

// GetAllKeyValues returns all valid key-value pairs, removing expired records
func (c *InMemoryCache[T]) GetAllKeyValues() map[string]T {
	results := make(map[string]T)
	c.elements.Range(func(key, value interface{}) bool {
		element := value.(CacheElement[T])
		if !c.isExpired(element) {
			results[key.(string)] = element.Value
		} else {
			c.elements.Delete(key)
		}
		return true
	})
	return results
}

// evict removed the expired keys
func (c *InMemoryCache[T]) evict() {
	c.elements.Range(func(key, value interface{}) bool {
		element := value.(CacheElement[T])
		if c.isExpired(element) {
			c.elements.Delete(key)
		}
		return true
	})
}

// isExpired checks if a record has exceeded its TTL and deletes it if so
func (c *InMemoryCache[T]) isExpired(record CacheElement[T]) bool {
	return time.Since(record.CreatedAt) > c.ttl
}

func (c *InMemoryCache[T]) sweep(ctx context.Context, sweep time.Duration) {
	ticker := time.NewTicker(sweep)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.evict()
		case <-ctx.Done():
			fmt.Println("sweeping closed")
			return
		}
	}
}
