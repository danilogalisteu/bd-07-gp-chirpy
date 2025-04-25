package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

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

func (cfg *ApiConfig) ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type paramRequest struct {
		Body string `json:"body"`
	}

	type returnClean struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Invalid JSON: %s", err)
		respondWithJSON(w, http.StatusBadRequest, returnError{Error: "Invalid JSON"})
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp is too long: %d characters", len(params.Body))
		respondWithJSON(w, http.StatusBadRequest, returnError{Error: "Chirp is too long"})
		return
	}

	resValid := returnClean{
		CleanedBody: cleanMessage(params.Body),
	}
	respondWithJSON(w, http.StatusOK, resValid)
}
