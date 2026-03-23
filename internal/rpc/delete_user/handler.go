package delete_user

import (
	"context"

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

func (h *Handler) Handle(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	tokenUserID, err := ctxutil.UserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	// проверяем что пользователь существует
	existing, err := h.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// проверяем ownership — нельзя удалить чужой аккаунт
	if tokenUserID != existing.ID {
		return nil, apperrors.ErrAccessDenied
	}

	if err := h.repo.Delete(ctx, req.Username); err != nil {
		return nil, err
	}

	h.log.Info("user deleted", zap.String("username", req.Username))

	return &pb.DeleteUserResponse{Success: true}, nil
}
