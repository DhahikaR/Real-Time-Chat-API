package mocks

import (
	"context"
	"real-time-chat-api/internal/models/domain"
	"real-time-chat-api/internal/models/web"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRoomService struct {
	mock.Mock
}

func (m *MockRoomService) CreateRoom(ctx context.Context, request web.CreateRoomRequest) (domain.Room, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(domain.Room), args.Error(1)
}

func (m *MockRoomService) ListRooms(ctx context.Context) ([]domain.Room, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Room), args.Error(1)
}

func (m *MockRoomService) GetMessages(ctx context.Context, roomID uuid.UUID, page, limit int) ([]domain.Message, error) {
	args := m.Called(ctx, roomID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}
