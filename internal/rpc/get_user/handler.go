package get_user

import (
	"context"
	"myCalendar/internal/domain"

	"go.uber.org/zap"

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

func (h *Handler) Handle(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// достаём user_id из контекста — положил туда middleware после проверки токена
	tokenUserID, err := ctxutil.UserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	// ищем пользователя в БД
	u, err := h.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// проверяем что токен принадлежит тому же пользователю
	// нельзя смотреть чужой профиль
	if tokenUserID != u.ID {
		return nil, apperrors.ErrAccessDenied
	}

	return &pb.GetUserResponse{User: toProto(u)}, nil
}

func toProto(u domain.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Surname:  u.Surname,
		Birth:    u.Birth,
	}
}
