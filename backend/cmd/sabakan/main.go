package main

import (
	"fmt"
	"os"

	"github.com/sweetfish329/sabakan/backend/internal/config"
	"github.com/sweetfish329/sabakan/backend/internal/container"
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

	// Run database migrations
	if err := db.Migrate(); err != nil {
		logger.Error("Failed to run database migrations", "error", err)
		os.Exit(1)
	}
	logger.Info("Database migrations completed")

	// Seed default data
	if err := db.Seed(); err != nil {
		logger.Error("Failed to seed database", "error", err)
		os.Exit(1)
	}
	logger.Info("Database seeding completed")

	// Initialize Container Service
	containerService := container.NewService(cfg.Podman.SocketPath)
	logger.Info("Container service initialized", "socket", cfg.Podman.SocketPath)

	// Create server dependencies
	deps := &server.Dependencies{
		ContainerService: containerService,
		DB:               db.GetDB(),
		Config:           cfg,
		SessionStore:     nil, // Redis session store (optional, can be nil for now)
	}

	// Initialize and Start Server
	s := server.New(deps)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Starting server", "address", addr)
	s.Logger.Fatal(s.Start(addr))
}
