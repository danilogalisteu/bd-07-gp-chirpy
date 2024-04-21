package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type paramEmail struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := paramEmail{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.DB.CreateUser(params.Email)
	if err != nil {
		log.Printf("Error creating user on DB:\n%v", err)
	}

	respondWithJSON(w, 201, user)
}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.DB.GetUsers()
	if err != nil {
		log.Printf("Error getting users from DB:\n%v", err)
	}

	respondWithJSON(w, 200, users)
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
			respondWithJSON(w, 200, user)
			return
		}
	}

	respondWithError(w, 404, "ID was not found")
}
