package games

import (
	"context"
)

// RustHandler handles Rust game server-specific operations.
// Uses max-pfeiffer/rust-game-server image.
type RustHandler struct{}

// GetDefaultImage returns the default Rust game server image.
func (h *RustHandler) GetDefaultImage() string {
	return "maxpfeiffer/rust-game-server:latest"
}

// GetDefaultPorts returns the default Rust port mappings.
func (h *RustHandler) GetDefaultPorts() map[string]int {
	return map[string]int{
		"game": 28015,
		"rcon": 28016,
		"app":  28082,
	}
}

// GetDefaultEnv returns the default Rust environment variables.
func (h *RustHandler) GetDefaultEnv() map[string]string {
	return map[string]string{
		"SERVER_NAME":        "Sabakan Rust Server",
		"SERVER_DESCRIPTION": "A Rust server managed by Sabakan",
		"SERVER_HOSTNAME":    "0.0.0.0",
		"SERVER_SEED":        "",
		"SERVER_WORLDSIZE":   "3000",
		"SERVER_MAXPLAYERS":  "50",
		"SERVER_IDENTITY":    "default",
		"RCON_PASSWORD":      "changeme",
		"RCON_WEB":           "1",
		"APP_PORT":           "28082",
	}
}

// ValidateConfig validates Rust-specific configuration.
func (h *RustHandler) ValidateConfig(config map[string]any) error {
	// TODO: Implement Rust-specific validation
	return nil
}

// OnStart is called when a Rust server is started.
func (h *RustHandler) OnStart(_ context.Context, _ string) error {
	// TODO: Implement Rust-specific startup logic
	return nil
}

// OnStop is called when a Rust server is stopped.
func (h *RustHandler) OnStop(_ context.Context, _ string) error {
	// TODO: Implement Rust-specific shutdown logic (save world)
	return nil
}

func init() {
	Register("rust", &RustHandler{})
}
