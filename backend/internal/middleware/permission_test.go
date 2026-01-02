package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/sweetfish329/sabakan/backend/internal/auth"
	"github.com/sweetfish329/sabakan/backend/internal/models"
	"gorm.io/gorm"
)

const permTestSecret = "test-secret-for-permission-32bytes!"

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Migrate models
	err = db.AutoMigrate(&models.Role{}, &models.Permission{}, &models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

// seedTestData creates test roles, permissions, and users.
func seedTestData(t *testing.T, db *gorm.DB) {
	// Create permissions
	permissions := []models.Permission{
		{Resource: "system", Action: "admin", Description: "Full access"},
		{Resource: "game_server", Action: "create", Description: "Create servers"},
		{Resource: "game_server", Action: "read", Description: "View servers"},
		{Resource: "game_server", Action: "delete", Description: "Delete servers"},
	}
	for _, p := range permissions {
		db.Create(&p)
	}

	// Create admin role with system:admin permission
	adminRole := models.Role{Name: "admin", DisplayName: "Administrator", Priority: 100, IsSystem: true}
	db.Create(&adminRole)
	var systemAdminPerm models.Permission
	db.Where("resource = ? AND action = ?", "system", "admin").First(&systemAdminPerm)
	db.Model(&adminRole).Association("Permissions").Append(&systemAdminPerm)

	// Create user role with limited permissions
	userRole := models.Role{Name: "user", DisplayName: "User", Priority: 10, IsSystem: true}
	db.Create(&userRole)
	var readPerm models.Permission
	db.Where("resource = ? AND action = ?", "game_server", "read").First(&readPerm)
	db.Model(&userRole).Association("Permissions").Append(&readPerm)

	// Create guest role with no permissions
	guestRole := models.Role{Name: "guest", DisplayName: "Guest", Priority: 0, IsSystem: true}
	db.Create(&guestRole)

	// Create test users
	db.Create(&models.User{Username: "admin", RoleID: adminRole.ID, IsActive: true})
	db.Create(&models.User{Username: "regularuser", RoleID: userRole.ID, IsActive: true})
	db.Create(&models.User{Username: "guest", RoleID: guestRole.ID, IsActive: true})
}

func TestPermissionMiddleware_Allowed(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)
	jwtManager := auth.NewJWTManager(permTestSecret, 15*time.Minute, 7*24*time.Hour)
	permMiddleware := NewPermissionMiddleware(db)

	e := echo.New()

	t.Run("should allow user with required permission", func(t *testing.T) {
		// Get regular user (has game_server:read permission)
		var user models.User
		db.Where("username = ?", "regularuser").First(&user)

		token, _, _ := jwtManager.GenerateAccessToken(user.ID, user.Username)

		handler := permMiddleware.RequirePermission("game_server", "read")(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(ContextKeyUserID, user.ID)
		c.Set(ContextKeyUsername, user.Username)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestPermissionMiddleware_Denied(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)
	jwtManager := auth.NewJWTManager(permTestSecret, 15*time.Minute, 7*24*time.Hour)
	permMiddleware := NewPermissionMiddleware(db)

	e := echo.New()

	t.Run("should deny user without required permission", func(t *testing.T) {
		// Get regular user (has game_server:read but NOT game_server:delete)
		var user models.User
		db.Where("username = ?", "regularuser").First(&user)

		token, _, _ := jwtManager.GenerateAccessToken(user.ID, user.Username)

		handler := permMiddleware.RequirePermission("game_server", "delete")(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(ContextKeyUserID, user.ID)
		c.Set(ContextKeyUsername, user.Username)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("should deny guest with no permissions", func(t *testing.T) {
		var user models.User
		db.Where("username = ?", "guest").First(&user)

		token, _, _ := jwtManager.GenerateAccessToken(user.ID, user.Username)

		handler := permMiddleware.RequirePermission("game_server", "read")(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(ContextKeyUserID, user.ID)
		c.Set(ContextKeyUsername, user.Username)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestPermissionMiddleware_AdminBypass(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)
	jwtManager := auth.NewJWTManager(permTestSecret, 15*time.Minute, 7*24*time.Hour)
	permMiddleware := NewPermissionMiddleware(db)

	e := echo.New()

	t.Run("should allow admin to access any resource", func(t *testing.T) {
		var user models.User
		db.Where("username = ?", "admin").First(&user)

		token, _, _ := jwtManager.GenerateAccessToken(user.ID, user.Username)

		// Admin should access game_server:delete even without explicit permission
		handler := permMiddleware.RequirePermission("game_server", "delete")(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(ContextKeyUserID, user.ID)
		c.Set(ContextKeyUsername, user.Username)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("admin should access any arbitrary resource", func(t *testing.T) {
		var user models.User
		db.Where("username = ?", "admin").First(&user)

		token, _, _ := jwtManager.GenerateAccessToken(user.ID, user.Username)

		handler := permMiddleware.RequirePermission("nonexistent", "action")(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(ContextKeyUserID, user.ID)
		c.Set(ContextKeyUsername, user.Username)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestPermissionMiddleware_MissingUserContext(t *testing.T) {
	db := setupTestDB(t)
	permMiddleware := NewPermissionMiddleware(db)

	e := echo.New()

	t.Run("should deny when user ID is not in context", func(t *testing.T) {
		handler := permMiddleware.RequirePermission("game_server", "read")(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// No user ID set in context

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

// Helper to suppress unused variable warnings
var _ = context.Background
