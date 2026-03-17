package auth

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	pb "myCalendar/grpc/pb"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/jwt"
)

type Handler struct {
	repo       repo
	log        *zap.Logger
	jwtService jwt.IJWT
}

func New(repo repo, log *zap.Logger, jwtService jwt.IJWT) *Handler {
	return &Handler{repo: repo, log: log, jwtService: jwtService}
}

func (h *Handler) Handle(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	u, err := h.repo.GetByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidPassword
	}

	accessToken := h.jwtService.GenerateJWT(u.ID, true, jwt.Hour)
	refreshToken, _ := h.jwtService.GenerateRefreshToken(u.ID, true, jwt.Day*7)

	h.log.Info("user authenticated", zap.String("username", u.Username))

	return &pb.AuthResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
