package main

import (
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) postRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

	_, err := validateToken(cfg.jwtSecret, tokenString, "chirpy-refresh")
	if err != nil {
		log.Printf("Token validation error:\n%v", err)
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
