package service_test

import (
	"context"
	"errors"
	"real-time-chat-api/internal/models/domain"
	"real-time-chat-api/internal/models/web"
	"real-time-chat-api/internal/service"
	"real-time-chat-api/test/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := service.NewAuthService(mockRepo, "test-secret")

	req := web.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	// Email does not exist yet
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)
	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	err := authService.Register(context.Background(), req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := service.NewAuthService(mockRepo, "test-secret")

	req := web.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	existingUser := &domain.User{
		ID:    uuid.New(),
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	err := authService.Register(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, "email already exist", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := service.NewAuthService(mockRepo, "test-secret")

	req := web.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, errors.New("database error"))

	err := authService.Register(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_SaveError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := service.NewAuthService(mockRepo, "test-secret")

	req := web.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)
	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("save error"))

	err := authService.Register(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, "save error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := service.NewAuthService(mockRepo, "test-secret")

	req := web.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	existingUser := &domain.User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: string(hashed),
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	token, err := authService.Login(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := service.NewAuthService(mockRepo, "test-secret")

	req := web.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, errors.New("record not found"))

	token, err := authService.Login(context.Background(), req)
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "wrong email or password", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_NilUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := service.NewAuthService(mockRepo, "test-secret")

	req := web.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)

	token, err := authService.Login(context.Background(), req)
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "wrong email or password", err.Error())
	mockRepo.AssertExpectations(t)
}
