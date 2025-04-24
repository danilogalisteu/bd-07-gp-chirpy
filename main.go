package main

import (
	"flag"
	"internal/api"
	"internal/database"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading '.env' file")
	}

	apiCfg := api.ApiConfig{
		JwtSecret:      os.Getenv("JWT_SECRET"),
		PolkaApiKey:    os.Getenv("POLKA_API_KEY"),
		FileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	mux.Handle("GET /app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", api.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MiddlewareMetricsCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.MiddlewareMetricsReset)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.ValidateChirp)

	corsMux := middlewareCors(mux)
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}

	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
