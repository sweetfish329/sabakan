package games

import (
	"context"
)

// SevenDaysToDieHandler handles 7 Days to Die-specific operations.
// Uses vinanrra/7dtd-server image.
type SevenDaysToDieHandler struct{}

// GetDefaultImage returns the default 7 Days to Die server image.
func (h *SevenDaysToDieHandler) GetDefaultImage() string {
	return "vinanrra/7dtd-server:latest"
}

// GetDefaultPorts returns the default 7 Days to Die port mappings.
func (h *SevenDaysToDieHandler) GetDefaultPorts() map[string]int {
	return map[string]int{
		"game":   26900,
		"telnet": 8081,
		"web":    8082,
	}
}

// GetDefaultEnv returns the default 7 Days to Die environment variables.
func (h *SevenDaysToDieHandler) GetDefaultEnv() map[string]string {
	return map[string]string{
		"START_MODE":         "1",
		"VERSION":            "stable",
		"SERVER_NAME":        "Sabakan 7DTD Server",
		"SERVER_PASSWORD":    "",
		"SERVER_PORT":        "26900",
		"TELNET_PORT":        "8081",
		"WEB_PORT":           "8082",
		"MAX_PLAYERS":        "8",
		"TELNET_PASSWORD":    "changeme",
		"UPDATE_ON_START":    "NO",
		"BACKUP_ON_SHUTDOWN": "YES",
	}
}

// ValidateConfig validates 7 Days to Die-specific configuration.
func (h *SevenDaysToDieHandler) ValidateConfig(config map[string]any) error {
	// TODO: Implement 7DTD-specific validation
	return nil
}

// OnStart is called when a 7 Days to Die server is started.
func (h *SevenDaysToDieHandler) OnStart(_ context.Context, _ string) error {
	// TODO: Implement 7DTD-specific startup logic
	return nil
}

// OnStop is called when a 7 Days to Die server is stopped.
func (h *SevenDaysToDieHandler) OnStop(_ context.Context, _ string) error {
	// TODO: Implement 7DTD-specific shutdown logic (backup world)
	return nil
}

func init() {
	Register("7daystodie", &SevenDaysToDieHandler{})
}
