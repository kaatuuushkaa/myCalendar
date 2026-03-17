package create_user

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "myCalendar/grpc/pb"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/domain"
)

type Handler struct {
	repo repo
	log  *zap.Logger
}

func New(repo repo, log *zap.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}

func (h *Handler) Handle(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := validate(req); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("failed to hash password", zap.Error(err))
		return nil, apperrors.ErrInternal
	}

	u := domain.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Password: string(hash),
		Email:    req.Email,
		Name:     req.Name,
		Surname:  req.Surname,
		Birth:    req.Birth,
	}

	if err := h.repo.Create(ctx, u); err != nil {
		h.log.Warn("failed to create user",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		return nil, err
	}

	h.log.Info("user created",
		zap.String("id", u.ID),
		zap.String("username", u.Username),
	)

	return &pb.CreateUserResponse{Success: true, Id: u.ID}, nil
}

func validate(req *pb.CreateUserRequest) error {
	if req.Username == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}
	if len(req.Password) < 8 {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}
	if !strings.Contains(req.Email, "@") {
		return status.Error(codes.InvalidArgument, "invalid email")
	}
	return nil
}
