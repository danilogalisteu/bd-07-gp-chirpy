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
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	fname := "database.json"
	if *dbg {
		if _, err := os.Stat(fname); os.IsNotExist(err) {
			log.Printf("DB file '%s' doesn't exist:\n%v", fname, err)
		} else {
			err := os.Remove(fname)
			if err != nil {
				log.Printf("Error removing DB file '%s':\n%v", fname, err)
			}
		}
	}
	db, err := database.NewDB(fname)
	if err != nil {
		log.Printf("Error creating DB with file '%s':\n%v", fname, err)
	}

	err = godotenv.Load()
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
	mux.Handle("GET /app/*", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", api.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MiddlewareMetricsCount)
	mux.HandleFunc("GET /api/reset", apiCfg.MiddlewareMetricsReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.PostChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.GetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCfg.DeleteChirpById)
	mux.HandleFunc("POST /api/users", apiCfg.PostUser)
	mux.HandleFunc("GET /api/users", apiCfg.GetUsers)
	mux.HandleFunc("PUT /api/users", apiCfg.PutUser)
	mux.HandleFunc("GET /api/users/{id}", apiCfg.GetUserById)
	mux.HandleFunc("POST /api/login", apiCfg.PostLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.PostRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.PostRevoke)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.PostPolkaWebhooks)

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
