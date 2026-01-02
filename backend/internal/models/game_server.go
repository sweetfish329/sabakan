package models

import (
	"gorm.io/gorm"
)

// GameServerStatus represents the current status of a game server.
type GameServerStatus string

const (
	// GameServerStatusStopped indicates a stopped server.
	GameServerStatusStopped GameServerStatus = "stopped"
	// GameServerStatusRunning indicates a running server.
	GameServerStatusRunning GameServerStatus = "running"
	// GameServerStatusCreating indicates a server being created.
	GameServerStatusCreating GameServerStatus = "creating"
	// GameServerStatusError indicates a server in error state.
	GameServerStatusError GameServerStatus = "error"
)

// GameServer represents a managed game server container.
type GameServer struct {
	gorm.Model
	Slug        string           `gorm:"uniqueIndex;not null" json:"slug"`
	Name        string           `gorm:"not null" json:"name"`
	Description string           `json:"description,omitempty"`
	Image       string           `gorm:"not null" json:"image"`
	Status      GameServerStatus `gorm:"default:stopped" json:"status"`
	ContainerID string           `json:"containerId,omitempty"`
	OwnerID     uint             `gorm:"index" json:"ownerId"`
	Owner       User             `json:"owner,omitempty"`
	Ports       []GameServerPort `json:"ports,omitempty"`
	Envs        []GameServerEnv  `json:"envs,omitempty"`
	Mods        []GameServerMod  `json:"mods,omitempty"`
}

// GameServerPort represents a port mapping for a game server.
type GameServerPort struct {
	gorm.Model
	GameServerID  uint       `gorm:"not null;index" json:"gameServerId"`
	GameServer    GameServer `json:"-"`
	HostPort      int        `gorm:"not null" json:"hostPort"`
	ContainerPort int        `gorm:"not null" json:"containerPort"`
	Protocol      string     `gorm:"default:tcp" json:"protocol"`
}

// GameServerEnv represents an environment variable for a game server.
type GameServerEnv struct {
	gorm.Model
	GameServerID uint       `gorm:"not null;index" json:"gameServerId"`
	GameServer   GameServer `json:"-"`
	Key          string     `gorm:"not null" json:"key"`
	Value        string     `json:"value,omitempty"`
	IsSecret     bool       `gorm:"default:false" json:"isSecret"`
}

// GetVisibleValue returns the value or a masked string if it's a secret.
func (e *GameServerEnv) GetVisibleValue() string {
	if e.IsSecret {
		return "********"
	}
	return e.Value
}
