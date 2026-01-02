package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/sweetfish329/sabakan/backend/internal/auth"
)

const testSecret = "test-secret-key-for-middleware-32b!"

func TestAuthMiddleware_ValidToken(t *testing.T) {
	e := echo.New()
	jwtManager := auth.NewJWTManager(testSecret, 15*time.Minute, 7*24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager, nil)

	handler := middleware.Authenticate(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	t.Run("should allow request with valid token", func(t *testing.T) {
		token, _, _ := jwtManager.GenerateAccessToken(123, "testuser")

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	e := echo.New()
	jwtManager := auth.NewJWTManager(testSecret, 15*time.Minute, 7*24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager, nil)

	handler := middleware.Authenticate(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	t.Run("should reject request without Authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("should reject request with empty Authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("should reject request with missing Bearer prefix", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "some-token-without-bearer")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	e := echo.New()
	jwtManager := auth.NewJWTManager(testSecret, 15*time.Minute, 7*24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager, nil)

	handler := middleware.Authenticate(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	t.Run("should reject request with invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("should reject request with token signed by different secret", func(t *testing.T) {
		otherJWT := auth.NewJWTManager("other-secret-key-32bytes!!", 15*time.Minute, 7*24*time.Hour)
		token, _, _ := otherJWT.GenerateAccessToken(123, "testuser")

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	e := echo.New()
	jwtManager := auth.NewJWTManager(testSecret, 1*time.Millisecond, 7*24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager, nil)

	handler := middleware.Authenticate(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	t.Run("should reject expired token", func(t *testing.T) {
		token, _, _ := jwtManager.GenerateAccessToken(123, "testuser")
		time.Sleep(10 * time.Millisecond)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAuthMiddleware_ContextValues(t *testing.T) {
	e := echo.New()
	jwtManager := auth.NewJWTManager(testSecret, 15*time.Minute, 7*24*time.Hour)
	middleware := NewAuthMiddleware(jwtManager, nil)

	t.Run("should set user claims in context", func(t *testing.T) {
		var capturedUserID uint
		var capturedUsername string

		handler := middleware.Authenticate(func(c echo.Context) error {
			capturedUserID = GetUserID(c)
			capturedUsername = GetUsername(c)
			return c.String(http.StatusOK, "OK")
		})

		token, _, _ := jwtManager.GenerateAccessToken(456, "contextuser")

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, uint(456), capturedUserID)
		assert.Equal(t, "contextuser", capturedUsername)
	})
}
