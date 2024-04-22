package main

import (
	"log"
	"net/http"
	"strings"
)

type responseRefresh struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) postRefresh(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

	claims, err := validateToken(cfg.jwtSecret, tokenString)
	if err != nil {
		log.Printf("Token validation error:\n%v", err)
		w.WriteHeader(401)
		return
	}

	if claims.Issuer != "chirpy-refresh" {
		log.Printf("Invalid token type: %s", claims.Issuer)
		w.WriteHeader(401)
		return
	}

	valid, err := cfg.DB.ValidateToken(tokenString)
	if err != nil {
		log.Printf("Error validating refresh token:\n%v", err)
		w.WriteHeader(500)
		return
	}
	if !valid {
		log.Printf("The refresh token was revoked")
		w.WriteHeader(401)
		return
	}

	acessTokenString, err := generateToken(cfg.jwtSecret, "chirpy-access", claims.Subject, 3600)
	if err != nil {
		log.Printf("Error creating access token:\n%v", err)
		w.WriteHeader(500)
		return
	}

	respondWithJSON(w, 200, responseRefresh{Token: acessTokenString})
}
