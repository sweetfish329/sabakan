package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key-for-jwt-testing-32bytes!"

func TestNewJWTManager(t *testing.T) {
	t.Run("should create JWT manager with valid secret", func(t *testing.T) {
		manager := NewJWTManager(testSecret, 15*time.Minute, 7*24*time.Hour)

		assert.NotNil(t, manager)
	})
}

func TestGenerateAccessToken(t *testing.T) {
	manager := NewJWTManager(testSecret, 15*time.Minute, 7*24*time.Hour)

	t.Run("should generate access token with user claims", func(t *testing.T) {
		userID := uint(123)
		username := "testuser"

		token, jti, err := manager.GenerateAccessToken(userID, username)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, jti)
	})

	t.Run("should generate different tokens for same user", func(t *testing.T) {
		userID := uint(123)
		username := "testuser"

		token1, jti1, _ := manager.GenerateAccessToken(userID, username)
		token2, jti2, _ := manager.GenerateAccessToken(userID, username)

		assert.NotEqual(t, token1, token2)
		assert.NotEqual(t, jti1, jti2)
	})
}

func TestGenerateRefreshToken(t *testing.T) {
	manager := NewJWTManager(testSecret, 15*time.Minute, 7*24*time.Hour)

	t.Run("should generate refresh token", func(t *testing.T) {
		userID := uint(123)
		familyID := "test-family-id"

		token, err := manager.GenerateRefreshToken(userID, familyID)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestValidateAccessToken(t *testing.T) {
	manager := NewJWTManager(testSecret, 15*time.Minute, 7*24*time.Hour)

	t.Run("should validate valid access token", func(t *testing.T) {
		userID := uint(123)
		username := "testuser"
		token, jti, _ := manager.GenerateAccessToken(userID, username)

		claims, err := manager.ValidateAccessToken(token)

		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, jti, claims.JTI)
	})

	t.Run("should reject token with invalid signature", func(t *testing.T) {
		otherManager := NewJWTManager("different-secret-key-32bytes!!", 15*time.Minute, 7*24*time.Hour)
		token, _, _ := otherManager.GenerateAccessToken(123, "testuser")

		claims, err := manager.ValidateAccessToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("should reject malformed token", func(t *testing.T) {
		claims, err := manager.ValidateAccessToken("not-a-valid-jwt")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("should reject empty token", func(t *testing.T) {
		claims, err := manager.ValidateAccessToken("")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestValidateAccessToken_Expired(t *testing.T) {
	// Create manager with very short expiration
	manager := NewJWTManager(testSecret, 1*time.Millisecond, 7*24*time.Hour)

	t.Run("should reject expired token", func(t *testing.T) {
		token, _, _ := manager.GenerateAccessToken(123, "testuser")

		// Wait for token to expire
		time.Sleep(10 * time.Millisecond)

		claims, err := manager.ValidateAccessToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestValidateRefreshToken(t *testing.T) {
	manager := NewJWTManager(testSecret, 15*time.Minute, 7*24*time.Hour)

	t.Run("should validate valid refresh token", func(t *testing.T) {
		userID := uint(123)
		familyID := "test-family-id"
		token, _ := manager.GenerateRefreshToken(userID, familyID)

		claims, err := manager.ValidateRefreshToken(token)

		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, familyID, claims.FamilyID)
	})

	t.Run("should reject refresh token with invalid signature", func(t *testing.T) {
		otherManager := NewJWTManager("different-secret-key-32bytes!!", 15*time.Minute, 7*24*time.Hour)
		token, _ := otherManager.GenerateRefreshToken(123, "family-id")

		claims, err := manager.ValidateRefreshToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}
