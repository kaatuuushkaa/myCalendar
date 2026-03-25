package middleware

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/ctxutil"
	"myCalendar/internal/jwt"
	"strings"
)

var publicMethods = map[string]bool{
	"/userGRPC.UserService/HealthCheck": true,
	"/userGRPC.UserService/CreateUser":  true,
	"/userGRPC.UserService/Auth":        true,
}

func AuthInterceptor(jwtService jwt.IJWT) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		token, err := extractToken(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		claims, err := jwtService.ParseJWT(token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return nil, status.Error(codes.Unauthenticated, "Token expired")
			}
			return nil, apperrors.ErrInvalidToken
		}

		if strings.TrimSpace(claims.ID) == "" {
			return nil, apperrors.ErrInvalidTokenWithoutID
		}

		if !claims.IsValid {
			return nil, apperrors.ErrIsValidFalse
		}

		ctx = ctxutil.NewContextWithUserID(ctx, claims.ID)
		return handler(ctx, req)
	}
}

func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("Metadata is unavailable")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", errors.New("Authorization header missing")
	}

	parts := strings.SplitN(values[0], " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Incorrect token format")
	}

	return parts[1], nil
}
