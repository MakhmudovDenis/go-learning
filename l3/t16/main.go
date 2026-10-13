// Инициализация плагинов с `sync.Once`
package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var errPluginNotExist = errors.New("plugin does not exist")

type Plugin interface {
	Execute() string
}

type PluginManager struct {
	plugins map[string]*pluginEntry
	mu      sync.RWMutex
}

type pluginEntry struct {
	once   sync.Once
	plugin Plugin
	initFn func() (Plugin, error)
	err    error
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*pluginEntry),
	}
}

func (pm *PluginManager) RegisterPlugin(name string, initFn func() (Plugin, error)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[name] = &pluginEntry{
		initFn: initFn,
	}
}

func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	entry, exists := pm.plugins[name]
	if !exists {
		return nil, errPluginNotExist
	}

	if entry.err != nil {
		fmt.Printf("cached error for %s:%v\n", name, entry.err)
		return nil, entry.err
	}

	entry.once.Do(func() {
		plugin, err := entry.initFn()
		if err != nil {
			entry.err = err
			return
		}

		entry.plugin = plugin
	})

	return entry.plugin, entry.err
}

type DemoPlugin struct{}

func (p *DemoPlugin) Execute() string {
	return "DemoPlugin executed successfully!"
}

func initDemo() (Plugin, error) {
	time.Sleep(500 * time.Millisecond)
	return &DemoPlugin{}, nil
}

func main() {
	pm := NewPluginManager()

	pm.RegisterPlugin("demo", initDemo)
	pm.RegisterPlugin("broken", func() (Plugin, error) {
		return nil, fmt.Errorf("simulated error")
	})

	var wg sync.WaitGroup

	// Тестирование рабочего плагина
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p, err := pm.GetPlugin("demo")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
				return
			}
			log.Printf("Goroutine %d: %s", id, p.Execute())
		}(i)
	}

	// Тестирование плагина с ошибкой
	for i := 5; i < 7; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := pm.GetPlugin("broken")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
}
