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
	if err == ErrTokenParsing {
		log.Printf("Invalid token parsing:\n%v", err)
		w.WriteHeader(401)
		return
	}
	if err == ErrTokenInvalid {
		log.Printf("Invalid token")
		w.WriteHeader(401)
		return
	}
	if err == ErrTokenClaimsParsing {
		log.Printf("Unable to extract token claims:\n%v", err)
		w.WriteHeader(401)
		return
	}

	if claims.Issuer != "chirpy-refresh" {
		log.Printf("Invalid token type:\n%v", err)
		w.WriteHeader(401)
		return
	}

	// TODO check token status

	acessTokenString, err := generateToken(cfg.jwtSecret, "chirpy-access", claims.Subject, 3600)
	if err != nil {
		log.Printf("Error creating access token:\n%v", err)
		w.WriteHeader(500)
		return
	}

	respondWithJSON(w, 200, responseRefresh{Token: acessTokenString})
}
