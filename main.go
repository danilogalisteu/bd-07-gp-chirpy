package main

import (
	"flag"
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

	apiCfg := apiConfig{}

	fname := "database.json"
	if *dbg {
		err := os.Remove(fname)
		if err != nil {
			log.Printf("Error removing DB file %s:\n%v", fname, err)
		}
	}
	db, err := NewDB(fname)
	if err != nil {
		log.Printf("Error creating DB with file %s:\n%v", fname, err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.middlewareMetricsCount)
	mux.HandleFunc("GET /api/reset", apiCfg.middlewareMetricsReset)
	mux.HandleFunc("POST /api/chirps", db.postChirp)
	mux.HandleFunc("GET /api/chirps", db.getChirps)
	mux.HandleFunc("GET /api/chirps/{id}", db.getChirpById)
	mux.HandleFunc("POST /api/users", db.postUser)

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
