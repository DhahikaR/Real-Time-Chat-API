package repository

import (
	"context"
	"real-time-chat-api/internal/models/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*domain.User, error)
}
