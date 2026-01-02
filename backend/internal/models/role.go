// Package models provides data structures for the Sabakan application.
package models

import (
	"gorm.io/gorm"
)

// Role represents a user role with permissions.
type Role struct {
	gorm.Model
	Name        string       `gorm:"uniqueIndex;not null" json:"name"`
	DisplayName string       `gorm:"not null" json:"displayName"`
	Description string       `json:"description,omitempty"`
	Priority    int          `gorm:"default:0" json:"priority"`
	IsSystem    bool         `gorm:"default:false" json:"isSystem"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `json:"-"`
}

// Permission represents an action that can be performed on a resource.
type Permission struct {
	gorm.Model
	Resource    string `gorm:"not null;uniqueIndex:idx_resource_action" json:"resource"`
	Action      string `gorm:"not null;uniqueIndex:idx_resource_action" json:"action"`
	Description string `json:"description,omitempty"`
	Roles       []Role `gorm:"many2many:role_permissions;" json:"-"`
}

// DefaultRoles returns the default system roles to be seeded.
func DefaultRoles() []Role {
	return []Role{
		{Name: "admin", DisplayName: "Administrator", Priority: 100, IsSystem: true},
		{Name: "moderator", DisplayName: "Moderator", Priority: 50, IsSystem: true},
		{Name: "user", DisplayName: "User", Priority: 10, IsSystem: true},
		{Name: "guest", DisplayName: "Guest", Priority: 0, IsSystem: true},
	}
}

// DefaultPermissions returns the default permissions to be seeded.
func DefaultPermissions() []Permission {
	return []Permission{
		{Resource: "system", Action: "admin", Description: "Full system access"},
		{Resource: "game_server", Action: "create", Description: "Create servers"},
		{Resource: "game_server", Action: "read", Description: "View servers"},
		{Resource: "game_server", Action: "update", Description: "Edit servers"},
		{Resource: "game_server", Action: "delete", Description: "Delete servers"},
		{Resource: "game_server", Action: "start", Description: "Start servers"},
		{Resource: "game_server", Action: "stop", Description: "Stop servers"},
		{Resource: "mod", Action: "create", Description: "Add mods"},
		{Resource: "mod", Action: "read", Description: "View mods"},
		{Resource: "mod", Action: "update", Description: "Edit mods"},
		{Resource: "mod", Action: "delete", Description: "Delete mods"},
		{Resource: "user", Action: "create", Description: "Create users"},
		{Resource: "user", Action: "read", Description: "View users"},
		{Resource: "user", Action: "update", Description: "Edit users"},
		{Resource: "user", Action: "delete", Description: "Delete users"},
		{Resource: "role", Action: "manage", Description: "Manage roles"},
		{Resource: "audit_log", Action: "read", Description: "View audit logs"},
	}
}
