package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_Register(t *testing.T) {
	e := echo.New()
	handler := NewAuthHandler(nil, nil, nil) // Will be mocked

	t.Run("should register user successfully", func(t *testing.T) {
		reqBody := `{"username":"newuser","email":"new@example.com","password":"securePass123!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		// For now, just check that the handler exists and can be called
		// Full implementation will come after this test
		assert.NoError(t, err)
	})

	t.Run("should reject registration with missing username", func(t *testing.T) {
		reqBody := `{"email":"new@example.com","password":"securePass123!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should reject registration with short password", func(t *testing.T) {
		reqBody := `{"username":"newuser","email":"new@example.com","password":"123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	e := echo.New()
	handler := NewAuthHandler(nil, nil, nil)

	t.Run("should login successfully with valid credentials", func(t *testing.T) {
		reqBody := `{"username":"testuser","password":"correctPassword"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
	})

	t.Run("should reject login with missing credentials", func(t *testing.T) {
		reqBody := `{}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthHandler_Refresh(t *testing.T) {
	e := echo.New()
	handler := NewAuthHandler(nil, nil, nil)

	t.Run("should refresh token successfully", func(t *testing.T) {
		reqBody := `{"refresh_token":"valid-refresh-token"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Refresh(c)

		assert.NoError(t, err)
	})

	t.Run("should reject refresh with missing token", func(t *testing.T) {
		reqBody := `{}`
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Refresh(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	e := echo.New()
	handler := NewAuthHandler(nil, nil, nil)

	t.Run("should logout successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Logout(c)

		assert.NoError(t, err)
	})
}

// Test request/response types
func TestAuthRequestTypes(t *testing.T) {
	t.Run("RegisterRequest should unmarshal correctly", func(t *testing.T) {
		jsonStr := `{"username":"test","email":"test@example.com","password":"password123"}`
		var req RegisterRequest
		err := json.Unmarshal([]byte(jsonStr), &req)

		require.NoError(t, err)
		assert.Equal(t, "test", req.Username)
		assert.Equal(t, "test@example.com", req.Email)
		assert.Equal(t, "password123", req.Password)
	})

	t.Run("LoginRequest should unmarshal correctly", func(t *testing.T) {
		jsonStr := `{"username":"test","password":"password123"}`
		var req LoginRequest
		err := json.Unmarshal([]byte(jsonStr), &req)

		require.NoError(t, err)
		assert.Equal(t, "test", req.Username)
		assert.Equal(t, "password123", req.Password)
	})
}
