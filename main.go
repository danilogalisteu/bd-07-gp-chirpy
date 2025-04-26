package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/danilogalisteu/bd-07-gp-chirpy/internal/api"
	"github.com/danilogalisteu/bd-07-gp-chirpy/internal/database"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
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

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	apiCfg := api.ApiConfig{
		JwtSecret:      os.Getenv("JWT_SECRET"),
		PolkaApiKey:    os.Getenv("POLKA_API_KEY"),
		FileserverHits: 0,
		DbQueries:      dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("GET /app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", api.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MiddlewareMetricsCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.MiddlewareMetricsReset)
	mux.HandleFunc("POST /api/users", apiCfg.CreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.UpdateUser)
	mux.HandleFunc("POST /api/login", apiCfg.GetUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.GetToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.UpdateToken)
	mux.HandleFunc("POST /api/chirps", apiCfg.CreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirp)

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
