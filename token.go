package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenParsing = errors.New("token parsing failed")
var ErrTokenInvalid = errors.New("invalid token")
var ErrTokenClaimsParsing = errors.New("token claims parsing failed")

func generateToken(secret string, issuer string, subject string, expires_in_seconds int) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expires_in_seconds) * time.Second)),
			Subject:   subject,
		},
	)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateToken(secret string, tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrTokenParsing
	}
	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrTokenClaimsParsing
	}

	return claims, nil
}
