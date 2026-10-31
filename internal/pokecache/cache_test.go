package pokecache

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(time.Minute)
	if val, ok := cache.Get("missing"); ok || val != nil {
		t.Fatal("expected a cache miss")
	}
	cases := []struct {
		key string
		val []byte
	}{
		{key: "https://example.com", val: []byte("testdata")},
		{key: "https://example.com/path", val: []byte("moretestdata")},
		{key: "https://example.com", val: []byte("replacement")},
	}
	for _, testCase := range cases {
		cache.Add(testCase.key, testCase.val)
		val, ok := cache.Get(testCase.key)
		if !ok || !bytes.Equal(val, testCase.val) {
			t.Errorf("Get(%q) = %q, %v; want %q, true", testCase.key, val, ok, testCase.val)
		}
	}
}

func TestCacheOwnsBytes(t *testing.T) {
	cache := NewCache(time.Minute)
	input := []byte("original")
	cache.Add("key", input)
	input[0] = 'X'
	val, _ := cache.Get("key")
	if string(val) != "original" {
		t.Fatal("mutating input changed cached bytes")
	}
	val[0] = 'Y'
	val, _ = cache.Get("key")
	if string(val) != "original" {
		t.Fatal("mutating a result changed cached bytes")
	}
}

func TestReapLoop(t *testing.T) {
	const interval = 10 * time.Millisecond
	cache := NewCache(interval)
	cache.mu.Lock()
	cache.entries["old"] = cacheEntry{createdAt: time.Now().Add(-time.Hour), val: []byte("old")}
	cache.mu.Unlock()
	deadline := time.NewTimer(time.Second)
	defer deadline.Stop()
	poll := time.NewTicker(time.Millisecond)
	defer poll.Stop()
	for {
		select {
		case <-poll.C:
			if _, ok := cache.Get("old"); !ok {
				cache.Add("fresh", []byte("fresh"))
				return
			}
		case <-deadline.C:
			t.Fatal("reaper did not remove the expired entry")
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewCache(time.Millisecond)
	var workers sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		workers.Go(func() {
			for index := 0; index < 100; index++ {
				key := fmt.Sprintf("key-%d", index)
				cache.Add(key, []byte("value"))
				cache.Get(key)
			}
		})
	}
	workers.Wait()
}
