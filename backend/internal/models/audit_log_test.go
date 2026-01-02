package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuditLogAction_Constants(t *testing.T) {
	t.Run("should have correct action constants", func(t *testing.T) {
		assert.Equal(t, AuditLogAction("create"), AuditLogActionCreate)
		assert.Equal(t, AuditLogAction("update"), AuditLogActionUpdate)
		assert.Equal(t, AuditLogAction("delete"), AuditLogActionDelete)
		assert.Equal(t, AuditLogAction("start"), AuditLogActionStart)
		assert.Equal(t, AuditLogAction("stop"), AuditLogActionStop)
		assert.Equal(t, AuditLogAction("login"), AuditLogActionLogin)
		assert.Equal(t, AuditLogAction("logout"), AuditLogActionLogout)
	})
}

func TestAuditLogTargetType_Constants(t *testing.T) {
	t.Run("should have correct target type constants", func(t *testing.T) {
		assert.Equal(t, AuditLogTargetType("user"), AuditLogTargetUser)
		assert.Equal(t, AuditLogTargetType("game_server"), AuditLogTargetGameServer)
		assert.Equal(t, AuditLogTargetType("mod"), AuditLogTargetMod)
		assert.Equal(t, AuditLogTargetType("role"), AuditLogTargetRole)
		assert.Equal(t, AuditLogTargetType("session"), AuditLogTargetSession)
	})
}

func TestAuditLog_TableName(t *testing.T) {
	t.Run("should return correct table name", func(t *testing.T) {
		log := AuditLog{}
		assert.Equal(t, "audit_logs", log.TableName())
	})
}
