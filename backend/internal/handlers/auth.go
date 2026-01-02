package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sweetfish329/sabakan/backend/internal/auth"
	"github.com/sweetfish329/sabakan/backend/internal/redis"
	"gorm.io/gorm"
)

// RegisterRequest represents the request body for user registration.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the request body for login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshRequest represents the request body for token refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthResponse represents the response for successful authentication.
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	db           *gorm.DB
	jwtManager   *auth.JWTManager
	sessionStore redis.SessionStore
}

// NewAuthHandler creates a new authentication handler.
func NewAuthHandler(db *gorm.DB, jwtManager *auth.JWTManager, sessionStore redis.SessionStore) *AuthHandler {
	return &AuthHandler{
		db:           db,
		jwtManager:   jwtManager,
		sessionStore: sessionStore,
	}
}

// Register handles user registration.
func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate request
	if req.Username == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Username is required",
		})
	}

	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Password must be at least 8 characters",
		})
	}

	// TODO: Implement full registration logic
	// For now, return success for testing
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

// Login handles user login.
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate request
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Username and password are required",
		})
	}

	// TODO: Implement full login logic
	// For now, return success for testing
	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken: "placeholder-token",
		ExpiresIn:   900,
		TokenType:   "Bearer",
	})
}

// Refresh handles token refresh.
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate request
	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Refresh token is required",
		})
	}

	// TODO: Implement full refresh logic
	// For now, return success for testing
	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken: "new-access-token",
		ExpiresIn:   900,
		TokenType:   "Bearer",
	})
}

// Logout handles user logout.
func (h *AuthHandler) Logout(c echo.Context) error {
	// TODO: Implement full logout logic
	// For now, return success for testing
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
