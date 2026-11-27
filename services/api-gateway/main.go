package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ride-sharing/shared/env"
)

var httpAddr = env.GetString("HTTP_ADDR", ":8081")

func main() {
	slog.Info("Starting API Gateway")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /trip/preview", tripPreviewHandler())
	mux.HandleFunc("/ws/drivers", driversWebsocketHandler())
	mux.HandleFunc("/ws/riders", ridersWebsocketHandler())

	handler := applyMiddleware(mux, prelight("*"))
	server := &http.Server{Addr: httpAddr, Handler: handler}
	serverErrCh := make(chan error, 1)

	go func() {
		slog.Info("HTTP server listening", "address", httpAddr)
		serverErrCh <- server.ListenAndServe()
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrCh:
		slog.Error("Failed to start HTTP server", "error", err)
	case sig := <-shutdownCh:
		slog.Info("HTTP server is shutting down", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown HTTP server gracefully", "error", err)
			if err := server.Close(); err != nil {
				slog.Error("Failed to close HTTP server immediately", "error", err)
			}
		}
	}
}
