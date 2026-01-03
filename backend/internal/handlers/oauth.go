package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sweetfish329/sabakan/backend/internal/auth"
	"github.com/sweetfish329/sabakan/backend/internal/config"
	"github.com/sweetfish329/sabakan/backend/internal/models"
	"github.com/sweetfish329/sabakan/backend/internal/oauth"
	"github.com/sweetfish329/sabakan/backend/internal/redis"
	"gorm.io/gorm"
)

// OAuthHandler handles OAuth authentication endpoints.
type OAuthHandler struct {
	db           *gorm.DB
	jwtManager   *auth.JWTManager
	sessionStore redis.SessionStore
	oauthConfig  *config.OAuthConfig
	frontendURL  string
}

// NewOAuthHandler creates a new OAuth handler.
func NewOAuthHandler(
	db *gorm.DB,
	jwtManager *auth.JWTManager,
	sessionStore redis.SessionStore,
	oauthConfig *config.OAuthConfig,
	frontendURL string,
) *OAuthHandler {
	return &OAuthHandler{
		db:           db,
		jwtManager:   jwtManager,
		sessionStore: sessionStore,
		oauthConfig:  oauthConfig,
		frontendURL:  frontendURL,
	}
}

// Authorize redirects to the OAuth provider's authorization URL.
func (h *OAuthHandler) Authorize(c echo.Context) error {
	providerName := c.Param("provider")

	provider, err := oauth.NewProviderFromConfig(providerName, h.oauthConfig)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_provider",
			Message: "Unsupported OAuth provider",
		})
	}

	// Generate state for CSRF protection
	state := generateState()

	// Store state in cookie
	c.SetCookie(&http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		MaxAge:   300,   // 5 minutes
		SameSite: http.SameSiteLaxMode,
	})

	return c.Redirect(http.StatusTemporaryRedirect, provider.AuthURL(state))
}

// Callback handles the OAuth provider's callback.
func (h *OAuthHandler) Callback(c echo.Context) error {
	providerName := c.Param("provider")

	provider, err := oauth.NewProviderFromConfig(providerName, h.oauthConfig)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?error=invalid_provider")
	}

	// Verify state
	state := c.QueryParam("state")
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || stateCookie.Value != state {
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?error=invalid_state")
	}

	// Clear state cookie
	c.SetCookie(&http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Exchange code for user info
	code := c.QueryParam("code")
	if code == "" {
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?error=missing_code")
	}

	userInfo, err := provider.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?error=exchange_failed")
	}

	// Find or create user
	user, err := h.findOrCreateUser(providerName, userInfo)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?error=user_creation_failed")
	}

	// Generate tokens
	accessToken, jti, err := h.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?error=token_generation_failed")
	}

	familyID := uuid.New().String()
	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID, familyID)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?error=token_generation_failed")
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

	// Redirect to frontend with tokens
	redirectURL := h.frontendURL + "?access_token=" + accessToken + "&refresh_token=" + refreshToken
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// findOrCreateUser finds an existing user by OAuth account or creates a new one.
func (h *OAuthHandler) findOrCreateUser(providerName string, userInfo *oauth.UserInfo) (*models.User, error) {
	// First, try to find existing OAuth account
	var oauthAccount models.OAuthAccount
	err := h.db.Where("provider = ? AND provider_user_id = ?", providerName, userInfo.ProviderID).First(&oauthAccount).Error
	if err == nil {
		// Found existing OAuth account, get the user
		var user models.User
		if err := h.db.First(&user, oauthAccount.UserID).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	// If user not found by OAuth, try to find by email and link
	if userInfo.Email != "" {
		var existingUser models.User
		if err := h.db.Where("email = ?", userInfo.Email).First(&existingUser).Error; err == nil {
			// Link OAuth account to existing user
			newOAuthAccount := models.OAuthAccount{
				UserID:         existingUser.ID,
				Provider:       providerName,
				ProviderUserID: userInfo.ProviderID,
				Email:          userInfo.Email,
			}
			if err := h.db.Create(&newOAuthAccount).Error; err != nil {
				return nil, err
			}
			return &existingUser, nil
		}
	}

	// Create new user
	// Get default user role
	var userRole models.Role
	if err := h.db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		return nil, err
	}

	// Generate unique username
	username := userInfo.Name
	if username == "" {
		username = "user_" + userInfo.ProviderID[:8]
	}

	// Check for username uniqueness and append suffix if needed
	var count int64
	h.db.Model(&models.User{}).Where("username LIKE ?", username+"%").Count(&count)
	if count > 0 {
		username = username + "_" + userInfo.ProviderID[:4]
	}

	var email *string
	if userInfo.Email != "" {
		email = &userInfo.Email
	}

	newUser := models.User{
		Username: username,
		Email:    email,
		RoleID:   userRole.ID,
		IsActive: true,
	}

	if err := h.db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	// Create OAuth account link
	newOAuthAccount := models.OAuthAccount{
		UserID:         newUser.ID,
		Provider:       providerName,
		ProviderUserID: userInfo.ProviderID,
		Email:          userInfo.Email,
	}
	if err := h.db.Create(&newOAuthAccount).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

// generateState generates a random state string for CSRF protection.
func generateState() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
