package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSystemConfig_Success(t *testing.T) {
	// Arrange: Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")
	content := `
[server]
host = "127.0.0.1"
port = 8080

[database]
path = "./test.db"

[logging]
level = "debug"
format = "json"
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	// Act
	cfg, err := LoadSystemConfig(configPath)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "./test.db", cfg.Database.Path)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
}

func TestLoadSystemConfig_FileNotFound(t *testing.T) {
	// Act
	_, err := LoadSystemConfig("nonexistent.toml")

	// Assert
	assert.Error(t, err)
}

func TestLoadSystemConfig_InvalidTOML(t *testing.T) {
	// Arrange: Create an invalid TOML file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.toml")
	content := `
[server
host = "broken
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	// Act
	_, err = LoadSystemConfig(configPath)

	// Assert
	assert.Error(t, err)
}

func TestDefaultSystemConfig(t *testing.T) {
	// Act
	cfg := DefaultSystemConfig()

	// Assert
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 1323, cfg.Server.Port)
	assert.Equal(t, "./sabakan.db", cfg.Database.Path)
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "text", cfg.Logging.Format)
}

func TestLoadGameConfig_Success(t *testing.T) {
	// Arrange: Create a temporary game config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "game.toml")
	content := `
[game]
id = "minecraft-1"
name = "My Server"
description = "Test server"

[container]
image = "itzg/minecraft-server:latest"
ports = ["25565:25565"]
volumes = ["./data:/data"]

[container.env]
EULA = "TRUE"
MEMORY = "4G"

[[mods]]
id = "worldedit"
name = "WorldEdit"
enabled = true
path = "./mods/worldedit.jar"
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	// Act
	cfg, err := LoadGameConfig(configPath)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "minecraft-1", cfg.Game.ID)
	assert.Equal(t, "My Server", cfg.Game.Name)
	assert.Equal(t, "itzg/minecraft-server:latest", cfg.Container.Image)
	assert.Len(t, cfg.Container.Ports, 1)
	assert.Equal(t, "TRUE", cfg.Container.Env["EULA"])
	assert.Len(t, cfg.Mods, 1)
	assert.Equal(t, "worldedit", cfg.Mods[0].ID)
	assert.True(t, cfg.Mods[0].Enabled)
}

func TestSaveGameConfig_Success(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "game.toml")
	cfg := &GameConfig{
		Game: GameInfo{
			ID:          "test-game",
			Name:        "Test Game",
			Description: "A test game",
		},
		Container: ContainerInfo{
			Image:   "test-image:latest",
			Ports:   []string{"8080:8080"},
			Volumes: []string{"./data:/data"},
			Env:     map[string]string{"KEY": "VALUE"},
		},
		Mods: []ModInfo{
			{ID: "mod1", Name: "Mod 1", Enabled: true, Path: "./mods/mod1.jar"},
		},
	}

	// Act
	err := SaveGameConfig(configPath, cfg)
	require.NoError(t, err)

	// Assert: Reload and verify
	loaded, err := LoadGameConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, cfg.Game.ID, loaded.Game.ID)
	assert.Equal(t, cfg.Container.Image, loaded.Container.Image)
	assert.Len(t, loaded.Mods, 1)
}
