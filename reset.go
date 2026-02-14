package main

import "net/http"

func (cfg *apiConfig) resetCounter(w http.ResponseWriter, r *http.Request) {
	if !cfg.dev {
		respondWithError(w, http.StatusForbidden, "Forbidden Outside Dev")
	}
	err := cfg.db.DeleteAllUser(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
