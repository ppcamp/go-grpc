package grpc

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type grpcHealthService struct {
	*grpc_health_v1.UnimplementedHealthServer
}

func NewHealthService() grpc_health_v1.HealthServer {
	return new(grpcHealthService)
}

func (m *grpcHealthService) Check(
	_ context.Context,
	_ *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}
