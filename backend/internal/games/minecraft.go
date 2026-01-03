package games

import (
	"context"
)

// MinecraftHandler handles Minecraft-specific operations.
// Uses itzg/docker-minecraft-server image.
type MinecraftHandler struct{}

// GetDefaultImage returns the default Minecraft server image.
func (h *MinecraftHandler) GetDefaultImage() string {
	return "itzg/minecraft-server:latest"
}

// GetDefaultPorts returns the default Minecraft port mappings.
func (h *MinecraftHandler) GetDefaultPorts() map[string]int {
	return map[string]int{
		"minecraft": 25565,
		"rcon":      25575,
	}
}

// GetDefaultEnv returns the default Minecraft environment variables.
func (h *MinecraftHandler) GetDefaultEnv() map[string]string {
	return map[string]string{
		"EULA":          "TRUE",
		"TYPE":          "VANILLA",
		"VERSION":       "LATEST",
		"MEMORY":        "2G",
		"OPS":           "",
		"ENABLE_RCON":   "true",
		"RCON_PASSWORD": "changeme",
	}
}

// ValidateConfig validates Minecraft-specific configuration.
func (h *MinecraftHandler) ValidateConfig(config map[string]any) error {
	// TODO: Implement Minecraft-specific validation
	return nil
}

// OnStart is called when a Minecraft server is started.
func (h *MinecraftHandler) OnStart(_ context.Context, _ string) error {
	// TODO: Implement Minecraft-specific startup logic
	return nil
}

// OnStop is called when a Minecraft server is stopped.
func (h *MinecraftHandler) OnStop(_ context.Context, _ string) error {
	// TODO: Implement Minecraft-specific shutdown logic (e.g., graceful save)
	return nil
}

func init() {
	Register("minecraft", &MinecraftHandler{})
}
