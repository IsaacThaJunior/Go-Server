package main

import (
	internal "Go-Server/internal/auth"
	"Go-Server/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}
type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Expires  int64  `json:"expires_in_seconds"`
}

const tokenTime = 1 * time.Hour

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var userBody UserReq
	err := decoder.Decode(&userBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := internal.HashPassword(userBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		Email:          userBody.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	timeThatWillBeUsed := tokenTime

	if userBody.Expires > 0 { // or whatever condition you mean
		timeThatWillBeUsed = time.Duration(userBody.Expires) * time.Minute
	}
	token, err := internal.MakeJWT(user.ID, cfg.secret, timeThatWillBeUsed)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	var userReq UserReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := cfg.db.GetUser(r.Context(), userReq.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	checkPassword, err := internal.CheckPasswordHash(userReq.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !checkPassword {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	timeThatWillBeUsed := tokenTime

	if userReq.Expires > 0 { // or whatever condition you mean
		timeThatWillBeUsed = time.Duration(userReq.Expires) * time.Minute
	}
	token, err := internal.MakeJWT(user.ID, cfg.secret, timeThatWillBeUsed)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
