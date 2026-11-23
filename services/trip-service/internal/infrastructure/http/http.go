package http

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
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
			if err := writeErrorInJSON(
				w,
				http.StatusBadRequest,
				"JSON serializing failed",
				err,
			); err != nil {
				slog.Error("Failed to write response", "err", err)
			}
			return
		}

		t, err := s.Service.GetRoute(r.Context(), &reqBody.Pickup, &reqBody.Destination)
		if err != nil {
			log.Println(err)
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
