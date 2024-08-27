package cache

import (
	"context"
	"time"
)

type CacheElement[T any] struct {
	Key       string
	Value     T
	CreatedAt time.Time
	TTL       time.Duration
}

type Cache[T any] interface {
	Get(key string) (CacheElement[T], bool)
	Put(key string, value T)
	GetAllKeyValues() map[string]T
	Sweep(ctx context.Context, checkInterval time.Duration)
}
