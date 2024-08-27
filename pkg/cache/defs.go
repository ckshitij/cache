package cache

import "time"

type CacheElement[T any] struct {
	Key       string
	Value     T
	CreatedAt time.Time
	TTL       time.Duration
}
