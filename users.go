package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type paramUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BasicUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
}

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := paramUser{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, string(passwordHash))
	if err != nil {
		log.Printf("Error creating user on DB:\n%v", err)
		w.WriteHeader(500)
		return
	}

	respondWithJSON(w, 201, BasicUser{ID: user.ID, Email: user.Email})
}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.DB.GetUsers()
	if err != nil {
		log.Printf("Error getting users from DB:\n%v", err)
	}

	response := make([]BasicUser, 0)
	for _, user := range users {
		response = append(response, BasicUser{ID: user.ID, Email: user.Email})
	}

	respondWithJSON(w, 200, response)
}

func (cfg *apiConfig) getUserById(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.DB.GetUsers()
	if err != nil {
		log.Printf("Error getting users from DB:\n%v", err)
		w.WriteHeader(500)
		return
	}

	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Printf("Error converting requested id '%s' to number:\n%v", strId, err)
		respondWithError(w, 400, "ID was not recognized as number")
		return
	}

	for _, user := range users {
		if user.ID == id {
			respondWithJSON(w, 200, BasicUser{ID: user.ID, Email: user.Email})
			return
		}
	}

	respondWithError(w, 404, "ID was not found")
}
