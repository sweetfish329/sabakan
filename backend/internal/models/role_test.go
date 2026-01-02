package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRoles(t *testing.T) {
	roles := DefaultRoles()

	t.Run("should return 4 default roles", func(t *testing.T) {
		assert.Len(t, roles, 4)
	})

	t.Run("should have admin role with highest priority", func(t *testing.T) {
		var admin *Role
		for i := range roles {
			if roles[i].Name == "admin" {
				admin = &roles[i]
				break
			}
		}
		assert.NotNil(t, admin)
		assert.Equal(t, "admin", admin.Name)
		assert.Equal(t, "Administrator", admin.DisplayName)
		assert.Equal(t, 100, admin.Priority)
		assert.True(t, admin.IsSystem)
	})

	t.Run("should have moderator role", func(t *testing.T) {
		var moderator *Role
		for i := range roles {
			if roles[i].Name == "moderator" {
				moderator = &roles[i]
				break
			}
		}
		assert.NotNil(t, moderator)
		assert.Equal(t, 50, moderator.Priority)
		assert.True(t, moderator.IsSystem)
	})

	t.Run("should have user role", func(t *testing.T) {
		var user *Role
		for i := range roles {
			if roles[i].Name == "user" {
				user = &roles[i]
				break
			}
		}
		assert.NotNil(t, user)
		assert.Equal(t, 10, user.Priority)
		assert.True(t, user.IsSystem)
	})

	t.Run("should have guest role with lowest priority", func(t *testing.T) {
		var guest *Role
		for i := range roles {
			if roles[i].Name == "guest" {
				guest = &roles[i]
				break
			}
		}
		assert.NotNil(t, guest)
		assert.Equal(t, 0, guest.Priority)
		assert.True(t, guest.IsSystem)
	})

	t.Run("all roles should be system roles", func(t *testing.T) {
		for _, role := range roles {
			assert.True(t, role.IsSystem, "Role %s should be a system role", role.Name)
		}
	})
}

func TestDefaultPermissions(t *testing.T) {
	permissions := DefaultPermissions()

	t.Run("should return expected number of permissions", func(t *testing.T) {
		assert.GreaterOrEqual(t, len(permissions), 17)
	})

	t.Run("should have system:admin permission", func(t *testing.T) {
		found := false
		for _, p := range permissions {
			if p.Resource == "system" && p.Action == "admin" {
				found = true
				assert.Equal(t, "Full system access", p.Description)
				break
			}
		}
		assert.True(t, found, "system:admin permission should exist")
	})

	t.Run("should have all game_server permissions", func(t *testing.T) {
		expectedActions := []string{"create", "read", "update", "delete", "start", "stop"}
		for _, action := range expectedActions {
			found := false
			for _, p := range permissions {
				if p.Resource == "game_server" && p.Action == action {
					found = true
					break
				}
			}
			assert.True(t, found, "game_server:%s permission should exist", action)
		}
	})

	t.Run("should have all mod permissions", func(t *testing.T) {
		expectedActions := []string{"create", "read", "update", "delete"}
		for _, action := range expectedActions {
			found := false
			for _, p := range permissions {
				if p.Resource == "mod" && p.Action == action {
					found = true
					break
				}
			}
			assert.True(t, found, "mod:%s permission should exist", action)
		}
	})

	t.Run("should have all user permissions", func(t *testing.T) {
		expectedActions := []string{"create", "read", "update", "delete"}
		for _, action := range expectedActions {
			found := false
			for _, p := range permissions {
				if p.Resource == "user" && p.Action == action {
					found = true
					break
				}
			}
			assert.True(t, found, "user:%s permission should exist", action)
		}
	})

	t.Run("all permissions should have descriptions", func(t *testing.T) {
		for _, p := range permissions {
			assert.NotEmpty(t, p.Description, "Permission %s:%s should have a description", p.Resource, p.Action)
		}
	})
}
