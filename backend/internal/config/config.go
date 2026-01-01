// Package config provides configuration loading and management for Sabakan.
// It supports system-wide configuration and per-game configuration using TOML format.
package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

// SystemConfig represents the system-wide configuration.
type SystemConfig struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	Logging  LoggingConfig  `toml:"logging"`
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

// DatabaseConfig contains database connection settings.
type DatabaseConfig struct {
	Path string `toml:"path"`
}

// LoggingConfig contains logging settings.
type LoggingConfig struct {
	Level  string `toml:"level"`  // debug, info, warn, error
	Format string `toml:"format"` // json, text
}

// GameConfig represents per-game configuration.
type GameConfig struct {
	Game      GameInfo      `toml:"game"`
	Container ContainerInfo `toml:"container"`
	Mods      []ModInfo     `toml:"mods"`
}

// GameInfo contains basic game metadata.
type GameInfo struct {
	ID          string `toml:"id"`
	Name        string `toml:"name"`
	Description string `toml:"description"`
}

// ContainerInfo contains container runtime configuration.
type ContainerInfo struct {
	Image   string            `toml:"image"`
	Ports   []string          `toml:"ports"`
	Volumes []string          `toml:"volumes"`
	Env     map[string]string `toml:"env"`
}

// ModInfo represents a single mod entry.
type ModInfo struct {
	ID      string `toml:"id"`
	Name    string `toml:"name"`
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`
}

// LoadSystemConfig loads the system configuration from a TOML file.
func LoadSystemConfig(path string) (*SystemConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg SystemConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// DefaultSystemConfig returns the default system configuration.
func DefaultSystemConfig() *SystemConfig {
	return &SystemConfig{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 1323,
		},
		Database: DatabaseConfig{
			Path: "./sabakan.db",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

// LoadGameConfig loads a game configuration from a TOML file.
func LoadGameConfig(path string) (*GameConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg GameConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveGameConfig saves a game configuration to a TOML file.
func SaveGameConfig(path string, cfg *GameConfig) error {
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
