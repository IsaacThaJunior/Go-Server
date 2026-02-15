package main

import (
	internal "Go-Server/internal/auth"
	"Go-Server/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const lengthOfChirp = 140

var badWords = map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}

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
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	log.Printf("Token: %s", getToken)

	userId, err := internal.ValidateJWT(getToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	log.Println("Here be the params user ID", userId)

	if len(params.Body) > lengthOfChirp {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	log.Println("Here be the params body", params.Body)

	cleaned_body := breakWordsReplacement(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:     uuid.New(),
		Body:   cleaned_body,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
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
		respondWithError(w, http.StatusInternalServerError, err.Error())
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
	vars := r.PathValue("chirpID")
	id, err := uuid.Parse(vars)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found")
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
