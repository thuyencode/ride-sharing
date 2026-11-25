package util

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func ResponseInJSON(w http.ResponseWriter, code int, payload any) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("Failed to write HTTP response", "error", err)
	}
}

func ResponseErrorInJSON(w http.ResponseWriter, code int, message string, err error) {
	respBody := make(map[string]string, 2)
	respBody["error"] = message

	if code >= 500 {
		slog.Error("Internal server", "error", err)
	} else {
		if err != nil {
			respBody["details"] = err.Error()
		}
	}

	ResponseInJSON(w, code, respBody)
}
