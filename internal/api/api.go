package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danilogalisteu/bd-07-gp-chirpy/internal/database"
)

type ApiConfig struct {
	JwtSecret      string
	PolkaApiKey    string
	FileserverHits int
	DbQueries      *database.Queries
}

type returnError struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	dat, err := json.Marshal(returnError{Error: msg})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
