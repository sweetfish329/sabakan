package games

import (
	"context"
)

// SatisfactoryHandler handles Satisfactory-specific operations.
// Uses wolveix/satisfactory-server image.
type SatisfactoryHandler struct{}

// GetDefaultImage returns the default Satisfactory server image.
func (h *SatisfactoryHandler) GetDefaultImage() string {
	return "wolveix/satisfactory-server:latest"
}

// GetDefaultPorts returns the default Satisfactory port mappings.
func (h *SatisfactoryHandler) GetDefaultPorts() map[string]int {
	return map[string]int{
		"game":   7777,
		"query":  15000,
		"beacon": 15777,
	}
}

// GetDefaultEnv returns the default Satisfactory environment variables.
func (h *SatisfactoryHandler) GetDefaultEnv() map[string]string {
	return map[string]string{
		"MAXPLAYERS":            "4",
		"PGID":                  "1000",
		"PUID":                  "1000",
		"ROOTLESS":              "false",
		"STEAMBETA":             "false",
		"AUTOPAUSE":             "true",
		"AUTOSAVEINTERVAL":      "300",
		"AUTOSAVENUM":           "3",
		"AUTOSAVEONDISCONNECT":  "true",
		"CRASHREPORT":           "true",
		"DEBUG":                 "false",
		"DISABLESEASONALEVENTS": "false",
		"NETWORKQUALITY":        "3",
	}
}

// ValidateConfig validates Satisfactory-specific configuration.
func (h *SatisfactoryHandler) ValidateConfig(config map[string]any) error {
	// TODO: Implement Satisfactory-specific validation
	return nil
}

// OnStart is called when a Satisfactory server is started.
func (h *SatisfactoryHandler) OnStart(_ context.Context, _ string) error {
	// TODO: Implement Satisfactory-specific startup logic
	return nil
}

// OnStop is called when a Satisfactory server is stopped.
func (h *SatisfactoryHandler) OnStop(_ context.Context, _ string) error {
	// TODO: Implement Satisfactory-specific shutdown logic (autosave)
	return nil
}

func init() {
	Register("satisfactory", &SatisfactoryHandler{})
}
