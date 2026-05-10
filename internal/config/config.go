package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	DatabaseURL      string
	ExecutorCapacity int
}

var (
	ENV  *Config = nil
	once sync.Once
)

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.ExecutorCapacity <= 0 {
		return fmt.Errorf("EXECUTOR_CAPACITY must be a positive integer")
	}

	return nil
}

func Load() {
	once.Do(func() {
		dbURL := os.Getenv("DATABASE_URL")
		capacityStr := os.Getenv("EXECUTOR_CAPACITY")
		if capacityStr == "" {
			capacityStr = "1024"
		}

		capacity, err := strconv.Atoi(capacityStr)
		if err != nil && capacityStr != "" {
			panic("invalid EXECUTOR_CAPACITY: " + err.Error())
		}

		cfg := &Config{
			DatabaseURL:      dbURL,
			ExecutorCapacity: capacity,
		}

		if err := cfg.Validate(); err != nil {
			panic("invalid config: " + err.Error())
		}

		ENV = cfg
	})
}

func LoadTest(executorCapacity int) {
	once.Do(func() {
		if ENV != nil {
			return
		}

		ENV = &Config{
			DatabaseURL:      "",
			ExecutorCapacity: executorCapacity,
		}
	})
}
