package api

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) MiddlewareMetricsCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	content := `<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	
	</html>`

	w.Write([]byte(fmt.Sprintf(content, cfg.FileserverHits)))
}

func (cfg *ApiConfig) MiddlewareMetricsReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits = 0

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	err := cfg.DbQueries.ResetUsers(r.Context())
	if err != nil {
		log.Printf("Error resetting users: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
