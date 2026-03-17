package ctxutil

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// типизированный ключ — string напрямую антипаттерн
type contextKey string

const userIDKey contextKey = "user_id"

func UserIDFromCtx(ctx context.Context) (string, error) {
	val := ctx.Value(userIDKey)
	if val == nil {
		return "", status.Error(codes.Unauthenticated, "user_id not found in context")
	}
	id, ok := val.(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "invalid user_id in context")
	}
	return id, nil
}

func NewContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
