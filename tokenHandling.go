package main

import (
	internal "Go-Server/internal/auth"
	"database/sql"
	"errors"
	"net/http"
)

func (cfg *apiConfig) handlerGetNewRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := internal.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized", err)
		return
	}

	token, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	getAccessToken, err := internal.MakeJWT(token.UserID, cfg.secret, tokenTime)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	respondWithJSON(w, http.StatusOK, TokenResponse{
		Token: getAccessToken,
	})

}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := internal.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
