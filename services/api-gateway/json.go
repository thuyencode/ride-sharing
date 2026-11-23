package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"ride-sharing/shared/contracts"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func writeErrorInJSON(w http.ResponseWriter, status int, message string, err error) error {
	resp := contracts.APIResponse{
		Data:  err.Error(),
		Error: &contracts.APIError{Code: http.StatusText(status), Message: message},
	}

	if status == http.StatusInternalServerError {
		resp.Data = ""
		slog.Error("Internal error", "err", err)
	}

	return writeJSON(w, status, resp)
}
