package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "google.golang.org/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
)

const GrpcAddr = ":9093"

func main() {
	inmemRepo := repository.NewInmemRepository()
	_ = service.NewService(inmemRepo)

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		slog.Error("Failed to listen", "address", GrpcAddr, "error", err)
	}

	grpcServer := grpcserver.NewServer()
	serverErrCh := make(chan error, 1)

	go func() {
		slog.Info("gRPC server listening", "address", lis.Addr().String())
		serverErrCh <- grpcServer.Serve(lis)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrCh:
		slog.Error("gRPC server has something wrong", "error", err)
	case sig := <-signalCh:
		slog.Info("gRPC server is shutting down", "signal", sig.String())
		grpcServer.GracefulStop()
	}
}
