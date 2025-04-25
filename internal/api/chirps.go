package api

import (
	"encoding/json"
	"internal/database"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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

func (cfg *ApiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	type paramRequest struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
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

	dbChirp, err := cfg.DbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      cleanMessage(params.Body),
		UserID:    uuid.MustParse(params.UserID),
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	resChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, resChirp)
}

func (cfg *ApiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DbQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	resChirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		resChirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}
	respondWithJSON(w, http.StatusOK, resChirps)
}
