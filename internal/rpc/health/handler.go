package health

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "myCalendar/grpc/pb"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(_ context.Context, _ *emptypb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Status: "ok"}, nil
}
