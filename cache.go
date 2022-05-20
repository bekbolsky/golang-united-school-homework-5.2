package cache

import (
	"sync"
	"time"
)

// Cache holds values with type string and
// allows to retrieve them using keys of time strings.
// Key/value pairs can expire if given a deadline using PutTill method.
type Cache struct {
	RWMutex  sync.RWMutex
	cacheMap map[string]cacheItem
}

// cacheItem holds the value and deadline of a key/value pair
type cacheItem struct {
	value    string
	deadline time.Time
	expired  bool
}

// NewCache returns a new cache instance
func NewCache() Cache {
	cachemap := make(map[string]cacheItem)
	return Cache{
		cacheMap: cachemap,
	}
}

// Get returns the value associated with the key and the boolean ok (true if exists, false if not),
// if the deadline of the key/value pair has not been exceeded yet.
func (c *Cache) Get(key string) (string, bool) {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()

	if item, ok := c.cacheMap[key]; ok {
		item.expired = item.deadline.Before(time.Now())
		if !item.expired || item.deadline.IsZero() {
			return item.value, true
		}
	}
	return "", false
}

// Put places a value with an associated key into cache.
// Value put with this method never expired (have infinite deadline).
// Putting into the existing key overwrites the value
func (c *Cache) Put(key, value string) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	c.cacheMap[key] = cacheItem{
		value:    value,
		deadline: time.Time{},
		expired:  false,
	}
}

// Keys returns a slice of keys in the cache (not expired)
func (c *Cache) Keys() []string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()

	keys := make([]string, 0, len(c.cacheMap))
	for key, item := range c.cacheMap {
		item.expired = item.deadline.Before(time.Now())
		if !item.expired || item.deadline.IsZero() {
			keys = append(keys, key)
		}
	}
	return keys
}

// PutTill places a value with an associated key into cache.
// Value put with this method expires after the given deadline.
// Putting into the existing key overwrites the value
func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	c.cacheMap[key] = cacheItem{
		value:    value,
		deadline: deadline,
		expired:  false,
	}
}
