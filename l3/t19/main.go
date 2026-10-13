// JSON-кэш с `sync.Pool` и `map + RWMutex`
package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	value     any
	expiresAt time.Time
}

type SafeCache struct {
	ttl  time.Duration
	mu   sync.RWMutex
	data map[string]*CacheEntry
	pool sync.Pool
}

func NewSafeCache(TTL time.Duration) *SafeCache {
	return &SafeCache{
		data: make(map[string]*CacheEntry),
		ttl:  TTL,
		pool: sync.Pool{
			New: func() any {
				return make(map[string]any)
			},
		},
	}
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
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
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
	c.mu.RLock()
	defer c.mu.RUnlock()

	tempMap := c.pool.Get().(map[string]any)
	defer func() {
		clear(tempMap)
		c.pool.Put(tempMap)
	}()

	for key, val := range c.data {
		tempMap[key] = val.value
	}

	jsonData, err := json.Marshal(tempMap)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func main() {
	c := NewSafeCache(5 * time.Second)
	c.Set("hehe", 123)
	fmt.Println(c.Get("hehe"))

	j, err := c.ToJSON()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(j))
}
