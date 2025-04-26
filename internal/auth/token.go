package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key), nil
}

func GetBearerToken(headers http.Header) (string, error) {
	// Get the Authorization header
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", jwt.ErrTokenUnverifiable
	}

	// Check if the header starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", jwt.ErrTokenUnverifiable
	}

	// Return the token
	return authHeader[7:], nil
}
