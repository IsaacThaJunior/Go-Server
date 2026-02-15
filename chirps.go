package main

import (
	internal "Go-Server/internal/auth"
	"Go-Server/internal/database"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

type Chirps []Chirp

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body string `json:"body"`
	}

	getToken, err := internal.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userId, err := internal.ValidateJWT(getToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong", err)
		return
	}

	if len(params.Body) > lengthOfChirp {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	cleaned_body := breakWordsReplacement(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:     uuid.New(),
		Body:   cleaned_body,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Body:      chirp.Body,
		UserID:    chirp.UserID,
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}

func breakWordsReplacement(sentence string) string {
	words := strings.Split(sentence, " ")

	var cleaned []string

	for _, word := range words {
		lowered := strings.ToLower(word)
		if _, exists := badWords[lowered]; exists {
			cleaned = append(cleaned, "****")
		} else {
			cleaned = append(cleaned, word)
		}
	}

	return strings.Join(cleaned, " ")
}

func (cfg *apiConfig) handlerGetAllChips(w http.ResponseWriter, r *http.Request) {
	var response Chirps
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}
	for _, chirp := range chirps {
		response = append(response, Chirp{
			Body:      chirp.Body,
			UserID:    chirp.UserID,
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirp, err := cfg.getChirpFromDb(r)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp doesnt exist", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		Body:      chirp.Body,
		UserID:    chirp.UserID,
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request, userID uuid.UUID, tokenString string) {
	chirp, err := cfg.getChirpFromDb(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Unauthorized", errors.New("User can't perform this action as the user id doesnt match the chirp user id"))
		return
	}

	err = cfg.db.DeleteOneChirp(r.Context(), chirp.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "No chirp found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}

func (cfg *apiConfig) getChirpFromDb(r *http.Request) (Chirp, error) {
	vars := r.PathValue("chirpID")
	id, err := uuid.Parse(vars)
	if err != nil {
		return Chirp{}, err
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		return Chirp{}, err
	}
	return Chirp{
		ID:     chirp.ID,
		Body:   chirp.Body,
		UserID: chirp.UserID,
	}, nil
}

func (cfg *apiConfig) handlerUpdateChirpyRed(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	parsedId, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID", err)
		return
	}

	_, err = cfg.db.UpgradeChirpyStat(r.Context(), database.UpgradeChirpyStatParams{
		ID:          parsedId,
		IsChirpyRed: true,
	})

	if err != nil {
		respondWithError(w, http.StatusNotFound, "No user found", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
