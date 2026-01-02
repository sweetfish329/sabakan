package models

import (
	"time"
)

// AuditLogAction represents the type of action performed.
type AuditLogAction string

const (
	// AuditLogActionCreate indicates a create action.
	AuditLogActionCreate AuditLogAction = "create"
	// AuditLogActionUpdate indicates an update action.
	AuditLogActionUpdate AuditLogAction = "update"
	// AuditLogActionDelete indicates a delete action.
	AuditLogActionDelete AuditLogAction = "delete"
	// AuditLogActionStart indicates a start action (for game servers).
	AuditLogActionStart AuditLogAction = "start"
	// AuditLogActionStop indicates a stop action (for game servers).
	AuditLogActionStop AuditLogAction = "stop"
	// AuditLogActionLogin indicates a user login.
	AuditLogActionLogin AuditLogAction = "login"
	// AuditLogActionLogout indicates a user logout.
	AuditLogActionLogout AuditLogAction = "logout"
)

// AuditLogTargetType represents the type of target entity.
type AuditLogTargetType string

const (
	// AuditLogTargetUser indicates a user target.
	AuditLogTargetUser AuditLogTargetType = "user"
	// AuditLogTargetGameServer indicates a game server target.
	AuditLogTargetGameServer AuditLogTargetType = "game_server"
	// AuditLogTargetMod indicates a mod target.
	AuditLogTargetMod AuditLogTargetType = "mod"
	// AuditLogTargetRole indicates a role target.
	AuditLogTargetRole AuditLogTargetType = "role"
	// AuditLogTargetSession indicates a session target.
	AuditLogTargetSession AuditLogTargetType = "session"
)

// AuditLog records user actions for auditing purposes.
type AuditLog struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time          `json:"createdAt"`
	UserID      *uint              `gorm:"index" json:"userId,omitempty"`
	User        *User              `json:"user,omitempty"`
	TargetType  AuditLogTargetType `gorm:"not null;index:idx_audit_target" json:"targetType"`
	TargetID    uint               `gorm:"index:idx_audit_target" json:"targetId"`
	Action      AuditLogAction     `gorm:"not null" json:"action"`
	DetailsJSON string             `json:"details,omitempty"`
	IPAddress   string             `json:"ipAddress,omitempty"`
}

// TableName specifies the table name for AuditLog.
func (AuditLog) TableName() string {
	return "audit_logs"
}
