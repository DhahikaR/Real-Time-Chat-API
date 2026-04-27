package repository

import (
	"context"
	"real-time-chat-api/internal/models/domain"

	"github.com/google/uuid"
)

type MessageRepository interface {
	Save(ctx context.Context, message *domain.Message) error
	GetByRoom(ctx context.Context, roomID uuid.UUID, page, limit int) ([]domain.Message, error)
}
