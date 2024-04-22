package main

import (
	"encoding/json"
	"internal/database"
	"log"
	"net/http"
	"strconv"
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
	if (err == database.ErrUserEmailNotFound) || (err == database.ErrUserInfoNotValid) {
		log.Printf("Logical error validating user on DB:\n%v", err)
		w.WriteHeader(401)
		return
	}
	if err != nil {
		log.Printf("Server error validating user on DB:\n%v", err)
		w.WriteHeader(500)
		return
	}

	expDuration := 24 * 3600
	if params.ExpiresInSeconds > 0 {
		expDuration = min(24*3600, params.ExpiresInSeconds)
	}

	tokenString, err := generateToken(cfg.jwtSecret, "chirpy", strconv.Itoa(user.ID), expDuration)
	if err != nil {
		log.Printf("Error creating auth token:\n%v", err)
		w.WriteHeader(500)
		return
	}

	respondWithJSON(w, 200, responseAuth{ID: user.ID, Email: user.Email, Token: tokenString})
}
