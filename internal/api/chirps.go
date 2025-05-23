package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
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
	var userID uuid.NullUUID
	var dbChirps []database.Chirp
	var err error

	queryUserID := r.URL.Query().Get("author_id")
	if queryUserID != "" {
		if userID.UUID, err = uuid.Parse(queryUserID); err != nil {
			log.Printf("Invalid author_id: %s", err)
			respondWithJSON(w, http.StatusBadRequest, returnError{Error: "Invalid author_id"})
			return
		}
		userID.Valid = true
	}

	if userID.Valid {
		dbChirps, err = cfg.DbQueries.GetChirpsByUser(r.Context(), userID.UUID)
		if err != nil {
			log.Printf("Error getting chirps by user ID: %s", err)
			respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
			return
		}
	} else {
		dbChirps, err = cfg.DbQueries.GetChirps(r.Context())
		if err != nil {
			log.Printf("Error getting chirps: %s", err)
			respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
			return
		}
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

	querySort := r.URL.Query().Get("sort")
	if querySort == "desc" {
		sort.Slice(resChirps, func(i, j int) bool {
			return resChirps[i].CreatedAt.After(resChirps[j].CreatedAt)
		})
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

func (cfg *ApiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		log.Printf("Missing chirpID")
		respondWithJSON(w, http.StatusBadRequest, returnError{Error: "Missing chirp ID"})
		return
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

	if dbChirp.UserID != userID {
		log.Printf("User %s is not authorized to delete chirp %s", userID, chirpID)
		respondWithJSON(w, http.StatusForbidden, returnError{Error: "Forbidden"})
		return
	}

	err = cfg.DbQueries.DeleteChirp(r.Context(), dbChirp.ID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Printf("Chirp not found: %s", err)
			respondWithJSON(w, http.StatusNotFound, returnError{Error: "Chirp not found"})
			return
		}
		log.Printf("Error deleting chirp: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
