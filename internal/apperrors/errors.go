package apperrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// user errors
var (
	ErrUserNotFound    = status.Error(codes.NotFound, "user not found")           //404
	ErrUserExists      = status.Error(codes.AlreadyExists, "user already exists") //409
	ErrInvalidPassword = status.Error(codes.Unauthenticated, "invalid password")  //401
	ErrAccessDenied    = status.Error(codes.PermissionDenied, "access denied")    //403
)

// event errors
var (
	ErrEventNotFound  = status.Error(codes.NotFound, "event not found")
	ErrInvalidTime    = status.Error(codes.InvalidArgument, "invalid time")
	ErrInvadArgument  = status.Error(codes.InvalidArgument, "event id is required")
	ErrEndBeforeStart = status.Error(codes.InvalidArgument, "end_at must be after start_at")
)

// auth errors
var (
	ErrInvalidToken = status.Error(codes.Unauthenticated, "invalid token")
	ErrTokenExpired = status.Error(codes.Unauthenticated, "token expired")
)

// common errors
var (
	ErrInternal = status.Error(codes.Internal, "internal server error") //500
)
