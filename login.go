package main

import (
	"encoding/json"
	"internal/database"
	"log"
	"net/http"
	"strconv"
)

type paramLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type responseAuth struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) postLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := paramLogin{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := cfg.DB.ValidateUser(params.Email, params.Password)
	if (err == database.ErrUserEmailNotFound) || (err == database.ErrUserInfoNotValid) {
		log.Printf("Logical error validating user on DB:\n%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("Server error validating user on DB:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	acessTokenString, err := generateToken(cfg.jwtSecret, "chirpy-access", strconv.Itoa(user.ID), 3600)
	if err != nil {
		log.Printf("Error creating access token:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshTokenString, err := generateToken(cfg.jwtSecret, "chirpy-refresh", strconv.Itoa(user.ID), 60*24*3600)
	if err != nil {
		log.Printf("Error creating refresh token:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = cfg.DB.CreateToken(refreshTokenString)
	if err != nil {
		if err != database.ErrTokenExists { // possible collision because of 1 second precision
			log.Printf("Error storing refresh token:\n%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, responseAuth{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed, Token: acessTokenString, RefreshToken: refreshTokenString})
}
