package mocks

import (
	"context"
	"real-time-chat-api/internal/models/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Save(ctx context.Context, message *domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetByRoom(ctx context.Context, roomID uuid.UUID, page, limit int) ([]domain.Message, error) {
	args := m.Called(ctx, roomID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}
