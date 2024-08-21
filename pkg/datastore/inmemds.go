package datastore

import (
	"fmt"
	"sync"
	"time"
)

type inMemoryKeyValueDatastore struct {
	// Map maintains the value
	elements map[string]DataStoreElement
	// element will not be
	// valid after the ttl duration
	ttl time.Duration
	// for handling concurrency using mutex
	sync.RWMutex
}

/*
Create a key value datastore instance with given ttl
(time to live).
*/
func NewKeyValueDataStore(ttl time.Duration) DataStore {
	return &inMemoryKeyValueDatastore{
		elements: make(map[string]DataStoreElement),
		ttl:      ttl,
	}
}

/*
Return the value and a boolean value which indicate
Whether the element is present on or not in the
datastore.
If the element exceed the ttl then it will get deleted
and return nil and false value.
*/
func (ele *inMemoryKeyValueDatastore) Get(key string) (any, bool) {
	ele.RLock()
	val, ok := ele.elements[key]
	if ok && ele.isExpired(val) {
		return nil, false
	}
	return val, ok
}

/*
Add new record into the datastore.
*/
func (ele *inMemoryKeyValueDatastore) Put(key string, value any) {
	ele.Lock()
	defer ele.Unlock()
	ele.elements[key] = DataStoreElement{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().UTC(),
		TTL:       ele.ttl,
	}
}

/*
Return all the valid key-values to the user and also remove the
expired record from the datastore
*/
func (ele *inMemoryKeyValueDatastore) GetAllKeyValues() map[string]any {
	var allRecord = make(map[string]any)
	for key, val := range ele.elements {
		if ele.isExpired(val) {
			continue
		}
		allRecord[key] = val.Value
	}
	return allRecord
}

func (ele *inMemoryKeyValueDatastore) isExpired(record DataStoreElement) bool {
	if time.Since(record.CreatedAt) > ele.ttl {
		ele.Lock()
		fmt.Println("Removing data with key: ", record.Key, " value : ", record.Value)
		delete(ele.elements, record.Key)
		ele.Unlock()
		return true
	}
	return false
}

/*
Run it as go-routine so that it will automatically
remove the data which exceeds ttl.
*/
func (ele *inMemoryKeyValueDatastore) AutoCleanUp(checkInterval time.Duration, done <-chan bool) {
	ticker := time.NewTicker(checkInterval)
	for {
		select {
		case <-ticker.C:
			ele.GetAllKeyValues()
		case <-done:
			fmt.Println("Closing the AutoCleanUp")
			return
		}
	}
}
