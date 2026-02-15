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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ChirpyRed    bool      `json:"is_chirpy_red"`
}
type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var userBody UserReq
	err := decoder.Decode(&userBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	hashedPassword, err := internal.HashPassword(userBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		Email:          userBody.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	token, err := internal.MakeJWT(user.ID, cfg.secret, tokenTime)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
		ChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerEditUser(w http.ResponseWriter, r *http.Request, userId uuid.UUID, tokenString string) {

	decoder := json.NewDecoder(r.Body)

	// Get and decode body from request
	var userBody UserReq
	err := decoder.Decode(&userBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// hash password and then update the user password and email
	hashedPassword, err := internal.HashPassword(userBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          userBody.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	// Send back user if everything works out
	respondWithJSON(w, http.StatusOK, User{
		ID:        userId,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
		Token:     tokenString,
		ChirpyRed: updatedUser.IsChirpyRed,
	})
}
