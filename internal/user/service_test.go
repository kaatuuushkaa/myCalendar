package user_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"myCalendar/grpc/pb"
	"myCalendar/internal/jwt"
	"myCalendar/internal/user"
	"myCalendar/internal/user/mocks"
	"testing"
	"time"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	//mockJWT := new(mocks.MockJWT)
	svc := user.NewService(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType(("*user.User"))).Return(nil)

	resp, err := svc.CreateUser(context.Background(), &pb.CreateUserRequest{
		Username: "stepa",
		Password: "ssecret",
		Email:    "stepa.com",
		Name:     "Stepa",
		Surname:  "Ivanov",
		Birth:    "2000-01-01",
	})

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotEmpty(t, resp.Id)
	mockRepo.AssertExpectations(t)
	//_ = mockJWT
}

func TestCreateUser_DBError(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	svc := user.NewService(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("db error"))

	resp, err := svc.CreateUser(context.Background(), &pb.CreateUserRequest{
		Username: "stepa",
		Password: "ssecret",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestAuth_Success(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockJWT := new(mocks.MockJWT)
	svc := user.NewService(mockRepo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("ssecret"), bcrypt.DefaultCost)

	mockRepo.On("GetByLogin", mock.Anything, "stepa").Return(&user.User{
		ID:       "uuid-stepa",
		Username: "stepa",
		Password: string(hash),
	}, nil)

	mockJWT.On("GenerateJWT", "uuid-stepa", true, jwt.Hour).Return("access-token")

	mockJWT.On("GenerateRefreshToken", "uuid-stepa", true, jwt.Day*7).Return("refresh-token", time.Now().Add(time.Hour*24*7))

	resp, err := svc.Auth(context.Background(), &pb.AuthRequest{
		Login:    "stepa",
		Password: "ssecret",
	}, mockJWT)

	assert.NoError(t, err)
	assert.True(t, resp.Success)

	assert.Equal(t, "access-token", resp.AccessToken)
	assert.Equal(t, "refresh-token", resp.RefreshToken)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestAuth_WrongPassword(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockJWT := new(mocks.MockJWT)
	svc := user.NewService(mockRepo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)

	mockRepo.On("GetByLogin", mock.Anything, "stepa").Return(&user.User{
		ID:       "uuid-stepa",
		Username: "stepa",
		Password: string(hash),
	}, nil)

	resp, err := svc.Auth(context.Background(), &pb.AuthRequest{
		Login:    "stepa",
		Password: "wrong",
	}, mockJWT)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockJWT.AssertNotCalled(t, "GenerateJWT")
	mockRepo.AssertExpectations(t)
}

func TestAuth_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockJWT := new(mocks.MockJWT)
	svc := user.NewService(mockRepo)

	mockRepo.On("GetByLogin", mock.Anything, "ghost").Return(nil, errors.New("user not found"))

	resp, err := svc.Auth(context.Background(), &pb.AuthRequest{
		Login:    "ghost",
		Password: "secret",
	}, mockJWT)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockJWT.AssertNotCalled(t, "GenerateJWT")
	mockRepo.AssertExpectations(t)
}
