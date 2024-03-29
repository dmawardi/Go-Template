package cache

import (
	"sync"
	"time"
)

// Entry represents a cache entry with a value and an expiration timestamp
type Entry struct {
	Value      interface{}
	Expiration int64
}

// CacheMap wraps sync.Map and provides basic TTL functionality
// Example usage: m.Store("key", "value", 10 * time.Second)
type CacheMap struct {
	internal sync.Map
}

// Store adds a value to the map with a specified TTL (in seconds)
// Example usage: m.Store("key", "value", 10 * time.Second)
func (m *CacheMap) Store(key, value interface{}, ttl time.Duration) {
	// Create expiration timestamp
	expiration := time.Now().Add(ttl).UnixNano()
	// Build entry
	entry := Entry{Value: value, Expiration: expiration}
	// Store entry
	m.internal.Store(key, entry)
}

// Load retrieves a value from the map, considering its TTL
// Example usage: value, ok := m.Load("key")
func (m *CacheMap) Load(key interface{}) (interface{}, bool) {
	// Load entry using key
	result, ok := m.internal.Load(key)
	if !ok {
		return nil, false
	}
	// If found,
	entry := result.(Entry)
	// check if expired
	if time.Now().UnixNano() > entry.Expiration {
		// If expired, delete entry and return false
		m.internal.Delete(key) // Remove expired entry
		return nil, false
	}
	// If not expired, return value and true
	return entry.Value, true
}
