package main

import (
	internal "Go-Server/internal/auth"
	"Go-Server/internal/database"
	"encoding/json"
	"net/http"
	"time"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	var userReq UserReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), userReq.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	checkPassword, err := internal.CheckPasswordHash(userReq.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	if !checkPassword {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	token, err := internal.MakeJWT(user.ID, cfg.secret, tokenTime)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	getRefreshToken, err := internal.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
	}

	storeToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     getRefreshToken,
		ExpiresAt: time.Now().Add(RefreshTokenExipires),
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: storeToken.Token,
		ChirpyRed:    user.IsChirpyRed,
	})
}
