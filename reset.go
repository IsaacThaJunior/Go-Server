package main

import "net/http"

func (cfg *apiConfig) resetCounter(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
