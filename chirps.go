package main

import (
	"encoding/json"
	"internal/database"
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

	claims, err := validateToken(cfg.jwtSecret, tokenString, "chirpy-access")
	if err != nil {
		log.Printf("Token validation error:\n%v", err)
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
	strAuthorId := r.URL.Query().Get("author_id")

	if strAuthorId == "" {
		chirps, err := cfg.DB.GetChirps()
		if err != nil {
			log.Printf("Error getting messages from DB:\n%v", err)
		}

		respondWithJSON(w, 200, chirps)
	} else {
		authorId, err := strconv.Atoi(strAuthorId)
		if err != nil {
			log.Printf("Error converting requested id '%d' to number:\n%v", authorId, err)
			respondWithError(w, 400, "Author ID was not recognized as number")
			return
		}

		chirps, err := cfg.DB.GetChirpsByAuthor(authorId)
		if err != nil {
			log.Printf("Error getting messages from DB:\n%v", err)
			w.WriteHeader(500)
		}

		respondWithJSON(w, 200, chirps)
	}
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

func (cfg *apiConfig) deleteChirpById(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

	claims, err := validateToken(cfg.jwtSecret, tokenString, "chirpy-access")
	if err != nil {
		log.Printf("Token validation error:\n%v", err)
		w.WriteHeader(401)
		return
	}

	user_id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Printf("Invalid token ID value:\n%v", err)
		w.WriteHeader(401)
		return
	}

	strId := r.PathValue("id")
	chirp_id, err := strconv.Atoi(strId)
	if err != nil {
		log.Printf("Error converting requested id '%s' to number:\n%v", strId, err)
		respondWithError(w, 400, "ID was not recognized as number")
		return
	}

	err = cfg.DB.DeleteChirpById(chirp_id, user_id)
	if err == database.ErrChirpIdNotFound {
		log.Printf("Chirp ID was not found: %d", chirp_id)
		w.WriteHeader(404)
		return
	}
	if err == database.ErrChirpAuthorInvalid {
		log.Printf("User ID '%d' is different from author ID", user_id)
		w.WriteHeader(403)
		return
	}
	if err != nil {
		log.Printf("Error deleting chirp on DB:\n%v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}
