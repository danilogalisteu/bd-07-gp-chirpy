package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type paramBody struct {
	Body string `json:"body"`
}

func cleanMessage(msg string) string {
	profane := []string{"kerfuffle", "sharbert", "fornax"}
	clean := make([]string, 0)
	for _, word := range strings.Split(msg, " ") {
		for _, pword := range profane {
			if strings.ToLower(word) == pword {
				word = "****"
				break
			}
		}
		clean = append(clean, word)
	}
	return strings.Join(clean, " ")
}

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := paramBody{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

	claims, err := validateToken(cfg.jwtSecret, tokenString)
	if err != nil {
		log.Printf("Token validation error:\n%v", err)
		w.WriteHeader(401)
		return
	}

	if claims.Issuer != "chirpy-access" {
		log.Printf("Invalid token type: %s", claims.Issuer)
		w.WriteHeader(401)
		return
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Printf("Invalid token ID value:\n%v", err)
		w.WriteHeader(401)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleaned := cleanMessage(params.Body)

	chirp, err := cfg.DB.CreateChirp(id, cleaned)
	if err != nil {
		log.Printf("Error creating chirp on DB:\n%v", err)
	}

	respondWithJSON(w, 201, chirp)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		log.Printf("Error getting messages from DB:\n%v", err)
	}

	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		log.Printf("Error getting items from DB:\n%v", err)
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

	for _, chirp := range chirps {
		if chirp.ID == id {
			respondWithJSON(w, 200, chirp)
			return
		}
	}

	respondWithError(w, 404, "ID was not found")
}
