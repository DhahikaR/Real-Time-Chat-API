package mocks

import (
	"context"
	"real-time-chat-api/internal/models/web"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, request web.RegisterRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockAuthService) Login(ctx context.Context, request web.LoginRequest) (string, error) {
	args := m.Called(ctx, request)
	return args.String(0), args.Error(1)
}
