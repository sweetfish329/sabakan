package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sweetfish329/sabakan/backend/internal/models"
	"gorm.io/gorm"
)

// PermissionMiddleware handles permission-based authorization.
type PermissionMiddleware struct {
	db *gorm.DB
}

// NewPermissionMiddleware creates a new permission middleware.
func NewPermissionMiddleware(db *gorm.DB) *PermissionMiddleware {
	return &PermissionMiddleware{db: db}
}

// RequirePermission returns a middleware that checks if the user has the required permission.
func (m *PermissionMiddleware) RequirePermission(resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user ID from context (set by auth middleware)
			userID := GetUserID(c)
			if userID == 0 {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "unauthorized",
					"message": "User not authenticated",
				})
			}

			// Check if user has the required permission
			hasPermission, err := m.checkPermission(userID, resource, action)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error":   "internal_error",
					"message": "Failed to check permissions",
				})
			}

			if !hasPermission {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error":   "forbidden",
					"message": "You do not have permission to perform this action",
				})
			}

			return next(c)
		}
	}
}

// checkPermission checks if a user has the required permission.
func (m *PermissionMiddleware) checkPermission(userID uint, resource, action string) (bool, error) {
	// Get user with role and permissions
	var user models.User
	if err := m.db.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return false, err
	}

	// Check for system:admin permission (admin bypass)
	for _, perm := range user.Role.Permissions {
		if perm.Resource == "system" && perm.Action == "admin" {
			return true, nil
		}
	}

	// Check for the specific permission
	for _, perm := range user.Role.Permissions {
		if perm.Resource == resource && perm.Action == action {
			return true, nil
		}
	}

	return false, nil
}

// RequireRole returns a middleware that checks if the user has the required role.
func (m *PermissionMiddleware) RequireRole(roleName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := GetUserID(c)
			if userID == 0 {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "unauthorized",
					"message": "User not authenticated",
				})
			}

			// Get user with role
			var user models.User
			if err := m.db.Preload("Role").First(&user, userID).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error":   "internal_error",
					"message": "Failed to get user",
				})
			}

			if user.Role.Name != roleName {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error":   "forbidden",
					"message": "You do not have the required role",
				})
			}

			return next(c)
		}
	}
}

// RequireAdmin is a convenience method for requiring admin role.
func (m *PermissionMiddleware) RequireAdmin() echo.MiddlewareFunc {
	return m.RequirePermission("system", "admin")
}
