package main

import (
	"fmt"
	"os"

	"github.com/sweetfish329/sabakan/backend/internal/config"
	"github.com/sweetfish329/sabakan/backend/internal/db"
	"github.com/sweetfish329/sabakan/backend/internal/logger"
	"github.com/sweetfish329/sabakan/backend/internal/server"
)

func main() {
	// Determine config path from environment or default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.toml"
	}

	// Load Configuration (fallback to defaults if not found)
	cfg, err := config.LoadSystemConfig(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = config.DefaultSystemConfig()
		} else {
			// Use fmt before logger is initialized
			fmt.Printf("Failed to load config: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize Logger
	logger.Init(cfg.Logging.Level, cfg.Logging.Format)

	if err != nil {
		logger.Info("Config file not found, using defaults")
	}

	// Initialize Database
	if err := db.Init(cfg.Database.Path); err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	// Initialize and Start Server
	s := server.New()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Starting server", "address", addr)
	s.Logger.Fatal(s.Start(addr))
}
