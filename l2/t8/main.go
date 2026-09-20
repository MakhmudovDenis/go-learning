// Универсальный потокобезопасный кэш с TTL, очисткой и JSON-сериализацией :wq
package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	val       interface{}
	ttl       time.Duration
	addedAt   time.Time
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
		val:     value,
		ttl:     ttl,
		addedAt: time.Now(),
	}
}

func (c *Cache) Get(key string) (val interface{}, ok bool) {
	defer c.RUnlock()
	if c.Exists(key) {
		c.RLock()
		val, ok = c.m[key]
		c.RUnlock()
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

	if v.addedAt.Add(v.ttl).After(time.Now()) {
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
	jsonMap, err := json.Marshal(c.m)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func (c *Cache) GetAs[T any](key string) (T, error) {
	val, ok := c.Get(key)
	if !ok {
		return errors.New("Значение не существует или просрочено")
	}

	retVal, ok := val.(T)
	if ok!= nil {
		return nil, errors.New("Значение не приводится к типу T")
	}

	return retVal, nil
}

//
