package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	service := service.NewService(inmemRepo)
	httpHandler := h.HttpHandler{Service: service}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /preview", httpHandler.TripPreviewHandler())

	httpAddr := ":8083"
	server := &http.Server{Addr: httpAddr, Handler: mux}
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
