package games

import (
	"context"
)

// PalworldHandler handles Palworld-specific operations.
// Uses thijsvanloef/palworld-server-docker image.
type PalworldHandler struct{}

// GetDefaultImage returns the default Palworld server image.
func (h *PalworldHandler) GetDefaultImage() string {
	return "thijsvanloef/palworld-server-docker:latest"
}

// GetDefaultPorts returns the default Palworld port mappings.
func (h *PalworldHandler) GetDefaultPorts() map[string]int {
	return map[string]int{
		"game":  8211,
		"query": 27015,
		"rcon":  25575,
	}
}

// GetDefaultEnv returns the default Palworld environment variables.
func (h *PalworldHandler) GetDefaultEnv() map[string]string {
	return map[string]string{
		"PUID":               "1000",
		"PGID":               "1000",
		"PORT":               "8211",
		"PLAYERS":            "16",
		"MULTITHREADING":     "true",
		"RCON_ENABLED":       "true",
		"RCON_PORT":          "25575",
		"ADMIN_PASSWORD":     "changeme",
		"COMMUNITY":          "false",
		"SERVER_NAME":        "Sabakan Palworld Server",
		"SERVER_DESCRIPTION": "A Palworld server managed by Sabakan",
	}
}

// ValidateConfig validates Palworld-specific configuration.
func (h *PalworldHandler) ValidateConfig(config map[string]any) error {
	// TODO: Implement Palworld-specific validation
	return nil
}

// OnStart is called when a Palworld server is started.
func (h *PalworldHandler) OnStart(_ context.Context, _ string) error {
	// TODO: Implement Palworld-specific startup logic
	return nil
}

// OnStop is called when a Palworld server is stopped.
func (h *PalworldHandler) OnStop(_ context.Context, _ string) error {
	// TODO: Implement Palworld-specific shutdown logic
	return nil
}

func init() {
	Register("palworld", &PalworldHandler{})
}
