package api

import (
	"encoding/json"
	"internal/database"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type paramRequest struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Invalid JSON: %s", err)
		respondWithJSON(w, http.StatusBadRequest, returnError{Error: "Invalid JSON"})
		return
	}

	dbUser, err := cfg.DbQueries.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     params.Email,
	})
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			log.Printf("Email already in use: %s", err)
			respondWithJSON(w, http.StatusConflict, returnError{Error: "Email already in use"})
			return
		}
		log.Printf("Error creating user: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	resUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, resUser)
}
