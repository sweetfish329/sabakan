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
	discordAuthURL     = "https://discord.com/api/oauth2/authorize"
	discordTokenURL    = "https://discord.com/api/oauth2/token"
	discordUserInfoURL = "https://discord.com/api/users/@me"
)

// DiscordProvider implements OAuth for Discord.
type DiscordProvider struct {
	BaseProvider
}

// NewDiscordProvider creates a new Discord OAuth provider.
func NewDiscordProvider(cfg config.OAuthProviderConfig) *DiscordProvider {
	return &DiscordProvider{
		BaseProvider: BaseProvider{
			name:         "discord",
			clientID:     cfg.ClientID,
			clientSecret: cfg.ClientSecret,
			redirectURL:  cfg.RedirectURL,
			authURL:      discordAuthURL,
			tokenURL:     discordTokenURL,
			scopes:       []string{"identify", "email"},
		},
	}
}

// discordUserInfo represents the Discord user info response.
type discordUserInfo struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
	Avatar        string `json:"avatar"`
}

// Exchange exchanges the authorization code for user info.
func (p *DiscordProvider) Exchange(ctx context.Context, code string) (*UserInfo, error) {
	token, err := p.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Fetch user info from Discord
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discordUserInfoURL, nil)
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

	var info discordUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	// Build avatar URL
	avatarURL := ""
	if info.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", info.ID, info.Avatar)
	}

	return &UserInfo{
		ProviderID: info.ID,
		Email:      info.Email,
		Name:       info.Username,
		AvatarURL:  avatarURL,
	}, nil
}
