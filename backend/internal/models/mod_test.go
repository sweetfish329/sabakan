package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGameServerMod_BeforeCreate(t *testing.T) {
	t.Run("should set InstalledAt when zero", func(t *testing.T) {
		gsm := &GameServerMod{}
		assert.True(t, gsm.InstalledAt.IsZero())

		err := gsm.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.False(t, gsm.InstalledAt.IsZero())
		assert.WithinDuration(t, time.Now(), gsm.InstalledAt, time.Second)
	})

	t.Run("should not override InstalledAt when already set", func(t *testing.T) {
		customTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
		gsm := &GameServerMod{
			InstalledAt: customTime,
		}

		err := gsm.BeforeCreate(&gorm.DB{})
		assert.NoError(t, err)
		assert.Equal(t, customTime, gsm.InstalledAt)
	})
}
