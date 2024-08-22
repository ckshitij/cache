package inmemds

import "time"

type DatastoreEntity struct {
	Key       string
	Value     any
	CreatedAt time.Time
	TTL       time.Duration
}

type KeyValueDataStore interface {
	Get(key string) (DatastoreEntity, bool)
	Put(key string, value any)
	GetAllKeyValues() map[string]any
	AutoCleanUp(checkInterval time.Duration, done <-chan bool)
}
