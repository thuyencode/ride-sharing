package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, _ *http.Request, status int, reason error) {
		util.ResponseErrorInJSON(w, status, "Failed to upgrade to WebSocket", reason)
	},
}

func driversWebsocketHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		defer func() {
			if err := conn.Close(); err != nil {
				slog.Error("Failed to close WebSocket connection", "error", err)
			}
		}()

		userID := r.URL.Query().Get("userID")
		if userID == "" {
			if err := conn.WriteJSON(
				map[string]string{"error": "No user id provided"},
			); err != nil {
				slog.Error("Failed to write WebSocket response", "error", err)
			}
			return
		}

		packageSlug := r.URL.Query().Get("packageSlug")
		if packageSlug == "" {
			if err := conn.WriteJSON(
				map[string]string{"error": "No package slug provided"},
			); err != nil {
				slog.Error("Failed to write WebSocket response", "error", err)
			}
			return
		}

		type Driver struct {
			Id             string `json:"id"`
			Name           string `json:"name"`
			ProfilePicture string `json:"profilePicture"`
			CarPlate       string `json:"carPlate"`
			PackageSluq    string `json:"packageSluq"`
		}

		msg := contracts.WSMessage{
			Type: contracts.DriverCmdRegister,
			Data: Driver{
				Id:             userID,
				Name:           "67",
				ProfilePicture: util.GetRandomAvatar(1),
				CarPlate:       "A113",
				PackageSluq:    packageSlug,
			},
		}

		if err := conn.WriteJSON(msg); err != nil {
			slog.Error("Failed to send message through WebSocket", "error", err)
			return
		}

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				slog.Error("Failed to reading messages from WebSocket connection", "error", err)
				break
			}

			slog.Info("Received message from Websocket connection", "message", message)
		}
	}
}

func ridersWebsocketHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		defer func() {
			if err := conn.Close(); err != nil {
				slog.Error("Failed to close WebSocket connection", "error", err)
			}
		}()

		userID := r.URL.Query().Get("userID")
		if userID == "" {
			if err := conn.WriteJSON(
				map[string]string{"error": "No user id provided"},
			); err != nil {
				slog.Error("Failed to write WebSocket response", "error", err)
			}
			return
		}

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				slog.Error("Failed to reading messages from WebSocket connection", "error", err)
				break
			}

			slog.Info("Received message from Websocket connection", "message", message)
		}
	}
}
