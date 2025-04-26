package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/danilogalisteu/bd-07-gp-chirpy/internal/auth"
	"github.com/danilogalisteu/bd-07-gp-chirpy/internal/database"
)

func (cfg *ApiConfig) GetToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting refresh token: %s", err)
		respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Unauthorized"})
		return
	}

	dbToken, err := cfg.DbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if err.Error() == "pq: no rows in result set" {
			log.Printf("Refresh token not found: %s", err)
			respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Refresh token not found"})
			return
		}
		log.Printf("Error getting refresh token: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	if dbToken.ExpiresAt.Before(time.Now()) {
		log.Printf("Refresh token expired on %v", dbToken.ExpiresAt)
		respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Refresh token expired"})
		return
	}

	if dbToken.RevokedAt.Valid {
		log.Printf("Refresh token revoked on %v", dbToken.RevokedAt.Time)
		respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Refresh token revoked"})
		return
	}

	token, err := auth.MakeJWT(dbToken.UserID, cfg.JwtSecret, time.Duration(3600)*time.Second)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	type resultToken struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, resultToken{Token: token})
}

func (cfg *ApiConfig) UpdateToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting refresh token: %s", err)
		respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Unauthorized"})
		return
	}

	dbToken, err := cfg.DbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if err.Error() == "pq: no rows in result set" {
			log.Printf("Refresh token not found: %s", err)
			respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Refresh token not found"})
			return
		}
		log.Printf("Error getting refresh token: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	err = cfg.DbQueries.UpdateToken(r.Context(), database.UpdateTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Token: dbToken.Token,
	})
	if err != nil {
		log.Printf("Error revoking refresh token: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
