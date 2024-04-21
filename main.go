package main

import (
	"flag"
	"internal/database"
	"log"
	"net/http"
	"os"
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

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	mux.Handle("GET /app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.middlewareMetricsCount)
	mux.HandleFunc("GET /api/reset", apiCfg.middlewareMetricsReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.getChirpById)
	mux.HandleFunc("POST /api/users", apiCfg.postUser)
	mux.HandleFunc("GET /api/users", apiCfg.getUsers)
	mux.HandleFunc("GET /api/users/{id}", apiCfg.getUserById)
	mux.HandleFunc("POST /api/login", apiCfg.postLogin)

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
