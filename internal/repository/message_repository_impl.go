package repository

import (
	"context"
	"real-time-chat-api/internal/models/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepositoryImpl struct {
	*gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &MessageRepositoryImpl{
		DB: db,
	}
}

func (repository *MessageRepositoryImpl) Save(ctx context.Context, message *domain.Message) error {
	return repository.DB.WithContext(ctx).Create(message).Error
}

func (repository *MessageRepositoryImpl) GetByRoom(ctx context.Context, roomId uuid.UUID, page, limit int) ([]domain.Message, error) {
	var messages []domain.Message
	result := repository.DB.WithContext(ctx).Where("room_id = ? AND is_deleted = false", roomId).
		Order("sent_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&messages)
	return messages, result.Error
}
