package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthCheckSrv struct {
	Status grpc_health_v1.HealthCheckResponse_ServingStatus
	Reason string
}

func (h *HealthCheckSrv) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

func (h *HealthCheckSrv) OffLine(reason string) {
	h.Status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	h.Reason = reason
	fmt.Println(reason)
}
func (h *HealthCheckSrv) OnLine(reason string) {
	h.Status = grpc_health_v1.HealthCheckResponse_SERVING
	h.Reason = reason
	fmt.Println(reason)
}

func (h *HealthCheckSrv) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: h.Status,
	}, nil
}
