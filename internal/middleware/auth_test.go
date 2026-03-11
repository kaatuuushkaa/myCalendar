package middleware_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"myCalendar/internal/jwt"
	"myCalendar/internal/middleware"
	"myCalendar/internal/user/mocks"
	"testing"
)

func ctxWithToken(token string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewIncomingContext(context.Background(), md)
}

func TestAuthInterceptor_PublicMethod(t *testing.T) {
	mockJWT := new(mocks.MockJWT)
	interceptor := middleware.AuthInterceptor(mockJWT)

	called := false
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		called = true
		return "ok", nil
	}

	_, err := interceptor(
		context.Background(), // без токена — публичный метод не требует
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/userGRPC.UserService/Auth"},
		handler,
	)

	assert.NoError(t, err)
	assert.True(t, called)
	mockJWT.AssertNotCalled(t, "ParseJWT")
}

func TestAuthInterceptor_ValidToken(t *testing.T) {
	mockJWT := new(mocks.MockJWT)
	interceptor := middleware.AuthInterceptor(mockJWT)

	mockJWT.On("ParseJWT", "valid-token").
		Return(jwt.Claims{ID: "uuid-alice"}, nil)

	var receivedID string
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		receivedID = ctx.Value("user_id").(string)
		return "ok", nil
	}

	_, err := interceptor(
		ctxWithToken("valid-token"),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/userGRPC.UserService/GetUser"},
		handler,
	)

	assert.NoError(t, err)
	assert.Equal(t, "uuid-alice", receivedID)
	mockJWT.AssertExpectations(t)
}

func TestAuthInterceptor_NoToken(t *testing.T) {
	mockJWT := new(mocks.MockJWT)
	interceptor := middleware.AuthInterceptor(mockJWT)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	_, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/userGRPC.UserService/GetUser"},
		handler,
	)

	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	mockJWT.AssertNotCalled(t, "ParseJWT")
}

func TestAuthInterceptor_ExpiredToken(t *testing.T) {
	mockJWT := new(mocks.MockJWT)
	interceptor := middleware.AuthInterceptor(mockJWT)

	// мок имитирует истёкший токен
	mockJWT.On("ParseJWT", "expired-token").
		Return(jwt.Claims{}, jwt.ErrTokenExpired)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	_, err := interceptor(
		ctxWithToken("expired-token"),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/userGRPC.UserService/GetUser"},
		handler,
	)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "expired")
	mockJWT.AssertExpectations(t)
}
