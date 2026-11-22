package main

import (
	"encoding/json"
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

		// TODO: Call trip service

		writeJSON(
			w,
			http.StatusCreated,
			contracts.APIResponse{Data: http.StatusText(http.StatusCreated)},
		)
	}
}
