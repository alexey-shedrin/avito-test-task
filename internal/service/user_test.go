package service_test

import (
	"errors"
	"testing"

	"github.com/alexey-shedrin/avito-test-task/internal/repository/mocks"

	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/request"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/alexey-shedrin/avito-test-task/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserService_DummyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := service.NewUserService(mockRepo)

	t.Run("should return token when role is valid", func(t *testing.T) {
		resp, err := userService.DummyLogin(&request.DummyLogin{Role: "user"})

		require.NoError(t, err)
		require.NotEmpty(t, resp.Token)
	})

	t.Run("should return token when role is admin", func(t *testing.T) {
		resp, err := userService.DummyLogin(&request.DummyLogin{Role: "admin"})

		require.NoError(t, err)
		require.NotEmpty(t, resp.Token)
	})
}

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := service.NewUserService(mockRepo)

	t.Run("should register user successfully", func(t *testing.T) {
		req := &request.Register{
			Email:    "test@example.com",
			Password: "password123",
			Role:     "user",
		}

		mockRepo.EXPECT().
			Create(gomock.Any()).
			DoAndReturn(func(user *entity.User) error {
				require.Equal(t, req.Email, user.Email)
				require.Equal(t, req.Role, user.Role)
				require.NotEqual(t, req.Password, user.Password)
				return nil
			}).
			Times(1)

		user, err := userService.Register(req)

		require.NoError(t, err)
		require.Equal(t, req.Email, user.Email)
		require.Equal(t, req.Role, user.Role)
		require.NotEqual(t, req.Password, user.Password)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		req := &request.Register{
			Email:    "test@example.com",
			Password: "password123",
			Role:     "user",
		}

		mockRepo.EXPECT().
			Create(gomock.Any()).
			Return(errors.New("db error")).
			Times(1)

		user, err := userService.Register(req)

		require.Error(t, err)
		require.Nil(t, user)
	})
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := service.NewUserService(mockRepo)

	t.Run("should login successfully with valid credentials", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"
		role := "user"

		user := &entity.User{
			Id:       uuid.New(),
			Email:    email,
			Password: password,
			Role:     role,
		}

		user.HashPassword()

		mockRepo.EXPECT().
			GetByEmail(email).
			Return(user, nil).
			Times(1)

		resp, err := userService.Login(&request.Login{
			Email:    email,
			Password: password,
		})

		require.NoError(t, err)
		require.NotEmpty(t, resp.Token)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		email := "nonexistent@example.com"
		password := "password123"

		mockRepo.EXPECT().
			GetByEmail(email).
			Return(nil, errors.New("user not found")).
			Times(1)

		resp, err := userService.Login(&request.Login{
			Email:    email,
			Password: password,
		})

		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("should return invalid credentials when password is incorrect", func(t *testing.T) {
		email := "test@example.com"
		password := "wrongpassword"
		hashedPassword := "$2a$10$abcdefghijklmnopqrstuv"
		role := "user"

		user := &entity.User{
			Id:       uuid.New(),
			Email:    email,
			Password: hashedPassword,
			Role:     role,
		}

		mockRepo.EXPECT().
			GetByEmail(email).
			Return(user, nil).
			Times(1)

		resp, err := userService.Login(&request.Login{
			Email:    email,
			Password: password,
		})

		require.Error(t, err)
		require.Nil(t, resp)
		require.EqualError(t, err, "invalid credentials")
	})
}
