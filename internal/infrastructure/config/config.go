// Package config provides configuration management.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application configuration.
type Config struct {
	// DefaultDateRange is the default date range preset.
	DefaultDateRange string `json:"default_date_range"`
	// RepoFilter is the repository filter pattern (glob).
	RepoFilter string `json:"repo_filter"`
	// OutputFormat is the output format: "text", "markdown", "json".
	OutputFormat string `json:"output_format"`
	// CustomTemplate is a custom template for output.
	CustomTemplate string `json:"custom_template"`
	// AutoCopy enables automatic copying to clipboard.
	AutoCopy bool `json:"auto_copy"`
	// ShowStats enables statistics display.
	ShowStats bool `json:"show_stats"`
}

// Default returns a config with default values.
func Default() Config {
	return Config{
		DefaultDateRange: "today",
		RepoFilter:       "",
		OutputFormat:     "text",
		CustomTemplate:   "",
		AutoCopy:         false,
		ShowStats:        true,
	}
}

// Path returns the path to the config file.
func Path() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "commitsum", "config.json"), nil
}

// Load loads configuration from file or returns defaults.
func Load() Config {
	configPath, err := Path()
	if err != nil {
		return Default()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Default()
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default()
	}

	return cfg
}

// Save saves configuration to file.
func Save(cfg Config) error {
	configPath, err := Path()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist.
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
