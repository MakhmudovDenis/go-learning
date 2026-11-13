// Конфигуратор приложения с `sync.Once`
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type loader func() (map[string]string, error)

type source struct {
	name string
	load loader
}

func fakeLoader(name string, delay time.Duration) loader {
	return func() (map[string]string, error) {
		time.Sleep(delay)
		if rand.Intn(2) == 0 {
			return map[string]string{
				"app_name":  "MyApp",
				"port":      "8080",
				"log_level": "debug",
			}, nil
		}
		return nil, fmt.Errorf("loading config from %s failed", name)
	}
}

var defaultSources = []source{
	{"env", fakeLoader("env", 1*time.Second)},
	{"file", fakeLoader("file", 2*time.Second)},
	{"db", fakeLoader("db", 3*time.Second)},
}

type ConfigManager struct {
	sources []source
	once    sync.Once
	config  map[string]string
	loadErr error
}

func NewConfigManager(preferred string) (*ConfigManager, error) {
	ordered := make([]source, 0, len(defaultSources))
	for _, s := range defaultSources {
		if s.name == preferred {
			ordered = append(ordered, s)
		}
	}
	if len(ordered) == 0 {
		names := make([]string, len(defaultSources))
		for i, s := range defaultSources {
			names[i] = s.name
		}
		return nil, fmt.Errorf("invalid init source %q, expected one of %v", preferred, names)
	}
	for _, s := range defaultSources {
		if s.name != preferred {
			ordered = append(ordered, s)
		}
	}

	return &ConfigManager{sources: ordered}, nil
}

func (cm *ConfigManager) LoadConfig() error {
	cm.once.Do(func() {
		cm.config, cm.loadErr = cm.loadFirstAvailable()
	})
	return cm.loadErr
}

func (cm *ConfigManager) loadFirstAvailable() (map[string]string, error) {
	var errs []error
	for _, s := range cm.sources {
		cfg, err := s.load()
		if err == nil {
			return cfg, nil
		}
		errs = append(errs, err)
	}
	return nil, fmt.Errorf("all config sources failed: %w", errors.Join(errs...))
}

func (cm *ConfigManager) Get(key string) string {
	if cm.LoadConfig() != nil {
		return ""
	}
	return cm.config[key]
}

func (cm *ConfigManager) PrintConfig() {
	if err := cm.LoadConfig(); err != nil {
		fmt.Println(err)
		return
	}

	jsonConfig, err := json.Marshal(cm.config)
	if err != nil {
		fmt.Printf("marshalling json failed with error: %v\n", err)
		return
	}

	fmt.Println(string(jsonConfig))
}

func main() {
	cm, err := NewConfigManager("env")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := range 5 {
		wg.Go(func() {
			if err := cm.LoadConfig(); err != nil {
				fmt.Printf("goroutine %d: %v\n", i, err)
				return
			}
			fmt.Printf("goroutine %d: port=%s\n", i, cm.Get("port"))
		})
	}
	wg.Wait()

	cm.PrintConfig()
}
