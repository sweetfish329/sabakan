package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sweetfish329/sabakan/backend/internal/auth"
	"github.com/sweetfish329/sabakan/backend/internal/models"
	"gorm.io/gorm"
)

const testJWTSecret = "test-secret-key-for-auth-handler-32b!"

// setupAuthTestDB creates an in-memory SQLite database for testing.
func setupAuthTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.User{},
		&models.RefreshToken{},
	)
	require.NoError(t, err)

	// Seed default role
	userRole := models.Role{Name: "user", DisplayName: "User", Priority: 10, IsSystem: true}
	db.Create(&userRole)

	return db
}

func TestAuthHandler_Register_Integration(t *testing.T) {
	db := setupAuthTestDB(t)
	jwtManager := auth.NewJWTManager(testJWTSecret, 15*time.Minute, 7*24*time.Hour)
	handler := NewAuthHandler(db, jwtManager, nil)
	e := echo.New()

	t.Run("should register new user and create in database", func(t *testing.T) {
		reqBody := `{"username":"newuser","email":"new@example.com","password":"securePass123!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Verify user was created in database
		var user models.User
		result := db.Where("username = ?", "newuser").First(&user)
		assert.NoError(t, result.Error)
		assert.Equal(t, "newuser", user.Username)
		assert.NotEmpty(t, user.PasswordHash)
		assert.True(t, user.IsActive)
	})

	t.Run("should reject duplicate username", func(t *testing.T) {
		// First registration
		reqBody := `{"username":"duplicate","email":"dup1@example.com","password":"securePass123!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = handler.Register(c)

		// Second registration with same username
		reqBody2 := `{"username":"duplicate","email":"dup2@example.com","password":"securePass123!"}`
		req2 := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody2))
		req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec2 := httptest.NewRecorder()
		c2 := e.NewContext(req2, rec2)

		err := handler.Register(c2)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec2.Code)
	})
}

func TestAuthHandler_Login_Integration(t *testing.T) {
	db := setupAuthTestDB(t)
	jwtManager := auth.NewJWTManager(testJWTSecret, 15*time.Minute, 7*24*time.Hour)
	handler := NewAuthHandler(db, jwtManager, nil)
	e := echo.New()

	// First, register a user
	registerBody := `{"username":"logintest","email":"login@example.com","password":"correctPassword123!"}`
	registerReq := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerBody))
	registerReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()
	registerCtx := e.NewContext(registerReq, registerRec)
	_ = handler.Register(registerCtx)

	t.Run("should login with correct credentials and return JWT", func(t *testing.T) {
		reqBody := `{"username":"logintest","password":"correctPassword123!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response AuthResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.Equal(t, "Bearer", response.TokenType)
		assert.Greater(t, response.ExpiresIn, 0)

		// Verify access token is valid
		claims, err := jwtManager.ValidateAccessToken(response.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, "logintest", claims.Username)
	})

	t.Run("should reject login with wrong password", func(t *testing.T) {
		reqBody := `{"username":"logintest","password":"wrongPassword!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("should reject login with non-existent user", func(t *testing.T) {
		reqBody := `{"username":"nonexistent","password":"anyPassword123!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAuthHandler_Refresh_Integration(t *testing.T) {
	db := setupAuthTestDB(t)
	jwtManager := auth.NewJWTManager(testJWTSecret, 15*time.Minute, 7*24*time.Hour)
	handler := NewAuthHandler(db, jwtManager, nil)
	e := echo.New()

	// Register and login to get tokens
	registerBody := `{"username":"refreshtest","email":"refresh@example.com","password":"testPassword123!"}`
	registerReq := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerBody))
	registerReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()
	_ = handler.Register(e.NewContext(registerReq, registerRec))

	loginBody := `{"username":"refreshtest","password":"testPassword123!"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginBody))
	loginReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	loginRec := httptest.NewRecorder()
	_ = handler.Login(e.NewContext(loginReq, loginRec))

	var loginResponse AuthResponse
	_ = json.Unmarshal(loginRec.Body.Bytes(), &loginResponse)

	t.Run("should refresh tokens with valid refresh token", func(t *testing.T) {
		reqBody := `{"refresh_token":"` + loginResponse.RefreshToken + `"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Refresh(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response AuthResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEqual(t, loginResponse.AccessToken, response.AccessToken)
	})

	t.Run("should reject invalid refresh token", func(t *testing.T) {
		reqBody := `{"refresh_token":"invalid-token"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Refresh(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
