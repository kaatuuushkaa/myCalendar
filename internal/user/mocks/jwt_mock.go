package mocks

import (
	"github.com/stretchr/testify/mock"
	"myCalendar/internal/jwt"
	"net/http"
	"time"
)

type MockJWT struct {
	mock.Mock
}

func (m *MockJWT) GenerateJWT(id string, isValid bool, seconds int) string {
	args := m.Called(id, isValid, seconds)
	return args.String(0)
}

func (m *MockJWT) GenerateRefreshToken(id string, isValid bool, seconds int) (string, time.Time) {
	args := m.Called(id, isValid, seconds)
	t, ok := args.Get(1).(time.Time)
	if !ok {
		return args.String(0), time.Time{}
	}
	return args.String(0), t
}

func (m *MockJWT) GenerateTokenCookie(access, refresh string, expAfter time.Time) *http.Cookie {
	args := m.Called(access, refresh, expAfter)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*http.Cookie)
}

func (m *MockJWT) RefreshAccessToken(token string) (string, time.Time, error) {
	args := m.Called(token)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockJWT) ParseJWT(token string) (jwt.Claims, error) {
	args := m.Called(token)
	return args.Get(0).(jwt.Claims), args.Error(1)
}

func (m *MockJWT) SetRefreshTokenValidator(fn func(string) (bool, error)) {
	m.Called(fn)
}

func (m *MockJWT) SetInvalidateToken(fn func(string) (bool, error)) {
	m.Called(fn)
}

func (m *MockJWT) ValidateRefreshToken(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}

func (m *MockJWT) InvalidateRefreshToken(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}
