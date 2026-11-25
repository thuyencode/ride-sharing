package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/util"
)

type HttpHandler struct {
	Service domain.TripService
}

func (s *HttpHandler) TripPreviewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var reqBody types.TripPreviewRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			util.ResponseErrorInJSON(w, http.StatusBadRequest, "JSON serializing failed", err)
			return
		}

		t, err := s.Service.GetRoute(r.Context(), &reqBody.Pickup, &reqBody.Destination)
		if err != nil {
			slog.Error("Failed to get route from trip service", "error", err)
			return
		}

		util.ResponseInJSON(w, http.StatusOK, t)
	}
}
