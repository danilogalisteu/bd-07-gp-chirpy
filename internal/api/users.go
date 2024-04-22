package api

import (
	"encoding/json"
	"internal/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type paramUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BasicUser struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (cfg *ApiConfig) PostUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := paramUser{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err == database.ErrUserExists {
		log.Printf("User already exists on DB:\n%v", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		log.Printf("Error creating user on DB:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, BasicUser{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed})
}

func (cfg *ApiConfig) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.DB.GetUsers()
	if err != nil {
		log.Printf("Error getting users from DB:\n%v", err)
	}

	response := make([]BasicUser, 0)
	for _, user := range users {
		response = append(response, BasicUser{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *ApiConfig) GetUserById(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Printf("Error converting requested id '%s' to number:\n%v", strId, err)
		respondWithError(w, http.StatusBadRequest, "ID was not recognized as number")
		return
	}

	user, err := cfg.DB.GetUserById(id)
	if err == database.ErrUserIdNotFound {
		respondWithError(w, http.StatusNotFound, "ID was not found")
		return
	}
	if err != nil {
		log.Printf("Error getting user from DB:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, BasicUser{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed})
}

func (cfg *ApiConfig) PutUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := paramUser{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

	claims, err := validateToken(cfg.JwtSecret, tokenString, "chirpy-access")
	if err != nil {
		log.Printf("Token validation error:\n%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Printf("Invalid token ID value:\n%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := cfg.DB.UpdateUserById(id, params.Email, params.Password)
	if err == database.ErrUserIdNotFound {
		log.Printf("ID was not found:\n%v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error updating user on DB:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, BasicUser{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed})
}
