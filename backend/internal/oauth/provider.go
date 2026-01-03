// Package oauth provides OAuth authentication with external providers.
package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sweetfish329/sabakan/backend/internal/config"
)

// Common errors for OAuth providers.
var (
	ErrInvalidState       = errors.New("invalid state parameter")
	ErrCodeExchangeFailed = errors.New("failed to exchange code for token")
	ErrUserInfoFailed     = errors.New("failed to fetch user info")
)

// UserInfo represents user information from an OAuth provider.
type UserInfo struct {
	ProviderID string // Unique ID from the provider
	Email      string // User's email address
	Name       string // Display name
	AvatarURL  string // Profile picture URL
}

// Provider defines the interface for OAuth providers.
type Provider interface {
	// Name returns the provider name (e.g., "google", "discord").
	Name() string

	// AuthURL returns the authorization URL with the given state.
	AuthURL(state string) string

	// Exchange exchanges the authorization code for tokens and user info.
	Exchange(ctx context.Context, code string) (*UserInfo, error)
}

// TokenResponse represents the OAuth token response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// BaseProvider contains common OAuth provider logic.
type BaseProvider struct {
	name         string
	clientID     string
	clientSecret string
	redirectURL  string
	authURL      string
	tokenURL     string
	scopes       []string
}

// Name returns the provider name.
func (p *BaseProvider) Name() string {
	return p.name
}

// AuthURL returns the authorization URL.
func (p *BaseProvider) AuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", p.clientID)
	params.Set("redirect_uri", p.redirectURL)
	params.Set("response_type", "code")
	params.Set("scope", strings.Join(p.scopes, " "))
	params.Set("state", state)

	return fmt.Sprintf("%s?%s", p.authURL, params.Encode())
}

// ExchangeCode exchanges the authorization code for an access token.
func (p *BaseProvider) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", p.clientID)
	data.Set("client_secret", p.clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", p.redirectURL)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrCodeExchangeFailed, string(body))
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

// NewProviderFromConfig creates a Provider based on the provider name and config.
func NewProviderFromConfig(name string, cfg *config.OAuthConfig) (Provider, error) {
	switch name {
	case "google":
		return NewGoogleProvider(cfg.Google), nil
	case "discord":
		return NewDiscordProvider(cfg.Discord), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
