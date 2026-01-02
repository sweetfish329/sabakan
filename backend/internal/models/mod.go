package models

import (
	"time"

	"gorm.io/gorm"
)

// Mod represents a mod in the catalog.
type Mod struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null" json:"slug"`
	Description string `json:"description,omitempty"`
	SourceURL   string `json:"sourceUrl,omitempty"`
	Version     string `json:"version,omitempty"`
}

// GameServerMod represents a mod installed on a game server.
type GameServerMod struct {
	gorm.Model
	GameServerID uint       `gorm:"not null;uniqueIndex:idx_server_mod" json:"gameServerId"`
	GameServer   GameServer `json:"-"`
	ModID        uint       `gorm:"not null;uniqueIndex:idx_server_mod" json:"modId"`
	Mod          Mod        `json:"mod,omitempty"`
	Enabled      bool       `gorm:"default:true" json:"enabled"`
	ConfigJSON   string     `json:"configJson,omitempty"`
	InstalledAt  time.Time  `json:"installedAt"`
}

// BeforeCreate sets the installed timestamp before creating.
func (gsm *GameServerMod) BeforeCreate(_ *gorm.DB) error {
	if gsm.InstalledAt.IsZero() {
		gsm.InstalledAt = time.Now()
	}
	return nil
}
