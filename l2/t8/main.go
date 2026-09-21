// Универсальный потокобезопасный кэш с TTL, очисткой и JSON-сериализацией :wq
package main

import (
	"encoding/json"
	_ "errors"
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	val       interface{}
	expiresAt time.Time
}

type Cache struct {
	m map[string]cacheEntry
	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{m: make(map[string]cacheEntry)}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	c.m[key] = cacheEntry{
		val:       value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (val interface{}, ok bool) {
	if c.Exists(key) {
		c.RLock()
		entry := c.m[key]
		c.RUnlock()
		val = entry.val
		ok = true
	}

	return
}

func (c *Cache) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.m, key)
}

func (c *Cache) Exists(key string) bool {
	c.RLock()
	v, ok := c.m[key]
	c.RUnlock()
	if !ok {
		return false
	}

	if v.expiresAt.Before(time.Now()) {
		return false
	}

	return true
}

func (c *Cache) Clear() {
	c.Lock()
	defer c.Unlock()

	for k := range c.m {
		delete(c.m, k)
	}
}

func (c *Cache) ToJSON() ([]byte, error) {
	tempMap := make(map[string]interface{})
	for key := range c.m {
		tempMap[key] = c.m[key].val
	}

	fmt.Println(tempMap)

	jsonMap, err := json.Marshal(tempMap)
	if err != nil {
		return nil, fmt.Errorf("ошибка кодировния: %w", err)
	}

	return jsonMap, nil
}

func main() {
	cache := NewCache()

	var t int8 = 2

	cache.Set("test", t, 1*time.Hour)
	cache.Set("test2", "asd", 1*time.Minute)
	fmt.Println(cache.Get("test"))
	cache.Set("a", struct{ B, a int }{a: 4, B: 1}, time.Second)
	fmt.Println(cache.Get("a"))
	time.Sleep(2 * time.Second)
	fmt.Println(cache.Exists("a"))
	json, err := cache.ToJSON()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(json))
	}
	cache.Clear()
	fmt.Println(cache.Exists("test2"))
}
