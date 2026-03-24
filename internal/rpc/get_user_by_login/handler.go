package get_user_by_login

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "myCalendar/grpc/pb"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/ctxutil"
)

type Handler struct {
	repo repo
	log  *zap.Logger
}

func New(repo repo, log *zap.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}

func (h *Handler) Handle(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.GetUserResponse, error) {
	if req.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	// проверяем что запрашивает авторизованный пользователь
	tokenUserID, err := ctxutil.UserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	u, err := h.repo.GetByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	// проверяем что пользователь запрашивает свои данные
	if tokenUserID != u.ID {
		return nil, apperrors.ErrAccessDenied
	}

	return &pb.GetUserResponse{User: &pb.UserResponse{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Surname:  u.Surname,
		Birth:    u.Birth,
	}}, nil
}
