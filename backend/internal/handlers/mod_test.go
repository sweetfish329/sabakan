package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sweetfish329/sabakan/backend/internal/models"
	"gorm.io/gorm"
)

// setupModTestDB creates an in-memory SQLite database for testing.
func setupModTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to create test database")

	err = db.AutoMigrate(&models.Mod{})
	require.NoError(t, err, "Failed to migrate")

	return db
}

// seedModTestData creates test mods.
func seedModTestData(t *testing.T, db *gorm.DB) {
	mods := []models.Mod{
		{Name: "Test Mod 1", Slug: "test-mod-1", Description: "First test mod", Version: "1.0.0"},
		{Name: "Test Mod 2", Slug: "test-mod-2", Description: "Second test mod", Version: "2.0.0"},
	}
	for _, m := range mods {
		require.NoError(t, db.Create(&m).Error)
	}
}

func TestModHandler_List(t *testing.T) {
	db := setupModTestDB(t)
	seedModTestData(t, db)
	handler := NewModHandler(db)
	e := echo.New()

	t.Run("should return list of mods", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/mods", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.List(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var mods []models.Mod
		err = json.Unmarshal(rec.Body.Bytes(), &mods)
		assert.NoError(t, err)
		assert.Len(t, mods, 2)
	})
}

func TestModHandler_Get(t *testing.T) {
	db := setupModTestDB(t)
	seedModTestData(t, db)
	handler := NewModHandler(db)
	e := echo.New()

	t.Run("should return mod by ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/mods/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.Get(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var mod models.Mod
		err = json.Unmarshal(rec.Body.Bytes(), &mod)
		assert.NoError(t, err)
		assert.Equal(t, "Test Mod 1", mod.Name)
	})

	t.Run("should return 404 for non-existent mod", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/mods/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := handler.Get(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	})
}

func TestModHandler_Create(t *testing.T) {
	db := setupModTestDB(t)
	handler := NewModHandler(db)
	e := echo.New()

	t.Run("should create a new mod", func(t *testing.T) {
		body := `{"name":"New Mod","slug":"new-mod","description":"A new mod","version":"1.0.0"}`
		req := httptest.NewRequest(http.MethodPost, "/api/mods", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var mod models.Mod
		err = json.Unmarshal(rec.Body.Bytes(), &mod)
		assert.NoError(t, err)
		assert.Equal(t, "New Mod", mod.Name)
		assert.NotZero(t, mod.ID)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		body := `{invalid json}`
		req := httptest.NewRequest(http.MethodPost, "/api/mods", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("should return 400 for missing required fields", func(t *testing.T) {
		body := `{"description":"Missing name and slug"}`
		req := httptest.NewRequest(http.MethodPost, "/api/mods", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	})
}

func TestModHandler_Update(t *testing.T) {
	db := setupModTestDB(t)
	seedModTestData(t, db)
	handler := NewModHandler(db)
	e := echo.New()

	t.Run("should update an existing mod", func(t *testing.T) {
		body := `{"name":"Updated Mod","description":"Updated description"}`
		req := httptest.NewRequest(http.MethodPut, "/api/mods/1", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var mod models.Mod
		err = json.Unmarshal(rec.Body.Bytes(), &mod)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Mod", mod.Name)
		assert.Equal(t, "Updated description", mod.Description)
	})

	t.Run("should return 404 for non-existent mod", func(t *testing.T) {
		body := `{"name":"Updated Mod"}`
		req := httptest.NewRequest(http.MethodPut, "/api/mods/999", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := handler.Update(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	})
}

func TestModHandler_Delete(t *testing.T) {
	db := setupModTestDB(t)
	seedModTestData(t, db)
	handler := NewModHandler(db)
	e := echo.New()

	t.Run("should delete an existing mod", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/mods/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify deletion
		var mod models.Mod
		err = db.First(&mod, 1).Error
		assert.Error(t, err) // Should not find
	})

	t.Run("should return 404 for non-existent mod", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/mods/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := handler.Delete(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	})
}
