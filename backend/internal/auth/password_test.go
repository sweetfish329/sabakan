package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		password := "securePassword123!"
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("should produce different hashes for same password", func(t *testing.T) {
		password := "securePassword123!"
		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2, "bcrypt should produce different hashes due to salt")
	})

	t.Run("should handle empty password", func(t *testing.T) {
		hash, err := HashPassword("")

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
	})
}

func TestVerifyPassword(t *testing.T) {
	t.Run("should return true for valid password", func(t *testing.T) {
		password := "securePassword123!"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		valid := VerifyPassword(password, hash)
		assert.True(t, valid)
	})

	t.Run("should return false for invalid password", func(t *testing.T) {
		password := "securePassword123!"
		wrongPassword := "wrongPassword456!"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		valid := VerifyPassword(wrongPassword, hash)
		assert.False(t, valid)
	})

	t.Run("should return false for empty password against valid hash", func(t *testing.T) {
		password := "securePassword123!"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		valid := VerifyPassword("", hash)
		assert.False(t, valid)
	})

	t.Run("should return false for invalid hash format", func(t *testing.T) {
		valid := VerifyPassword("password", "invalid-hash")
		assert.False(t, valid)
	})

	t.Run("should return true for empty password with empty hash", func(t *testing.T) {
		hash, err := HashPassword("")
		require.NoError(t, err)

		valid := VerifyPassword("", hash)
		assert.True(t, valid)
	})
}
