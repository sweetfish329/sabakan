package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameServerEnv_GetVisibleValue(t *testing.T) {
	t.Run("should return actual value when not secret", func(t *testing.T) {
		env := &GameServerEnv{
			Key:      "SERVER_PORT",
			Value:    "25565",
			IsSecret: false,
		}
		assert.Equal(t, "25565", env.GetVisibleValue())
	})

	t.Run("should return masked value when secret", func(t *testing.T) {
		env := &GameServerEnv{
			Key:      "RCON_PASSWORD",
			Value:    "supersecretpassword",
			IsSecret: true,
		}
		assert.Equal(t, "********", env.GetVisibleValue())
	})

	t.Run("should return empty string when value is empty and not secret", func(t *testing.T) {
		env := &GameServerEnv{
			Key:      "EMPTY_VAR",
			Value:    "",
			IsSecret: false,
		}
		assert.Equal(t, "", env.GetVisibleValue())
	})

	t.Run("should return masked value when value is empty but marked as secret", func(t *testing.T) {
		env := &GameServerEnv{
			Key:      "EMPTY_SECRET",
			Value:    "",
			IsSecret: true,
		}
		assert.Equal(t, "********", env.GetVisibleValue())
	})
}

func TestGameServerStatus_Constants(t *testing.T) {
	t.Run("should have correct status constants", func(t *testing.T) {
		assert.Equal(t, GameServerStatus("stopped"), GameServerStatusStopped)
		assert.Equal(t, GameServerStatus("running"), GameServerStatusRunning)
		assert.Equal(t, GameServerStatus("creating"), GameServerStatusCreating)
		assert.Equal(t, GameServerStatus("error"), GameServerStatusError)
	})
}
