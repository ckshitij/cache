package datastore

import "time"

type DataStoreElement struct {
	Key       string
	Value     any
	CreatedAt time.Time
	TTL       time.Duration
}

type DataStore interface {
	Get(key string) (any, bool)
	Put(key string, value any)
	GetAllKeyValues() map[string]any
	AutoCleanUp(checkInterval time.Duration, done <-chan bool)
}
