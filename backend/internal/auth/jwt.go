package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ErrInvalidToken is returned when token validation fails.
var ErrInvalidToken = errors.New("invalid token")

// ErrExpiredToken is returned when token has expired.
var ErrExpiredToken = errors.New("token has expired")

// AccessTokenClaims represents the claims in an access token.
type AccessTokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	JTI      string `json:"jti"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents the claims in a refresh token.
type RefreshTokenClaims struct {
	UserID   uint   `json:"user_id"`
	FamilyID string `json:"family_id"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token generation and validation.
type JWTManager struct {
	secret             []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// NewJWTManager creates a new JWT manager.
func NewJWTManager(secret string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secret:             []byte(secret),
		accessTokenExpiry:  accessExpiry,
		refreshTokenExpiry: refreshExpiry,
	}
}

// GenerateAccessToken creates a new access token for the given user.
func (m *JWTManager) GenerateAccessToken(userID uint, username string) (string, string, error) {
	jti := uuid.New().String()
	now := time.Now()

	claims := AccessTokenClaims{
		UserID:   userID,
		Username: username,
		JTI:      jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "sabakan",
			Subject:   username,
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", "", err
	}

	return signedToken, jti, nil
}

// GenerateRefreshToken creates a new refresh token for the given user.
func (m *JWTManager) GenerateRefreshToken(userID uint, familyID string) (string, error) {
	now := time.Now()
	jti := uuid.New().String()

	claims := RefreshTokenClaims{
		UserID:   userID,
		FamilyID: familyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "sabakan",
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateAccessToken validates an access token and returns its claims.
func (m *JWTManager) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(_ *jwt.Token) (any, error) {
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token and returns its claims.
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(_ *jwt.Token) (any, error) {
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
