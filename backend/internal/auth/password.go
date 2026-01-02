// Package auth provides authentication utilities including password hashing and JWT management.
package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// DefaultCost is the bcrypt cost factor for password hashing.
const DefaultCost = 12

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword checks if a password matches a bcrypt hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
