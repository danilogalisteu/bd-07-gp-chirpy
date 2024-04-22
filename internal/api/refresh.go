package api

import (
	"log"
	"net/http"
	"strings"
)

type responseRefresh struct {
	Token string `json:"token"`
}

func (cfg *ApiConfig) PostRefresh(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

	claims, err := validateToken(cfg.JwtSecret, tokenString, "chirpy-refresh")
	if err != nil {
		log.Printf("Token validation error:\n%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	valid, err := cfg.DB.ValidateToken(tokenString)
	if err != nil {
		log.Printf("Error validating refresh token:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !valid {
		log.Printf("The refresh token was revoked")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	acessTokenString, err := generateToken(cfg.JwtSecret, "chirpy-access", claims.Subject, 3600)
	if err != nil {
		log.Printf("Error creating access token:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, responseRefresh{Token: acessTokenString})
}
