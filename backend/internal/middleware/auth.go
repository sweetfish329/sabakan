// Package middleware provides HTTP middleware for the Sabakan API.
package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sweetfish329/sabakan/backend/internal/auth"
	"github.com/sweetfish329/sabakan/backend/internal/redis"
)

// Context keys for user information.
const (
	ContextKeyUserID   = "user_id"
	ContextKeyUsername = "username"
	ContextKeyJTI      = "jti"
	ContextKeyClaims   = "claims"
)

// AuthMiddleware handles JWT authentication.
type AuthMiddleware struct {
	jwtManager   *auth.JWTManager
	sessionStore redis.SessionStore
}

// NewAuthMiddleware creates a new authentication middleware.
func NewAuthMiddleware(jwtManager *auth.JWTManager, sessionStore redis.SessionStore) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager:   jwtManager,
		sessionStore: sessionStore,
	}
}

// Authenticate returns a middleware that validates JWT tokens.
func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error":   "unauthorized",
				"message": "Missing authorization header",
			})
		}

		// Check for Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error":   "unauthorized",
				"message": "Invalid authorization header format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error":   "unauthorized",
				"message": "Missing token",
			})
		}

		// Validate token
		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
			})
		}

		// Check if token is revoked (if session store is available)
		if m.sessionStore != nil {
			isRevoked, err := m.sessionStore.IsRevoked(c.Request().Context(), claims.JTI)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error":   "internal_error",
					"message": "Failed to verify session",
				})
			}
			if isRevoked {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "unauthorized",
					"message": "Token has been revoked",
				})
			}
		}

		// Set user info in context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeyJTI, claims.JTI)
		c.Set(ContextKeyClaims, claims)

		return next(c)
	}
}

// GetUserID retrieves the user ID from the context.
func GetUserID(c echo.Context) uint {
	userID, ok := c.Get(ContextKeyUserID).(uint)
	if !ok {
		return 0
	}
	return userID
}

// GetUsername retrieves the username from the context.
func GetUsername(c echo.Context) string {
	username, ok := c.Get(ContextKeyUsername).(string)
	if !ok {
		return ""
	}
	return username
}

// GetJTI retrieves the JWT ID from the context.
func GetJTI(c echo.Context) string {
	jti, ok := c.Get(ContextKeyJTI).(string)
	if !ok {
		return ""
	}
	return jti
}

// GetClaims retrieves the full claims from the context.
func GetClaims(c echo.Context) *auth.AccessTokenClaims {
	claims, ok := c.Get(ContextKeyClaims).(*auth.AccessTokenClaims)
	if !ok {
		return nil
	}
	return claims
}
