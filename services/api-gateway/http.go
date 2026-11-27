package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	grpcclients "ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"
)

func tripPreviewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody types.TripPreviewRequest

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			util.ResponseErrorInJSON(
				w,
				http.StatusBadRequest,
				"Failed to deserialize JSON body",
				err,
			)
			return
		}

		defer r.Body.Close()

		if reqBody.UserID == "" {
			util.ResponseErrorInJSON(w, http.StatusBadRequest, "user id must not be empty", nil)
			return
		}

		jsonBody, _ := json.Marshal(reqBody)
		reader := bytes.NewReader(jsonBody)

		tripService, err := grpcclients.NewTripServiceClient()
		if err != nil {
			util.ResponseErrorInJSON(
				w,
				http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				err,
			)
			return
		}

		defer tripService.Close()

		resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
		if err != nil {
			slog.Error("Request failed", "url", "http://trip-service:8083/preview", "err", err)
			return
		}

		var respBody any
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			util.ResponseErrorInJSON(
				w,
				http.StatusBadRequest,
				"Failed to parse JSON data from trip service",
				err,
			)
			return
		}

		util.ResponseInJSON(w, http.StatusOK, contracts.APIResponse{Data: respBody})
	}
}
