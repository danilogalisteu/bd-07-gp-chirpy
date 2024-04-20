package main

import (
	"fmt"
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
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("GET /app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.middlewareMetricsCount)
	mux.HandleFunc("GET /api/reset", apiCfg.middlewareMetricsReset)
	mux.HandleFunc("POST /api/chirps", postChirp)
	mux.HandleFunc("GET /api/chirps", getChirps)

	corsMux := middlewareCors(mux)
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
