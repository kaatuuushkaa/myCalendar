package reset_password

import (
	"context"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"myCalendar/grpc/pb"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/ctxutil"
)

type Handler struct {
	repo repo
	log  *zap.Logger
}

func New(repo repo, log *zap.Logger) *Handler {
	return &Handler{
		repo: repo,
		log:  log,
	}
}

func (h *Handler) Handle(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	tokenUserID, err := ctxutil.UserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := h.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if tokenUserID != existing.ID {
		return nil, apperrors.ErrAccessDenied
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(req.OldPassword)); err != nil {
		return nil, apperrors.ErrInvalidPassword
	}
	if len(req.NewPassword) < 8 {
		return nil, apperrors.ErrInvalidLenPassword
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("failed to hash password", zap.Error(err))
		return nil, apperrors.ErrInternal
	}

	err = h.repo.ResetPassword(ctx, req.Username, string(newPasswordHash))
	if err != nil {
		return nil, err
	}

	h.log.Info("password reseted", zap.String("username", req.Username))

	return &pb.ResetPasswordResponse{Success: true}, nil

}
