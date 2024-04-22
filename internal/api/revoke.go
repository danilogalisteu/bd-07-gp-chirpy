package api

import (
	"log"
	"net/http"
	"strings"
)

func (cfg *ApiConfig) PostRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

	_, err := validateToken(cfg.JwtSecret, tokenString, "chirpy-refresh")
	if err != nil {
		log.Printf("Token validation error:\n%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.DB.RevokeToken((tokenString))
	if err != nil {
		log.Printf("Error revoking refresh token:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
