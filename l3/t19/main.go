// JSON-кэш с `sync.Pool` и `map + RWMutex`
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	value     any
	expiresAt time.Time
}

func (e *CacheEntry) expired(now time.Time) bool {
	return now.After(e.expiresAt)
}

var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

type SafeCache struct {
	ttl  time.Duration
	mu   sync.RWMutex
	data map[string]*CacheEntry

	done      chan struct{}
	closeOnce sync.Once
}

func NewSafeCache(ttl time.Duration) *SafeCache {
	c := &SafeCache{
		ttl:  ttl,
		data: make(map[string]*CacheEntry),
		done: make(chan struct{}),
	}

	go c.janitor(ttl / 2)

	return c
}

func (c *SafeCache) janitor(every time.Duration) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.done:
			return
		}
	}
}

func (c *SafeCache) deleteExpired() {
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.data {
		if entry.expired(now) {
			delete(c.data, key)
		}
	}
}

func (c *SafeCache) Close() {
	c.closeOnce.Do(func() { close(c.done) })
}

func (c *SafeCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &CacheEntry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *SafeCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists || entry.expired(time.Now()) {
		return nil, false
	}

	return entry.value, true
}

func (c *SafeCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *SafeCache) ToJSON() ([]byte, error) {
	now := time.Now()
	c.mu.RLock()
	snapshot := make(map[string]any, len(c.data))
	for key, entry := range c.data {
		if !entry.expired(now) {
			snapshot[key] = entry.value
		}
	}
	c.mu.RUnlock()

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(snapshot); err != nil {
		return nil, err
	}

	return bytes.Clone(bytes.TrimSuffix(buf.Bytes(), []byte("\n"))), nil
}

func main() {
	cache := NewSafeCache(2 * time.Second)
	defer cache.Close()

	cache.Set("user:1", map[string]string{"name": "Alice", "role": "admin"})
	cache.Set("user:2", map[string]string{"name": "Bob", "role": "user"})

	if user, found := cache.Get("user:1"); found {
		fmt.Println("Найден:", user)
	}

	jsonData, err := cache.ToJSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Кэш в JSON:", string(jsonData))

	time.Sleep(3500 * time.Millisecond)

	_, found := cache.Get("user:1")
	fmt.Println("После TTL, user:1 найден?", found)

	cache.mu.RLock()
	fmt.Println("Записей в map после очистки:", len(cache.data)) // 0: really deleted, not just hidden
	cache.mu.RUnlock()
}
