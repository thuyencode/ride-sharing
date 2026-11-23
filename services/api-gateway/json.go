package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type errorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func writeErrorInJSON(w http.ResponseWriter, status int, message string, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := errorResponse{Error: message, Details: err.Error()}

	if status == http.StatusInternalServerError {
		resp.Details = ""
		slog.Error("Internal error", "err", err)
	}

	return json.NewEncoder(w).Encode(resp)
}
