// concurrent cache
package main

import (
	"fmt"
	"sync"
)

type SafeCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewSafeCache() *SafeCache {
	return &SafeCache{
		data: make(map[string]string),
	}
}

func (c *SafeCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *SafeCache) Get(key string) (val string, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok = c.data[key]
	return
}

func reader(wg *sync.WaitGroup, id int, cache *SafeCache) {
	defer wg.Done()
	for i := range 1000 {
		key := fmt.Sprintf("key-%d", i%100)
		v, ok := cache.Get(key)

		if !ok {
			continue
		}

		if i%200 == 0 {
			fmt.Printf("-%d: %s = %s\n", id, key, v)
		}
	}
}

func writer(wg *sync.WaitGroup, id int, cache *SafeCache) {
	defer wg.Done()
	for i := range 1000 {
		key := fmt.Sprintf("key-%d", i%100)
		value := fmt.Sprintf("writer-%d", id)
		cache.Set(key, value)
	}
}

func main() {
	var wg sync.WaitGroup
	cache := NewSafeCache()

	for i := range 100 {
		key := fmt.Sprintf("key-%d", i)
		cache.Set(key, "initial")
	}

	wg.Add(200)

	for i := range 100 {
		go writer(&wg, i, cache)
		go reader(&wg, i, cache)
	}
	wg.Wait()
}
