package cache

import (
	"time"
)

type CacheElement[T any] struct {
	Key       string
	Value     T
	CreatedAt time.Time
	TTL       time.Duration
}

/*
Create new cache element, if ttl is not given
then default value will be 1 day
*/
func NewCacheElement[T any](key string, value T, ttl ...time.Duration) CacheElement[T] {
	if len(ttl) < 1 {
		ttl[0] = time.Hour * 24
	}
	return CacheElement[T]{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().UTC(),
		TTL:       ttl[0],
	}
}
