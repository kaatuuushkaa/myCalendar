package user

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "myCalendar/grpc/pb"
	"myCalendar/internal/jwt"
)

type Handler struct {
	pb.UnimplementedUserServiceServer
	service    *Service
	jwtService jwt.IJWT
}

func NewHandler(service *Service, jwtService jwt.IJWT) *Handler {
	return &Handler{service: service, jwtService: jwtService}
}

func (h *Handler) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Status: "ok"}, nil
}

func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return h.service.CreateUser(ctx, req)
}

func (h *Handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if err := h.checkOwner(ctx, req.Username); err != nil {
		return nil, err
	}
	return h.service.GetUser(ctx, req)
}

func (h *Handler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if err := h.checkOwner(ctx, req.Username); err != nil {
		return nil, err
	}
	return h.service.UpdateUser(ctx, req)
}

func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := h.checkOwner(ctx, req.Username); err != nil {
		return nil, err
	}
	return h.service.DeleteUser(ctx, req)
}

func (h *Handler) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	return h.service.Auth(ctx, req, h.jwtService)
}

func userIDFromTokenCtx(ctx context.Context) (string, error) {
	val := ctx.Value("user_id")
	if val == nil {
		return "", status.Error(codes.Unauthenticated, "user_id not found in context")
	}
	id, ok := val.(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "invalid user_id in context")
	}
	return id, nil
}

func (h *Handler) checkOwner(ctx context.Context, username string) error {
	tokenUserID, err := userIDFromTokenCtx(ctx)
	if err != nil {
		return err
	}

	existing, err := h.service.GetUser(ctx, &pb.GetUserRequest{Username: username})
	if err != nil {
		return err
	}

	if tokenUserID != existing.User.Id {
		return status.Error(codes.PermissionDenied, "Access denied")
	}

	return nil
}
