package mocks

import (
	"context"
	"real-time-chat-api/internal/models/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) Save(ctx context.Context, room *domain.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockRoomRepository) FindAll(ctx context.Context) ([]domain.Room, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Room), args.Error(1)
}

func (m *MockRoomRepository) FindById(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Room), args.Error(1)
}
