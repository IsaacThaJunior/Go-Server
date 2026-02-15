package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	// Get caller info
	pc, _, line, ok := runtime.Caller(1) // caller of respondWithError
	funcName := "unknown"
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			// Strip package path, keep only function name
			fullName := fn.Name()
			parts := strings.Split(fullName, ".")
			funcName = parts[len(parts)-1]
		}
	}

	// Log concise error
	log.Printf("Error in %s:%d: %v\n", funcName, line, err)

	// Respond to client
	respondWithJSON(w, code, errorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
