package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestNewKeyValueDataStore(t *testing.T) {
	key, value := "DOB", "26/07/1996"
	ttl := time.Second
	ds := NewKeyValueCache[string](ttl)
	ds.Put(key, value)

	val, ok := ds.Get(key)

	if !ok {
		t.Errorf("expected to get true, got false")
	}

	if val.Value != value {
		t.Errorf("expected value to be %v, got %v", value, val.Value)
	}

	if val.TTL != ttl {
		t.Errorf("expected TTL to be %v, got %v", ttl, val.TTL)
	}
}

func TestPutAndGet(t *testing.T) {
	ds := NewKeyValueCache[string](5 * time.Second)
	key, value := "username", "user123"

	// Test putting and getting a value
	ds.Put(key, value)
	val, ok := ds.Get(key)

	if !ok {
		t.Fatalf("expected to get true, got false")
	}

	if val.Value != value {
		t.Errorf("expected value to be %v, got %v", value, val.Value)
	}
}

func TestGetWithTTLExpiry(t *testing.T) {
	ds := NewKeyValueCache[string](1 * time.Second)
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
	ds := NewKeyValueCache[string](5 * time.Second)
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
	ds := NewKeyValueCache[string](1 * time.Second)
	key, value := "tempKey", "tempValue"
	ds.Put(key, value)

	done := make(chan bool)
	go ds.AutoCleanUp(500*time.Millisecond, done)

	// Wait for TTL to expire and cleanup to run
	time.Sleep(2 * time.Second)

	// Verify that the key has been cleaned up
	_, ok := ds.Get(key)
	if ok {
		t.Errorf("expected key to be cleaned up, but it still exists")
	}

	done <- true
}

func TestConcurrentAccess(t *testing.T) {
	ds := NewKeyValueCache[string](5 * time.Second)

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
		if !ok || val.Value != fmt.Sprintf("value%d", i) {
			t.Errorf("expected value %v, got %v", fmt.Sprintf("value%d", i), val.Value)
		}
	}
}
