package datastore

import "time"

type inMemoryDataStore struct {
	elements map[string]DataStoreElement
	ttl      time.Duration
}

func NewDataStore(ttl time.Duration) DataStore {
	return &inMemoryDataStore{
		elements: make(map[string]DataStoreElement),
		ttl:      ttl,
	}
}

func (ele *inMemoryDataStore) Get(key string) (any, bool) {
	val, ok := ele.elements[key]
	return val, ok
}

func (ele *inMemoryDataStore) Put(key string, value string) {
	ele.elements[key] = DataStoreElement{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().UTC(),
		TTL:       ele.ttl,
	}
}

func (ele *inMemoryDataStore) GetAllKeys() ([]string, error) {
	return nil, nil
}
