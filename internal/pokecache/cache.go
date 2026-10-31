package pokecache

import (
	"bytes"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	if interval <= 0 {
		panic("cache interval must be positive")
	}
	cache := &Cache{entries: make(map[string]cacheEntry)}
	go cache.reapLoop(interval)
	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.entries[key] = cacheEntry{createdAt: time.Now(), val: bytes.Clone(val)}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entry, ok := cache.entries[key]
	if !ok {
		return nil, false
	}
	return bytes.Clone(entry.val), true
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for now := range ticker.C {
		cache.mu.Lock()
		for key, entry := range cache.entries {
			if now.Sub(entry.createdAt) > interval {
				delete(cache.entries, key)
			}
		}
		cache.mu.Unlock()
	}
}
