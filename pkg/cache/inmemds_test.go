package cache

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewKeyValueDataStore(t *testing.T) {
	key, value := "DOB", "26/07/1996"
	ttl := time.Second
	sweep := 3 * time.Second
	ds, err := NewKeyValueCache[string](context.Background(), ttl, WithSweeping(sweep))
	if err != nil {
		t.Fatalf("initialize failed")
	}
	ds.Put(key, value)

	val, ok := ds.Get(key)

	if !ok {
		t.Errorf("expected to get true, got false")
	}

	if val != value {
		t.Errorf("expected value to be %v, got %v", value, val)
	}

	if ds.ttl != ttl {
		t.Errorf("expected TTL to be %v, got %v", ttl, ds.ttl)
	}
}

func TestPutAndGet(t *testing.T) {
	ds, err := NewKeyValueCache[string](context.Background(), 5*time.Second)
	if err != nil {
		t.Fatalf("initialize failed")
	}
	key, value := "username", "user123"

	// Test putting and getting a value
	ds.Put(key, value)
	val, ok := ds.Get(key)

	if !ok {
		t.Fatalf("expected to get true, got false")
	}

	if val != value {
		t.Errorf("expected value to be %v, got %v", value, val)
	}
}

func TestGetWithTTLExpiry(t *testing.T) {
	ds, err := NewKeyValueCache[string](context.Background(), 1*time.Second)
	if err != nil {
		t.Fatalf("initialize failed")
	}
	key, value := "sessionID", "abc123"

	// Put value
	ds.Put(key, value)

	// Wait for TTL to expire
	time.Sleep(2 * time.Second)

	// Get should return false due to TTL expiration
	_, ok := ds.Get(key)
	if ok {
		t.Errorf("expected to get false, got true after TTL expiry")
	}
}

func TestGetAllKeyValues(t *testing.T) {
	ds, err := NewKeyValueCache[string](context.Background(), 5*time.Second)
	if err != nil {
		t.Fatalf("cache init failed")
	}
	ds.Put("key1", "value1")
	ds.Put("key2", "value2")
	ds.Put("key3", "value3")

	allRecords := ds.GetAllKeyValues()
	expectedLen := 3

	if len(allRecords) != expectedLen {
		t.Errorf("expected %d records, got %d", expectedLen, len(allRecords))
	}

	if allRecords["key1"] != "value1" || allRecords["key2"] != "value2" || allRecords["key3"] != "value3" {
		t.Errorf("unexpected values in GetAllKeyValues result")
	}
}

func TestAutoCleanUp(t *testing.T) {
	sweepInterval := 500 * time.Millisecond
	ds, err := NewKeyValueCache[string](context.Background(), 1*time.Second, WithSweeping(sweepInterval))
	if err != nil {
		t.Fatalf("init failed")
	}
	key, value := "tempKey", "tempValue"
	ds.Put(key, value)

	// Wait for TTL to expire and cleanup to run
	time.Sleep(2 * time.Second)

	// Verify that the key has been cleaned up
	_, ok := ds.Get(key)
	if ok {
		t.Errorf("expected key to be cleaned up, but it still exists")
	}
}

func TestConcurrentAccess(t *testing.T) {
	ds, err := NewKeyValueCache[string](context.Background(), 5*time.Second)
	if err != nil {
		t.Fail()
	}

	// Run concurrent writes using integer range
	for i := range make([]struct{}, 100) {
		go ds.Put(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	// Run concurrent reads using integer range
	for i := range make([]struct{}, 100) {
		go func(i int) {
			ds.Get(fmt.Sprintf("key%d", i))
		}(i)
	}

	// Ensure all values are accessible
	time.Sleep(1 * time.Second)
	for i := range make([]struct{}, 100) {
		val, ok := ds.Get(fmt.Sprintf("key%d", i))
		if !ok || val != fmt.Sprintf("value%d", i) {
			t.Errorf("expected value %v, got %v", fmt.Sprintf("value%d", i), val)
		}
	}
}
