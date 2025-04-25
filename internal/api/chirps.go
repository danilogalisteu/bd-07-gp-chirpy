package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/danilogalisteu/bd-07-gp-chirpy/internal/auth"
	"github.com/danilogalisteu/bd-07-gp-chirpy/internal/database"

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
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Unauthorized"})
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.JwtSecret)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Unauthorized"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := paramRequest{}
	err = decoder.Decode(&params)
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
		UserID:    userID,
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

func (cfg *ApiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		log.Printf("Missing chirpID")
		respondWithJSON(w, http.StatusBadRequest, returnError{Error: "Missing chirp ID"})
		return
	}

	dbChirp, err := cfg.DbQueries.GetChirp(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Printf("Chirp not found: %s", err)
			respondWithJSON(w, http.StatusNotFound, returnError{Error: "Chirp not found"})
			return
		}
		log.Printf("Error getting chirp: %s", err)
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
	respondWithJSON(w, http.StatusOK, resChirp)
}
