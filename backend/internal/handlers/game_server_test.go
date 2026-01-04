package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/sweetfish329/sabakan/backend/internal/middleware"
	"github.com/sweetfish329/sabakan/backend/internal/models"
	"gorm.io/gorm"
)

// setupGameServerTestDB creates an in-memory SQLite database for testing.
func setupGameServerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Migrate models
	err = db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.User{},
		&models.GameServer{},
		&models.GameServerPort{},
		&models.GameServerEnv{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create test role and user
	role := models.Role{Name: "admin", DisplayName: "Administrator"}
	db.Create(&role)
	user := models.User{Username: "testuser", RoleID: role.ID, IsActive: true}
	db.Create(&user)

	return db
}

func TestGameServerHandler_Create(t *testing.T) {
	db := setupGameServerTestDB(t)
	handler := NewGameServerHandler(db)

	e := echo.New()

	t.Run("should create a game server", func(t *testing.T) {
		reqBody := CreateGameServerRequest{
			Slug:        "my-minecraft",
			Name:        "My Minecraft Server",
			Game:        "minecraft",
			Description: "Test server",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/game-servers", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(middleware.ContextKeyUserID, uint(1))

		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.GameServer
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "my-minecraft", resp.Slug)
		assert.Equal(t, "My Minecraft Server", resp.Name)
	})
}

func TestGameServerHandler_Create_InvalidSlug(t *testing.T) {
	db := setupGameServerTestDB(t)
	handler := NewGameServerHandler(db)

	e := echo.New()

	t.Run("should reject invalid slug with spaces", func(t *testing.T) {
		reqBody := CreateGameServerRequest{
			Slug: "my minecraft",
			Name: "My Server",
			Game: "minecraft",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/game-servers", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(middleware.ContextKeyUserID, uint(1))

		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should reject empty slug", func(t *testing.T) {
		reqBody := CreateGameServerRequest{
			Slug: "",
			Name: "My Server",
			Game: "minecraft",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/game-servers", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(middleware.ContextKeyUserID, uint(1))

		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGameServerHandler_Create_DuplicateSlug(t *testing.T) {
	db := setupGameServerTestDB(t)
	handler := NewGameServerHandler(db)

	// Create existing server
	db.Create(&models.GameServer{
		Slug:    "existing-server",
		Name:    "Existing Server",
		Image:   "test:latest",
		OwnerID: 1,
	})

	e := echo.New()

	t.Run("should reject duplicate slug", func(t *testing.T) {
		reqBody := CreateGameServerRequest{
			Slug: "existing-server",
			Name: "New Server",
			Game: "minecraft",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/game-servers", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(middleware.ContextKeyUserID, uint(1))

		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestGameServerHandler_List(t *testing.T) {
	db := setupGameServerTestDB(t)
	handler := NewGameServerHandler(db)

	// Create test servers
	db.Create(&models.GameServer{Slug: "server1", Name: "Server 1", Image: "test:latest", OwnerID: 1})
	db.Create(&models.GameServer{Slug: "server2", Name: "Server 2", Image: "test:latest", OwnerID: 1})
	db.Create(&models.GameServer{Slug: "other-user", Name: "Other User Server", Image: "test:latest", OwnerID: 2})

	e := echo.New()

	t.Run("should list user's game servers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/game-servers", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(middleware.ContextKeyUserID, uint(1))

		err := handler.List(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var servers []models.GameServer
		err = json.Unmarshal(rec.Body.Bytes(), &servers)
		assert.NoError(t, err)
		assert.Len(t, servers, 2)
	})
}

func TestGameServerHandler_Get(t *testing.T) {
	db := setupGameServerTestDB(t)
	handler := NewGameServerHandler(db)

	// Create test server
	db.Create(&models.GameServer{Slug: "my-server", Name: "My Server", Image: "test:latest", OwnerID: 1})

	e := echo.New()

	t.Run("should get a game server by slug", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/game-servers/my-server", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues("my-server")
		c.Set(middleware.ContextKeyUserID, uint(1))

		err := handler.Get(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var server models.GameServer
		err = json.Unmarshal(rec.Body.Bytes(), &server)
		assert.NoError(t, err)
		assert.Equal(t, "my-server", server.Slug)
	})
}

func TestGameServerHandler_Get_NotFound(t *testing.T) {
	db := setupGameServerTestDB(t)
	handler := NewGameServerHandler(db)

	e := echo.New()

	t.Run("should return 404 for nonexistent server", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/game-servers/nonexistent", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues("nonexistent")
		c.Set(middleware.ContextKeyUserID, uint(1))

		err := handler.Get(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestGameServerHandler_Delete(t *testing.T) {
	db := setupGameServerTestDB(t)
	handler := NewGameServerHandler(db)

	// Create test server
	db.Create(&models.GameServer{Slug: "to-delete", Name: "To Delete", Image: "test:latest", OwnerID: 1})

	e := echo.New()

	t.Run("should delete a game server", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/game-servers/to-delete", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues("to-delete")
		c.Set(middleware.ContextKeyUserID, uint(1))

		err := handler.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify deleted
		var count int64
		db.Model(&models.GameServer{}).Where("slug = ?", "to-delete").Count(&count)
		assert.Equal(t, int64(0), count)
	})
}
