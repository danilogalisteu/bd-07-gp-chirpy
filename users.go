package main

import (
	"encoding/json"
	"internal/database"
	"log"
	"net/http"
	"strconv"
)

type paramUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BasicUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
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

	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err == database.ErrUserExists {
		log.Printf("User already exists on DB:\n%v", err)
		w.WriteHeader(403)
		return
	}
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
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Printf("Error converting requested id '%s' to number:\n%v", strId, err)
		respondWithError(w, 400, "ID was not recognized as number")
		return
	}

	user, err := cfg.DB.GetUserById(id)
	if err == database.ErrUserIdNotFound {
		respondWithError(w, 404, "ID was not found")
		return
	}
	if err != nil {
		log.Printf("Error getting user from DB:\n%v", err)
		w.WriteHeader(500)
		return
	}

	respondWithJSON(w, 200, BasicUser{ID: user.ID, Email: user.Email})
}
