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
	var ttlDuration = []time.Duration{time.Duration(1) * time.Hour}
	if len(ttl) > 0 {
		ttlDuration[0] = ttl[0]
	}
	return CacheElement[T]{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().UTC(),
		TTL:       ttlDuration[0],
	}
}
