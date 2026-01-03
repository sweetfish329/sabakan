package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sweetfish329/sabakan/backend/internal/config"
)

func TestNewGoogleProvider(t *testing.T) {
	cfg := config.OAuthProviderConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:1323/auth/oauth/google/callback",
	}

	provider := NewGoogleProvider(cfg)

	assert.Equal(t, "google", provider.Name())
	assert.Contains(t, provider.AuthURL("test-state"), "accounts.google.com")
	assert.Contains(t, provider.AuthURL("test-state"), "client_id=test-client-id")
	assert.Contains(t, provider.AuthURL("test-state"), "state=test-state")
}

func TestNewDiscordProvider(t *testing.T) {
	cfg := config.OAuthProviderConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:1323/auth/oauth/discord/callback",
	}

	provider := NewDiscordProvider(cfg)

	assert.Equal(t, "discord", provider.Name())
	assert.Contains(t, provider.AuthURL("test-state"), "discord.com")
	assert.Contains(t, provider.AuthURL("test-state"), "client_id=test-client-id")
	assert.Contains(t, provider.AuthURL("test-state"), "state=test-state")
}

func TestNewProviderFromConfig(t *testing.T) {
	cfg := &config.OAuthConfig{
		Google: config.OAuthProviderConfig{
			ClientID: "google-id",
		},
		Discord: config.OAuthProviderConfig{
			ClientID: "discord-id",
		},
	}

	t.Run("should return Google provider", func(t *testing.T) {
		provider, err := NewProviderFromConfig("google", cfg)
		assert.NoError(t, err)
		assert.Equal(t, "google", provider.Name())
	})

	t.Run("should return Discord provider", func(t *testing.T) {
		provider, err := NewProviderFromConfig("discord", cfg)
		assert.NoError(t, err)
		assert.Equal(t, "discord", provider.Name())
	})

	t.Run("should return error for unknown provider", func(t *testing.T) {
		_, err := NewProviderFromConfig("unknown", cfg)
		assert.Error(t, err)
	})
}
