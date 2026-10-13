// Конфигуратор приложения с `sync.Once`
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ConfigManager struct {
	config  map[string]string
	initSrc string
	once    sync.Once
}

var initSources = [3]string{"env", "file", "db"}

func NewConfigManager(initSrc string) (*ConfigManager, error) {
	validSrc := false
	for _, src := range initSources {
		if initSrc == src {
			validSrc = true
			break
		}
	}

	if !validSrc {
		return nil, fmt.Errorf("invalid init source: %s\nexpected one of %v", initSrc, initSources)
	}

	return &ConfigManager{
		initSrc: initSrc,
	}, nil
}

func (cm *ConfigManager) loadConfigFromEnv() error {
	time.Sleep(1 * time.Second)
	if rand.Intn(2) == 0 {
		cm.config = map[string]string{
			"app_name":  "MyApp",
			"port":      "8080",
			"log_level": "debug",
		}
		return nil
	}

	return errors.New("loading config from env failed")
}

func (cm *ConfigManager) loadConfigFromFile() error {
	time.Sleep(2 * time.Second)
	if rand.Intn(2) == 0 {
		cm.config = map[string]string{
			"app_name":  "MyApp",
			"port":      "8080",
			"log_level": "debug",
		}
		return nil
	}

	return errors.New("loading config from file failed")
}

func (cm *ConfigManager) loadConfigFromDB() error {
	time.Sleep(3 * time.Second)
	if rand.Intn(2) == 0 {
		cm.config = map[string]string{
			"app_name":  "MyApp",
			"port":      "8080",
			"log_level": "debug",
		}
		return nil
	}

	return errors.New("loading config from db failed")
}

func (cm *ConfigManager) loadConfig() error {
	errs := make([]error, 0)
	cm.once.Do(func() {
		loadConfigFunctions := map[string]func() error{
			"env":  cm.loadConfigFromEnv,
			"file": cm.loadConfigFromFile,
			"db":   cm.loadConfigFromDB,
		}

		err := loadConfigFunctions[cm.initSrc]()
		if err == nil {
			return
		}

		errs = append(errs, err)

		for src, loadConfigFunc := range loadConfigFunctions {
			if src == cm.initSrc {
				continue
			}
			err := loadConfigFunc()
			if err == nil {
				return
			}
			errs = append(errs, err)
		}
	})
	if len(errs) != 0 {
		return fmt.Errorf("config loading failed with errors: %v", errs)
	}
	return nil
}

func (cm *ConfigManager) Get(key string) (string, error) {
	if err := cm.loadConfig(); err != nil {
		return "", err
	}

	return cm.config[key], nil
}

func (cm *ConfigManager) PrintConfig() {
	if err := cm.loadConfig(); err != nil {
		fmt.Println(err)
		return
	}

	jsonConfig, err := json.Marshal(cm.config)
	if err != nil {
		fmt.Printf("marshalling json failed with error: %v", err)
		return
	}

	fmt.Println(string(jsonConfig))
}

func main() {
	cm, err := NewConfigManager("env")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cm.Get("port"))
	cm.PrintConfig()
}
