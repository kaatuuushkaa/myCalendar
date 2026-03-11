package user_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"myCalendar/grpc/pb"
	"myCalendar/internal/user"
	"myCalendar/internal/user/mocks"
	"testing"
)

func ctxWithUserID(id string) context.Context {
	return context.WithValue(context.Background(), "user_id", id)
}

func TestGetUser_OwnerAccess(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockJWT := new(mocks.MockJWT)
	svc := user.NewService(mockRepo)
	h := user.NewHandler(svc, mockJWT)

	mockRepo.On("GetByUsername", mock.Anything, "stepa").Return(&user.User{
		ID:       "uuid-stepa",
		Username: "stepa",
		Name:     "Stepa",
	}, nil)

	resp, err := h.GetUser(ctxWithUserID("uuid-stepa"), &pb.GetUserRequest{
		Username: "stepa",
	})

	assert.NoError(t, err)
	assert.Equal(t, "stepa", resp.User.Username)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_AccessDenied(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockJWT := new(mocks.MockJWT)
	svc := user.NewService(mockRepo)
	h := user.NewHandler(svc, mockJWT)

	mockRepo.On("GetByUsername", mock.Anything, "nikita").Return(&user.User{
		ID:       "uuid-nikita",
		Username: "nikita",
	}, nil)

	resp, err := h.GetUser(ctxWithUserID("uuid-stepa"), &pb.GetUserRequest{
		Username: "nikita",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.PermissionDenied, st.Code())
}

func TestDeleteUser_OwnerAccess(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockJWT := new(mocks.MockJWT)
	svc := user.NewService(mockRepo)
	h := user.NewHandler(svc, mockJWT)

	mockRepo.On("GetByUsername", mock.Anything, "stepa").Return(&user.User{
		ID:       "uuid-stepa",
		Username: "stepa",
	}, nil)

	mockRepo.On("Delete", mock.Anything, "stepa").Return(nil)

	resp, err := h.DeleteUser(ctxWithUserID("uuid-stepa"), &pb.DeleteUserRequest{
		Username: "stepa",
	})

	assert.NoError(t, err)
	assert.True(t, resp.Success)

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_AccessDenied(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockJWT := new(mocks.MockJWT)
	svc := user.NewService(mockRepo)
	h := user.NewHandler(svc, mockJWT)

	mockRepo.On("GetByUsername", mock.Anything, "nikita").Return(&user.User{
		ID:       "uuid-nikita",
		Username: "nikita",
	}, nil)

	resp, err := h.DeleteUser(ctxWithUserID("uuid-stepa"), &pb.DeleteUserRequest{
		Username: "nikita",
	})

	assert.Nil(t, resp)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.PermissionDenied, st.Code())

	mockRepo.AssertNotCalled(t, "Delete")
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockJWT := new(mocks.MockJWT)
	svc := user.NewService(mockRepo)
	h := user.NewHandler(svc, mockJWT)

	// репозиторий вернёт ошибку — такого юзера нет
	mockRepo.On("GetByUsername", mock.Anything, "ghost").
		Return(nil, errors.New("user not found"))

	resp, err := h.UpdateUser(ctxWithUserID("uuid-stepa"), &pb.UpdateUserRequest{
		Username: "ghost",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	// Update не должен вызываться если юзер не найден
	mockRepo.AssertNotCalled(t, "Update")
}
