package main

import (
	internal "Go-Server/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

func LogInMiddleware(cfg *apiConfig, handler func(w http.ResponseWriter, r *http.Request, userID uuid.UUID, tokenString string)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := internal.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Not allowed to perform this action", err)
			return
		}

		// Validate JWT
		userID, err := internal.ValidateJWT(tokenString, cfg.secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		// Call the handler with both userID and token
		handler(w, r, userID, tokenString)
	}
}
