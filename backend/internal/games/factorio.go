package games

import (
	"context"
)

// FactorioHandler handles Factorio-specific operations.
// Uses factoriotools/factorio image.
type FactorioHandler struct{}

// GetDefaultImage returns the default Factorio server image.
func (h *FactorioHandler) GetDefaultImage() string {
	return "factoriotools/factorio:stable"
}

// GetDefaultPorts returns the default Factorio port mappings.
func (h *FactorioHandler) GetDefaultPorts() map[string]int {
	return map[string]int{
		"game": 34197,
		"rcon": 27015,
	}
}

// GetDefaultEnv returns the default Factorio environment variables.
func (h *FactorioHandler) GetDefaultEnv() map[string]string {
	return map[string]string{
		"PORT":                 "34197",
		"RCON_PORT":            "27015",
		"UPDATE_MODS_ON_START": "false",
		"DLC_SPACE_AGE":        "false",
		"SAVE_NAME":            "world",
		"TOKEN":                "",
		"USERNAME":             "",
		"GENERATE_NEW_SAVE":    "true",
		"LOAD_LATEST_SAVE":     "true",
	}
}

// ValidateConfig validates Factorio-specific configuration.
func (h *FactorioHandler) ValidateConfig(config map[string]any) error {
	// TODO: Implement Factorio-specific validation
	return nil
}

// OnStart is called when a Factorio server is started.
func (h *FactorioHandler) OnStart(_ context.Context, _ string) error {
	// TODO: Implement Factorio-specific startup logic
	return nil
}

// OnStop is called when a Factorio server is stopped.
func (h *FactorioHandler) OnStop(_ context.Context, _ string) error {
	// TODO: Implement Factorio-specific shutdown logic (autosave)
	return nil
}

func init() {
	Register("factorio", &FactorioHandler{})
}
