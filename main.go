package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.Handle("/app/", http.StripPrefix("/app/", cfg.middlewareMetricsInc(fileServer)))
	mux.HandleFunc("GET /admin/metrics", cfg.readCounter)
	mux.HandleFunc("POST /admin/reset", cfg.resetCounter)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerBody)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) readCounter(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)
	w.WriteHeader(http.StatusOK)

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  	</body>
		</html>
	`, hits)

	w.Write([]byte(html))

}
