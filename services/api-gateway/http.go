package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	"ride-sharing/shared/contracts"
)

func tripPreviewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody tripPreviewRequest

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		if reqBody.UserID == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		jsonBody, _ := json.Marshal(reqBody)
		reader := bytes.NewReader(jsonBody)

		resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
		if err != nil {
			slog.Error("Request failed", "url", "http://trip-service:8083/preview", "err", err)
			return
		}

		var respBody any
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			writeErrorInJSON(
				w,
				http.StatusBadRequest,
				"Failed to parse JSON data from trip service",
				err,
			)
			return
		}

		writeJSON(w, http.StatusOK, contracts.APIResponse{Data: respBody})
	}
}
