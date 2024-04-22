package main

import (
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) postRevoke(w http.ResponseWriter, r *http.Request) {
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

	err = cfg.DB.RevokeToken((tokenString))
	if err != nil {
		log.Printf("Error revoking refresh token:\n%v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}
