package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrSessionNotFound is returned when a session does not exist.
var ErrSessionNotFound = errors.New("session not found")

// SessionData represents the data stored for a session.
type SessionData struct {
	UserID    uint   `json:"user_id"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

// SessionStore defines the interface for session storage.
type SessionStore interface {
	StoreSession(ctx context.Context, jti string, data *SessionData, expiry time.Duration) error
	GetSession(ctx context.Context, jti string) (*SessionData, error)
	RevokeSession(ctx context.Context, jti string, blacklistTTL time.Duration) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
	RevokeAllUserSessions(ctx context.Context, userID uint) error
}

// RedisSessionStore implements SessionStore using Redis.
type RedisSessionStore struct {
	client *Client
}

// NewRedisSessionStore creates a new Redis-based session store.
func NewRedisSessionStore(client *Client) *RedisSessionStore {
	return &RedisSessionStore{client: client}
}

// sessionKey returns the Redis key for a session.
func sessionKey(jti string) string {
	return fmt.Sprintf("session:%s", jti)
}

// revokedKey returns the Redis key for a revoked session.
func revokedKey(jti string) string {
	return fmt.Sprintf("revoked:%s", jti)
}

// userSessionsKey returns the Redis key for a user's session set.
func userSessionsKey(userID uint) string {
	return fmt.Sprintf("user:%d:sessions", userID)
}

// StoreSession stores a session in Redis.
func (s *RedisSessionStore) StoreSession(ctx context.Context, jti string, data *SessionData, expiry time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	pipe := s.client.rdb.Pipeline()

	// Store session data
	pipe.Set(ctx, sessionKey(jti), jsonData, expiry)

	// Add JTI to user's session set
	pipe.SAdd(ctx, userSessionsKey(data.UserID), jti)

	_, err = pipe.Exec(ctx)
	return err
}

// GetSession retrieves a session from Redis.
func (s *RedisSessionStore) GetSession(ctx context.Context, jti string) (*SessionData, error) {
	jsonData, err := s.client.rdb.Get(ctx, sessionKey(jti)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	var data SessionData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// RevokeSession removes a session and adds it to the blacklist.
func (s *RedisSessionStore) RevokeSession(ctx context.Context, jti string, blacklistTTL time.Duration) error {
	// Get session data to retrieve user ID
	data, err := s.GetSession(ctx, jti)
	if err != nil && !errors.Is(err, ErrSessionNotFound) {
		return err
	}

	pipe := s.client.rdb.Pipeline()

	// Delete session
	pipe.Del(ctx, sessionKey(jti))

	// Add to blacklist
	pipe.Set(ctx, revokedKey(jti), "1", blacklistTTL)

	// Remove from user's session set
	if data != nil {
		pipe.SRem(ctx, userSessionsKey(data.UserID), jti)
	}

	_, err = pipe.Exec(ctx)
	return err
}

// IsRevoked checks if a session has been revoked.
func (s *RedisSessionStore) IsRevoked(ctx context.Context, jti string) (bool, error) {
	exists, err := s.client.rdb.Exists(ctx, revokedKey(jti)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// RevokeAllUserSessions revokes all sessions for a user.
func (s *RedisSessionStore) RevokeAllUserSessions(ctx context.Context, userID uint) error {
	// Get all JTIs for the user
	jtis, err := s.client.rdb.SMembers(ctx, userSessionsKey(userID)).Result()
	if err != nil {
		return err
	}

	if len(jtis) == 0 {
		return nil
	}

	pipe := s.client.rdb.Pipeline()

	// Revoke each session
	for _, jti := range jtis {
		pipe.Del(ctx, sessionKey(jti))
		pipe.Set(ctx, revokedKey(jti), "1", 24*time.Hour)
	}

	// Clear user's session set
	pipe.Del(ctx, userSessionsKey(userID))

	_, err = pipe.Exec(ctx)
	return err
}
