package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type paramLogin struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type responseAuth struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *apiConfig) postLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := paramLogin{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.DB.ValidateUser(params.Email, params.Password)
	if err != nil {
		log.Printf("Error validating user on DB:\n%v", err)
		w.WriteHeader(401)
		return
	}

	now := time.Now()
	expDuration := 24 * 3600
	if params.ExpiresInSeconds > 0 {
		expDuration = min(24*3600, params.ExpiresInSeconds)
	}

	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Issuer":    "chirpy",
			"IssuedAt":  now,
			"ExpiresAt": now.Add(time.Duration(expDuration) * time.Second),
			"Subject":   strconv.Itoa(user.ID),
		},
	).SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		log.Printf("Error creating auth token:\n%v", err)
		w.WriteHeader(500)
		return
	}

	respondWithJSON(w, 200, responseAuth{ID: user.ID, Email: user.Email, Token: token})
}
