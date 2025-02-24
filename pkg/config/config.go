// pkg/config/config.go
package config

import (
	"elastic-btree/pkg/logger"
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration.
type Config struct {
	TreeDegree  int          // B-tree degree
	LogLevel    logger.Level // Logging level (debug, info, warn, error)
	StoragePath string       // Path to the storage file
}

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	// Default values
	cfg := &Config{
		TreeDegree:  3,
		LogLevel:    logger.Info,
		StoragePath: "data/tree.json",
	}

	// Load TreeDegree from environment
	if degreeStr := os.Getenv("TREE_DEGREE"); degreeStr != "" {
		degree, err := strconv.Atoi(degreeStr)
		if err != nil || degree < 2 {
			return nil, fmt.Errorf("invalid TREE_DEGREE: %s (must be >= 2)", degreeStr)
		}
		cfg.TreeDegree = degree
	}

	// Load LogLevel from environment
	if logLevelStr := os.Getenv("LOG_LEVEL"); logLevelStr != "" {
		logLevel, err := logger.ParseLevel(logLevelStr)
		if err != nil {
			return nil, fmt.Errorf("invalid LOG_LEVEL: %s", logLevelStr)
		}
		cfg.LogLevel = logLevel
	}

	// Load StoragePath from environment
	if storagePath := os.Getenv("STORAGE_PATH"); storagePath != "" {
		cfg.StoragePath = storagePath
	}

	return cfg, nil
}
