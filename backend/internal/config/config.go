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
	Podman   PodmanConfig   `toml:"podman"`
	JWT      JWTConfig      `toml:"jwt"`
	Redis    RedisConfig    `toml:"redis"`
	Auth     AuthConfig     `toml:"auth"`
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

// PodmanConfig contains Podman connection settings.
type PodmanConfig struct {
	// SocketPath is the path to the Podman socket.
	// For rootful: unix:///run/podman/podman.sock
	// For rootless: unix://$XDG_RUNTIME_DIR/podman/podman.sock
	SocketPath string `toml:"socket_path"`
}

// JWTConfig contains JWT authentication settings.
type JWTConfig struct {
	Secret             string `toml:"secret"`               // Secret key for signing tokens
	AccessTokenExpiry  int    `toml:"access_token_expiry"`  // Access token expiry in minutes
	RefreshTokenExpiry int    `toml:"refresh_token_expiry"` // Refresh token expiry in days
}

// RedisConfig contains Redis connection settings.
type RedisConfig struct {
	URL string `toml:"url"` // Redis connection URL (redis://host:port)
}

// AuthConfig contains authentication settings.
type AuthConfig struct {
	AllowRegistration bool `toml:"allow_registration"` // Whether to allow new user registration
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
		Podman: PodmanConfig{
			SocketPath: "unix:///run/podman/podman.sock",
		},
		JWT: JWTConfig{
			Secret:             "change-this-secret-in-production-32bytes!",
			AccessTokenExpiry:  15, // 15 minutes
			RefreshTokenExpiry: 7,  // 7 days
		},
		Redis: RedisConfig{
			URL: "redis://localhost:6379",
		},
		Auth: AuthConfig{
			AllowRegistration: true,
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
