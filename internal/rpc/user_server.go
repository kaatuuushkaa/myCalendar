package rpc

import (
	"context"
	"myCalendar/internal/rpc/get_user_by_login"
	"myCalendar/internal/rpc/reset_password"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "myCalendar/grpc/pb"
	"myCalendar/internal/rpc/auth"
	"myCalendar/internal/rpc/create_user"
	"myCalendar/internal/rpc/delete_user"
	"myCalendar/internal/rpc/get_user"
	"myCalendar/internal/rpc/health"
	"myCalendar/internal/rpc/update_user"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	health         *health.Handler
	createUser     *create_user.Handler
	auth           *auth.Handler
	getUser        *get_user.Handler
	updateUser     *update_user.Handler
	deleteUser     *delete_user.Handler
	resetPassword  *reset_password.Handler
	getUserByLogin *get_user_by_login.Handler
}

func NewUserServer(
	health *health.Handler,
	createUser *create_user.Handler,
	auth *auth.Handler,
	getUser *get_user.Handler,
	updateUser *update_user.Handler,
	deleteUser *delete_user.Handler,
	resetPassword *reset_password.Handler,
	getUserByLogin *get_user_by_login.Handler,
) *UserServer {
	return &UserServer{
		health:         health,
		createUser:     createUser,
		auth:           auth,
		getUser:        getUser,
		updateUser:     updateUser,
		deleteUser:     deleteUser,
		resetPassword:  resetPassword,
		getUserByLogin: getUserByLogin,
	}
}

func (s *UserServer) HealthCheck(ctx context.Context, e *emptypb.Empty) (*pb.HealthResponse, error) {
	return s.health.Handle(ctx, e)
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return s.createUser.Handle(ctx, req)
}

func (s *UserServer) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	return s.auth.Handle(ctx, req)
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return s.getUser.Handle(ctx, req)
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return s.updateUser.Handle(ctx, req)
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return s.deleteUser.Handle(ctx, req)
}

func (s *UserServer) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	return s.resetPassword.Handle(ctx, req)
}

func (s *UserServer) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.GetUserResponse, error) {
	return s.getUserByLogin.Handle(ctx, req)
}
