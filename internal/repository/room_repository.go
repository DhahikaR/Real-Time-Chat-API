package repository

import (
	"context"
	"real-time-chat-api/internal/models/domain"

	"github.com/google/uuid"
)

type RoomRepository interface {
	Save(ctx context.Context, room *domain.Room) error
	FindAll(ctx context.Context) ([]domain.Room, error)
	FindById(ctx context.Context, id uuid.UUID) (*domain.Room, error)
}
