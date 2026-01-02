package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sweetfish329/sabakan/backend/internal/auth"
	"github.com/sweetfish329/sabakan/backend/internal/models"
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

	// Check if username already exists
	var existingUser models.User
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "conflict",
			Message: "Username already exists",
		})
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to process password",
		})
	}

	// Get default role (user)
	var userRole models.Role
	if err := h.db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get default role",
		})
	}

	// Create user
	var email *string
	if req.Email != "" {
		email = &req.Email
	}

	user := models.User{
		Username:     req.Username,
		Email:        email,
		PasswordHash: passwordHash,
		RoleID:       userRole.ID,
		IsActive:     true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"message": "User registered successfully",
		"user_id": user.ID,
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

	// Find user
	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Invalid credentials",
		})
	}

	// Check if user is active
	if !user.IsActive {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Account is disabled",
		})
	}

	// Verify password
	if !auth.VerifyPassword(req.Password, user.PasswordHash) {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Invalid credentials",
		})
	}

	// Generate tokens
	accessToken, jti, err := h.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate access token",
		})
	}

	// Generate refresh token with family ID
	familyID := uuid.New().String()
	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID, familyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate refresh token",
		})
	}

	// Store session in Redis if available
	if h.sessionStore != nil {
		sessionData := &redis.SessionData{
			UserID:    user.ID,
			IPAddress: c.RealIP(),
			UserAgent: c.Request().UserAgent(),
		}
		_ = h.sessionStore.StoreSession(c.Request().Context(), jti, sessionData, 15*time.Minute)
	}

	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
		TokenType:    "Bearer",
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

	// Validate refresh token
	claims, err := h.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Invalid or expired refresh token",
		})
	}

	// Get user
	var user models.User
	if err := h.db.First(&user, claims.UserID).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not found",
		})
	}

	// Check if user is still active
	if !user.IsActive {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Account is disabled",
		})
	}

	// Generate new access token
	accessToken, jti, err := h.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate access token",
		})
	}

	// Store new session in Redis if available
	if h.sessionStore != nil {
		sessionData := &redis.SessionData{
			UserID:    user.ID,
			IPAddress: c.RealIP(),
			UserAgent: c.Request().UserAgent(),
		}
		_ = h.sessionStore.StoreSession(c.Request().Context(), jti, sessionData, 15*time.Minute)
	}

	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken: accessToken,
		ExpiresIn:   900,
		TokenType:   "Bearer",
	})
}

// Logout handles user logout.
func (h *AuthHandler) Logout(c echo.Context) error {
	// Revoke session in Redis if available
	if h.sessionStore != nil {
		// Get JTI from context (set by auth middleware)
		jti, ok := c.Get("jti").(string)
		if ok && jti != "" {
			_ = h.sessionStore.RevokeSession(c.Request().Context(), jti, 24*time.Hour)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
