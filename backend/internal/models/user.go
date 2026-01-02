package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a system user.
type User struct {
	gorm.Model
	Username      string         `gorm:"uniqueIndex;not null" json:"username"`
	Email         *string        `gorm:"uniqueIndex" json:"email,omitempty"`
	PasswordHash  string         `json:"-"`
	RoleID        uint           `gorm:"not null;index" json:"roleId"`
	Role          Role           `json:"role,omitempty"`
	IsActive      bool           `gorm:"default:true" json:"isActive"`
	OAuthAccounts []OAuthAccount `json:"oauthAccounts,omitempty"`
	APITokens     []APIToken     `json:"-"`
	RefreshTokens []RefreshToken `json:"-"`
	GameServers   []GameServer   `gorm:"foreignKey:OwnerID" json:"gameServers,omitempty"`
}

// OAuthAccount represents a linked OAuth provider account.
type OAuthAccount struct {
	gorm.Model
	UserID         uint       `gorm:"not null;index" json:"userId"`
	User           User       `json:"-"`
	Provider       string     `gorm:"not null;uniqueIndex:idx_oauth_provider_user" json:"provider"`
	ProviderUserID string     `gorm:"not null;uniqueIndex:idx_oauth_provider_user" json:"providerUserId"`
	Email          string     `json:"email,omitempty"`
	DisplayName    string     `json:"displayName,omitempty"`
	AvatarURL      string     `json:"avatarUrl,omitempty"`
	AccessToken    string     `json:"-"` // Encrypted
	RefreshToken   string     `json:"-"` // Encrypted
	TokenExpiresAt *time.Time `json:"-"`
}

// APIToken represents an API access token.
type APIToken struct {
	gorm.Model
	UserID      uint       `gorm:"not null;index" json:"userId"`
	User        User       `json:"-"`
	Name        string     `gorm:"not null" json:"name"`
	TokenHash   string     `gorm:"uniqueIndex;not null" json:"-"`
	TokenPrefix string     `gorm:"not null" json:"tokenPrefix"` // First 8 chars for display
	Scopes      string     `json:"scopes,omitempty"`            // JSON array of allowed scopes
	LastUsedAt  *time.Time `json:"lastUsedAt,omitempty"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

// RefreshToken represents a long-lived refresh token for JWT rotation.
type RefreshToken struct {
	gorm.Model
	UserID    uint       `gorm:"not null;index" json:"-"`
	User      User       `json:"-"`
	TokenHash string     `gorm:"uniqueIndex;not null" json:"-"`
	FamilyID  string     `gorm:"not null;index" json:"-"` // For rotation tracking
	IPAddress string     `json:"-"`
	UserAgent string     `json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"-"`
	RevokedAt *time.Time `json:"-"` // Null if not revoked
}

// IsRevoked returns true if the refresh token has been revoked.
func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

// IsExpired returns true if the refresh token has expired.
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid returns true if the refresh token is not revoked and not expired.
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked() && !rt.IsExpired()
}
