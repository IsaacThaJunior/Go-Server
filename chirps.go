package main

import (
	"Go-Server/internal/database"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const lengthOfChirp = 140

var badWords = map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(params.Body) > lengthOfChirp {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned_body := breakWordsReplacement(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:     uuid.New(),
		Body:   cleaned_body,
		UserID: params.UserID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, parameter{
		Body:   chirp.Body,
		UserID: chirp.UserID,
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
