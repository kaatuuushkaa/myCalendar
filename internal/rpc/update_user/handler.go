package update_user

import (
	"context"
	"myCalendar/internal/domain"
	"strings"

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

func (h *Handler) Handle(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if err := validate(req); err != nil {
		return nil, err
	}

	tokenUserID, err := ctxutil.UserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	// сначала проверяем что пользователь существует
	existing, err := h.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// проверяем ownership
	if tokenUserID != existing.ID {
		return nil, apperrors.ErrAccessDenied
	}

	// обновляем
	updated, err := h.repo.Update(ctx, req.Username, req.Email, req.Name, req.Surname, req.Birth)
	if err != nil {
		return nil, err
	}

	h.log.Info("user updated", zap.String("username", req.Username))

	return &pb.UpdateUserResponse{Success: true, User: toProto(updated)}, nil
}

func validate(req *pb.UpdateUserRequest) error {
	if req.Username == "" {
		return apperrors.ErrEmptyUsername
	}
	if !strings.Contains(req.Email, "@") {
		return apperrors.ErrInvalidEmail
	}
	return nil
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
