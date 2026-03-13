package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	pb "myCalendar/grpc/pb"
	"myCalendar/internal/jwt"
)

type Service struct {
	repo       RepositoryInterface
	log        *zap.Logger
	jwtService jwt.IJWT
}

func NewService(repo RepositoryInterface, log *zap.Logger, jwtService jwt.IJWT) *Service {
	return &Service{repo: repo, log: log, jwtService: jwtService}
}

func (s *Service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Password: string(hash),
		Email:    req.Email,
		Name:     req.Name,
		Surname:  req.Surname,
		Birth:    req.Birth,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &pb.CreateUserResponse{Success: true, Id: u.ID}, nil
}

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	u, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{User: modelToProto(u)}, nil
}

func (s *Service) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	u, err := s.repo.Update(ctx, req.Username, req.Email, req.Name, req.Surname, req.Birth)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{Success: true, User: modelToProto(u)}, nil
}

func (s *Service) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := s.repo.Delete(ctx, req.Username); err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}

func (s *Service) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	u, err := s.repo.GetByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(([]byte(u.Password)), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("Incorrect password")
	}

	accessToken := s.jwtService.GenerateJWT(u.ID, true, jwt.Hour)
	refreshToken, _ := s.jwtService.GenerateRefreshToken(u.ID, true, jwt.Day*7)

	return &pb.AuthResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func modelToProto(u *User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Surname:  u.Surname,
		Birth:    u.Birth,
	}
}
