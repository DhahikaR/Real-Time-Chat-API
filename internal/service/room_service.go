package service

import (
	"context"
	"real-time-chat-api/internal/models/domain"
	"real-time-chat-api/internal/models/web"

	"github.com/google/uuid"
)

type RoomService interface {
	CreateRoom(ctx context.Context, request web.CreateRoomRequest) (domain.Room, error)
	ListRooms(ctx context.Context) ([]domain.Room, error)
	GetMessages(ctx context.Context, roomID uuid.UUID, page, limit int) ([]domain.Message, error)
}
