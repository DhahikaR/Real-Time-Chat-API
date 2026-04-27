package repository

import (
	"context"
	"errors"
	"real-time-chat-api/internal/models/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, user *domain.User) error {
	return repository.DB.WithContext(ctx).Create(user).Error
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	result := repository.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	result := repository.DB.WithContext(ctx).First(&user, "id = ?", id)
	return &user, result.Error
}
