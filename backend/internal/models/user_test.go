package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRefreshToken_IsRevoked(t *testing.T) {
	t.Run("should return false when RevokedAt is nil", func(t *testing.T) {
		rt := &RefreshToken{
			RevokedAt: nil,
		}
		assert.False(t, rt.IsRevoked())
	})

	t.Run("should return true when RevokedAt is set", func(t *testing.T) {
		revokedTime := time.Now()
		rt := &RefreshToken{
			RevokedAt: &revokedTime,
		}
		assert.True(t, rt.IsRevoked())
	})
}

func TestRefreshToken_IsExpired(t *testing.T) {
	t.Run("should return false when token is not expired", func(t *testing.T) {
		rt := &RefreshToken{
			ExpiresAt: time.Now().Add(time.Hour),
		}
		assert.False(t, rt.IsExpired())
	})

	t.Run("should return true when token is expired", func(t *testing.T) {
		rt := &RefreshToken{
			ExpiresAt: time.Now().Add(-time.Hour),
		}
		assert.True(t, rt.IsExpired())
	})
}

func TestRefreshToken_IsValid(t *testing.T) {
	t.Run("should return true when not revoked and not expired", func(t *testing.T) {
		rt := &RefreshToken{
			ExpiresAt: time.Now().Add(time.Hour),
			RevokedAt: nil,
		}
		assert.True(t, rt.IsValid())
	})

	t.Run("should return false when revoked", func(t *testing.T) {
		revokedTime := time.Now()
		rt := &RefreshToken{
			ExpiresAt: time.Now().Add(time.Hour),
			RevokedAt: &revokedTime,
		}
		assert.False(t, rt.IsValid())
	})

	t.Run("should return false when expired", func(t *testing.T) {
		rt := &RefreshToken{
			ExpiresAt: time.Now().Add(-time.Hour),
			RevokedAt: nil,
		}
		assert.False(t, rt.IsValid())
	})

	t.Run("should return false when both revoked and expired", func(t *testing.T) {
		revokedTime := time.Now()
		rt := &RefreshToken{
			ExpiresAt: time.Now().Add(-time.Hour),
			RevokedAt: &revokedTime,
		}
		assert.False(t, rt.IsValid())
	})
}
