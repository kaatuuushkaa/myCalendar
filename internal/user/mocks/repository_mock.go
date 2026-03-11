package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"myCalendar/internal/user"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) GetByLogin(ctx context.Context, login string) (*user.User, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, username, email, name, surname, birth string) (*user.User, error) {
	args := m.Called(ctx, username, email, name, surname, birth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, username string) error {
	args := m.Called(ctx, username)
	return args.Error(0)
}
