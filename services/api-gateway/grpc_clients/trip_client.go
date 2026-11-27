package grpcclients

import (
	"log/slog"

	"google.golang.org/grpc"
	"ride-sharing/shared/env"
	"ride-sharing/shared/proto/trip"
)

type tripServiceClient struct {
	Client trip.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripServiceClient() (*tripServiceClient, error) {
	tripServiceURL := env.GetString("TRIP_SERVICE_URL", "trip-service:9093")
	conn, err := grpc.NewClient(tripServiceURL)
	if err != nil {
		return nil, err
	}

	return &tripServiceClient{trip.NewTripServiceClient(conn), conn}, nil
}

func (c *tripServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			slog.Error("Failed to close gRPC connection", "error", err)
			return
		}
	}
}
