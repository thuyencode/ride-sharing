package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("Failed to write HTTP response", "error", err)
	}
}

func writeErrorInJSON(w http.ResponseWriter, code int, message string, err error) {
	if code >= 500 {
		slog.Error("Internal server", "error", err)
	}
	writeJSON(w, code, map[string]string{"error": message})
}
