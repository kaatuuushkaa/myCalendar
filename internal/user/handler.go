package user

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "myCalendar/grpc/pb"
)

type Handler struct {
	pb.UnimplementedUserServiceServer
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Status: "ok"}, nil
}

func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return h.service.CreateUser(ctx, req)
}

func (h *Handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return h.service.GetUser(ctx, req)
}

func (h *Handler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return h.service.UpdateUser(ctx, req)
}

func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return h.service.DeleteUser(ctx, req)
}
