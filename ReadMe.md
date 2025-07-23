# Go Generic In-Memory Datastore 🗃️

A lightweight, concurrent-safe, generic in-memory datastore with TTL support and optional background sweeping for expired records.

## ✨ Features

- ✅ Type-safe generics for both keys and values
- 🕒 TTL (Time-To-Live) support per record
- 🔄 Optional sweeping goroutine to clear expired items
- 🔐 Thread-safe using RWMutex
- ⚙️ Flexible options via functional configuration

## 📦 Installation

Simply copy the package into your project or import from your module:

```go
import "your_project/datastore"
```

## 🔧 Usage Example

```go
ctx := context.Background()

// Create a new datastore with 5-minute TTL and sweeping every 1 minute
ds, err := datastore.NewDatastore[string, int](
	ctx,
	5 * time.Minute,
	datastore.WithSweeping(1 * time.Minute),
)
if err != nil {
	log.Fatal(err)
}

// Put data
ds.Put("user:123", 42)

// Get data
item, ok := ds.Get("user:123")
if ok {
	fmt.Println("Value:", item.Value)
}
```

## 🔍 API Overview

### `NewDatastore[K, T](ctx, ttl, opts...)`
Creates a new datastore with optional sweeping interval.

### `Put(key, value)`
Stores a value with the given key and TTL.

### `Get(key)`
Returns the value and whether it's still valid (not expired).

### `GetAllKeyValues()`
Returns all non-expired key-value pairs.

### `WithSweeping(duration)`
Optional background job that clears expired entries at regular intervals.

## 📁 Structure

- `Datastore[K, T]` - Main generic structure
- `CacheItem[T]` - Holds value, TTL, and creation time
- `Options` - Configurable behaviors like sweeping

## 🧪 Example Use-Cases

- Caching intermediate results in web services
- Session/token storage
- Lightweight in-memory KV stores

## 🛡️ License

MIT

---
Made with ❤️ in Go.
