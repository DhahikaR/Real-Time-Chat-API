package repository

import (
	"context"
	"errors"
	"real-time-chat-api/internal/models/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepositoryImpl struct {
	*gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &RoomRepositoryImpl{
		DB: db,
	}
}

func (repository *RoomRepositoryImpl) Save(ctx context.Context, room *domain.Room) error {
	return repository.DB.WithContext(ctx).Create(room).Error
}

func (repository *RoomRepositoryImpl) FindAll(ctx context.Context) ([]domain.Room, error) {
	var rooms []domain.Room
	result := repository.DB.WithContext(ctx).Where("type = ?", "public").Find(&rooms)
	return rooms, result.Error
}

func (repository *RoomRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	var room domain.Room
	result := repository.DB.WithContext(ctx).First(&room, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &room, result.Error
}
