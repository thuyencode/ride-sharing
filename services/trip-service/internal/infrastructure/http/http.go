package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"
)

type HttpHandler struct {
	Service domain.TripService
}

type tripPreviewRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (s *HttpHandler) TripPreviewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var reqBody tripPreviewRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			if err := writeErrorInJSON(w, http.StatusBadRequest, "JSON serializing failed", err); err != nil {
				slog.Error("Failed to write response", "err", err)
			}
			return
		}

		fare := &domain.RideFareModel{UserID: "67"}
		t, err := s.Service.CreateTrip(r.Context(), fare)
		if err != nil {
			slog.Error("Failed to create trip", "err", err)
			return
		}

		writeJSON(w, http.StatusOK, t)
	}
}

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
		slog.Error("Internal error", "err", err)
		resp.Details = ""
	}

	return json.NewEncoder(w).Encode(resp)
}
