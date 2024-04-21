package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type paramLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	respondWithJSON(w, 200, BasicUser{ID: user.ID, Email: user.Email})
}
