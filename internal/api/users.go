package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/danilogalisteu/bd-07-gp-chirpy/internal/database"

	"github.com/danilogalisteu/bd-07-gp-chirpy/internal/auth"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type paramRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Invalid JSON: %s", err)
		respondWithJSON(w, http.StatusBadRequest, returnError{Error: "Invalid JSON"})
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	dbUser, err := cfg.DbQueries.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          params.Email,
		HashedPassword: hash,
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

func (cfg *ApiConfig) GetUser(w http.ResponseWriter, r *http.Request) {
	type paramRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Invalid JSON: %s", err)
		respondWithJSON(w, http.StatusBadRequest, returnError{Error: "Invalid JSON"})
		return
	}

	dbUser, err := cfg.DbQueries.GetUser(r.Context(), params.Email)
	if err != nil {
		if err.Error() == "pq: no rows in result set" {
			log.Printf("User not found: %s", err)
			respondWithJSON(w, http.StatusNotFound, returnError{Error: "User not found"})
			return
		}
		log.Printf("Error getting user: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	if err := auth.CheckPasswordHash(dbUser.HashedPassword, params.Password); err != nil {
		log.Printf("Invalid password: %s", err)
		respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Invalid email or password"})
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.JwtSecret, time.Duration(3600)*time.Second)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	dbToken, err := cfg.DbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(time.Duration(60*24) * time.Hour),
	})
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	resUser := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		IsChirpyRed:  dbUser.IsChirpyRed,
		Token:        token,
		RefreshToken: dbToken.Token,
	}
	respondWithJSON(w, http.StatusOK, resUser)
}

func (cfg *ApiConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {
	type paramRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	dbUser, err := cfg.DbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hash,
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		log.Printf("Error updating user: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	resUser := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, resUser)
}

func (cfg *ApiConfig) UpdateUserRed(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error getting API key: %s", err)
		respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Unauthorized"})
		return
	}
	if apiKey != cfg.PolkaApiKey {
		log.Printf("Invalid API key: %s", apiKey)
		respondWithJSON(w, http.StatusUnauthorized, returnError{Error: "Unauthorized"})
		return
	}

	type paramRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Invalid JSON: %s", err)
		respondWithJSON(w, http.StatusBadRequest, returnError{Error: "Invalid JSON"})
		return
	}

	if params.Event != "user.upgraded" {
		log.Printf("Unhandled event: %s", params.Event)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Invalid user ID: %s", err)
		respondWithJSON(w, http.StatusBadRequest, returnError{Error: "Invalid user ID"})
		return
	}

	_, err = cfg.DbQueries.UpdateUserRed(r.Context(), database.UpdateUserRedParams{
		ID:          userID,
		UpdatedAt:   time.Now(),
		IsChirpyRed: true,
	})
	if err != nil {
		if err.Error() == "pq: no rows in result set" {
			log.Printf("User not found: %s", err)
			respondWithJSON(w, http.StatusNotFound, returnError{Error: "User not found"})
			return
		}
		log.Printf("Error updating user: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, returnError{Error: "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
