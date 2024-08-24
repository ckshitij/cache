package cache

import "time"

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
	AutoCleanUp(checkInterval time.Duration, done <-chan bool)
}
