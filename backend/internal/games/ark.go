package games

import (
	"context"
)

// ArkHandler handles ARK: Survival Evolved-specific operations.
// Uses hermsi/ark-server image.
type ArkHandler struct{}

// GetDefaultImage returns the default ARK server image.
func (h *ArkHandler) GetDefaultImage() string {
	return "hermsi/ark-server:latest"
}

// GetDefaultPorts returns the default ARK port mappings.
func (h *ArkHandler) GetDefaultPorts() map[string]int {
	return map[string]int{
		"game":  7777,
		"query": 27015,
		"rcon":  32330,
	}
}

// GetDefaultEnv returns the default ARK environment variables.
func (h *ArkHandler) GetDefaultEnv() map[string]string {
	return map[string]string{
		"SESSION_NAME":      "Sabakan ARK Server",
		"SERVER_MAP":        "TheIsland",
		"SERVER_PASSWORD":   "",
		"ADMIN_PASSWORD":    "changeme",
		"MAX_PLAYERS":       "70",
		"UPDATE_ON_START":   "true",
		"BACKUP_ON_STOP":    "true",
		"WARN_ON_STOP":      "true",
		"ENABLE_CROSSPLAY":  "false",
		"DISABLE_BATTLEYE":  "false",
		"ARK_MODS":          "",
		"GAME_CLIENT_PORT":  "7777",
		"UDP_SOCKET":        "7778",
		"RCON_PORT":         "32330",
		"SERVER_QUERY_PORT": "27015",
	}
}

// ValidateConfig validates ARK-specific configuration.
func (h *ArkHandler) ValidateConfig(config map[string]any) error {
	// TODO: Implement ARK-specific validation
	return nil
}

// OnStart is called when an ARK server is started.
func (h *ArkHandler) OnStart(_ context.Context, _ string) error {
	// TODO: Implement ARK-specific startup logic
	return nil
}

// OnStop is called when an ARK server is stopped.
func (h *ArkHandler) OnStop(_ context.Context, _ string) error {
	// TODO: Implement ARK-specific shutdown logic (save world, backup)
	return nil
}

func init() {
	Register("ark", &ArkHandler{})
}
