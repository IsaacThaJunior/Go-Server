package main

import (
	"errors"
	"net/http"
)

func (cfg *apiConfig) resetCounter(w http.ResponseWriter, r *http.Request) {
	if !cfg.dev {
		respondWithError(w, http.StatusForbidden, "Forbidden Outside Dev", errors.New("You cant perform this action unless in Dev"))
	}
	err := cfg.db.DeleteAllUser(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, InternalServerMessage, err)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
