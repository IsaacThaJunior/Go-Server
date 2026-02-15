package main

import (
	"Go-Server/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	dev            bool
	secret         string
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	environment := os.Getenv("PLATFORM")
	secret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("error connecting to db", err)
	}

	dbQueries := database.New(db)
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		dev:            environment == "dev",
		secret:         secret,
	}

	mux.Handle("/app/", http.StripPrefix("/app/", cfg.middlewareMetricsInc(fileServer)))
	mux.HandleFunc("GET /admin/metrics", cfg.readCounter)
	mux.HandleFunc("POST /admin/reset", cfg.resetCounter)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetAllChips)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", LogInMiddleware(cfg, cfg.handlerDeleteChirp))

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", LogInMiddleware(cfg, cfg.handlerEditUser))

	mux.HandleFunc("POST /api/login", cfg.handlerLoginUser)
	mux.HandleFunc("POST /api/refresh", cfg.handlerGetNewRefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevokeToken)

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
