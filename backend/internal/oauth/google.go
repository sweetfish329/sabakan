package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sweetfish329/sabakan/backend/internal/config"
)

const (
	googleAuthURL     = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// GoogleProvider implements OAuth for Google.
type GoogleProvider struct {
	BaseProvider
}

// NewGoogleProvider creates a new Google OAuth provider.
func NewGoogleProvider(cfg config.OAuthProviderConfig) *GoogleProvider {
	return &GoogleProvider{
		BaseProvider: BaseProvider{
			name:         "google",
			clientID:     cfg.ClientID,
			clientSecret: cfg.ClientSecret,
			redirectURL:  cfg.RedirectURL,
			authURL:      googleAuthURL,
			tokenURL:     googleTokenURL,
			scopes:       []string{"openid", "email", "profile"},
		},
	}
}

// googleUserInfo represents the Google user info response.
type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// Exchange exchanges the authorization code for user info.
func (p *GoogleProvider) Exchange(ctx context.Context, code string) (*UserInfo, error) {
	token, err := p.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Fetch user info from Google
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrUserInfoFailed, string(body))
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &UserInfo{
		ProviderID: info.ID,
		Email:      info.Email,
		Name:       info.Name,
		AvatarURL:  info.Picture,
	}, nil
}
