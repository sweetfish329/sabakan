// Package games provides game-specific handlers for different game servers.
package games

import (
	"context"
)

// GameHandler defines the interface for game-specific operations.
type GameHandler interface {
	// GetDefaultImage returns the default Docker image for this game.
	GetDefaultImage() string

	// GetDefaultPorts returns the default port mappings for this game.
	GetDefaultPorts() map[string]int

	// GetDefaultEnv returns the default environment variables for this game.
	GetDefaultEnv() map[string]string

	// ValidateConfig validates the game-specific configuration.
	ValidateConfig(config map[string]any) error

	// OnStart is called when a game server is started.
	OnStart(ctx context.Context, containerID string) error

	// OnStop is called when a game server is stopped.
	OnStop(ctx context.Context, containerID string) error
}

// Registry holds all registered game handlers.
var Registry = make(map[string]GameHandler)

// Register adds a game handler to the registry.
func Register(name string, handler GameHandler) {
	Registry[name] = handler
}

// Get retrieves a game handler by name.
func Get(name string) (GameHandler, bool) {
	handler, ok := Registry[name]
	return handler, ok
}

// List returns all registered game names.
func List() []string {
	names := make([]string, 0, len(Registry))
	for name := range Registry {
		names = append(names, name)
	}
	return names
}
