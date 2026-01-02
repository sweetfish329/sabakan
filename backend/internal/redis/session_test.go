package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockSessionStore is a mock implementation for testing without Redis.
type MockSessionStore struct {
	sessions map[string]*SessionData
	revoked  map[string]bool
}

// NewMockSessionStore creates a new mock session store.
func NewMockSessionStore() *MockSessionStore {
	return &MockSessionStore{
		sessions: make(map[string]*SessionData),
		revoked:  make(map[string]bool),
	}
}

func (m *MockSessionStore) StoreSession(_ context.Context, jti string, data *SessionData, _ time.Duration) error {
	m.sessions[jti] = data
	return nil
}

func (m *MockSessionStore) GetSession(_ context.Context, jti string) (*SessionData, error) {
	data, exists := m.sessions[jti]
	if !exists {
		return nil, ErrSessionNotFound
	}
	return data, nil
}

func (m *MockSessionStore) RevokeSession(_ context.Context, jti string, _ time.Duration) error {
	delete(m.sessions, jti)
	m.revoked[jti] = true
	return nil
}

func (m *MockSessionStore) IsRevoked(_ context.Context, jti string) (bool, error) {
	return m.revoked[jti], nil
}

func (m *MockSessionStore) RevokeAllUserSessions(_ context.Context, _ uint) error {
	m.sessions = make(map[string]*SessionData)
	return nil
}

func TestSessionData(t *testing.T) {
	t.Run("should create session data with required fields", func(t *testing.T) {
		data := &SessionData{
			UserID:    123,
			IPAddress: "192.168.1.1",
			UserAgent: "Mozilla/5.0",
		}

		assert.Equal(t, uint(123), data.UserID)
		assert.Equal(t, "192.168.1.1", data.IPAddress)
		assert.Equal(t, "Mozilla/5.0", data.UserAgent)
	})
}

func TestMockSessionStore_StoreSession(t *testing.T) {
	store := NewMockSessionStore()

	t.Run("should store session successfully", func(t *testing.T) {
		ctx := context.Background()
		jti := "test-jti-123"
		data := &SessionData{
			UserID:    123,
			IPAddress: "192.168.1.1",
			UserAgent: "Mozilla/5.0",
		}

		err := store.StoreSession(ctx, jti, data, 15*time.Minute)

		require.NoError(t, err)
	})
}

func TestMockSessionStore_GetSession(t *testing.T) {
	store := NewMockSessionStore()
	ctx := context.Background()

	t.Run("should get existing session", func(t *testing.T) {
		jti := "test-jti-456"
		data := &SessionData{
			UserID:    456,
			IPAddress: "10.0.0.1",
			UserAgent: "Chrome/100",
		}
		_ = store.StoreSession(ctx, jti, data, 15*time.Minute)

		result, err := store.GetSession(ctx, jti)

		require.NoError(t, err)
		assert.Equal(t, uint(456), result.UserID)
		assert.Equal(t, "10.0.0.1", result.IPAddress)
	})

	t.Run("should return error for non-existent session", func(t *testing.T) {
		result, err := store.GetSession(ctx, "non-existent-jti")

		assert.ErrorIs(t, err, ErrSessionNotFound)
		assert.Nil(t, result)
	})
}

func TestMockSessionStore_RevokeSession(t *testing.T) {
	store := NewMockSessionStore()
	ctx := context.Background()

	t.Run("should revoke session successfully", func(t *testing.T) {
		jti := "test-jti-revoke"
		data := &SessionData{UserID: 789}
		_ = store.StoreSession(ctx, jti, data, 15*time.Minute)

		err := store.RevokeSession(ctx, jti, 24*time.Hour)

		require.NoError(t, err)

		// Session should no longer exist
		_, getErr := store.GetSession(ctx, jti)
		assert.ErrorIs(t, getErr, ErrSessionNotFound)
	})

	t.Run("should mark session as revoked", func(t *testing.T) {
		jti := "test-jti-revoked-check"
		data := &SessionData{UserID: 999}
		_ = store.StoreSession(ctx, jti, data, 15*time.Minute)
		_ = store.RevokeSession(ctx, jti, 24*time.Hour)

		isRevoked, err := store.IsRevoked(ctx, jti)

		require.NoError(t, err)
		assert.True(t, isRevoked)
	})
}

func TestMockSessionStore_IsRevoked(t *testing.T) {
	store := NewMockSessionStore()
	ctx := context.Background()

	t.Run("should return false for non-revoked session", func(t *testing.T) {
		isRevoked, err := store.IsRevoked(ctx, "never-revoked-jti")

		require.NoError(t, err)
		assert.False(t, isRevoked)
	})
}
