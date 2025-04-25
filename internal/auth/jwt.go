package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
	})

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	// Check the claims
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// Check if the token is expired
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return uuid.Nil, jwt.ErrTokenExpired
		}
		// Convert the subject to a UUID
		userId, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, err
		}
		// Return the user ID
		return userId, nil
	}

	return uuid.Nil, jwt.ErrTokenSignatureInvalid
}
