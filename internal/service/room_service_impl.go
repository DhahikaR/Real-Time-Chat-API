package service

import (
	"context"
	"real-time-chat-api/internal/models/domain"
	"real-time-chat-api/internal/models/web"
	"real-time-chat-api/internal/repository"

	"github.com/google/uuid"
)

type RoomServiceImpl struct {
	roomRepository    repository.RoomRepository
	messageRepository repository.MessageRepository
}

func NewRoomService(roomRepository repository.RoomRepository, messageRepository repository.MessageRepository) RoomService {
	return &RoomServiceImpl{
		roomRepository:    roomRepository,
		messageRepository: messageRepository}
}

func (service *RoomServiceImpl) CreateRoom(ctx context.Context, request web.CreateRoomRequest) (domain.Room, error) {
	room := &domain.Room{
		ID:      uuid.New(),
		Name:    request.Name,
		Type:    request.Type,
		OwnerID: request.OwnerID,
	}

	if err := service.roomRepository.Save(ctx, room); err != nil {
		return domain.Room{}, err
	}
	return *room, nil
}

func (service *RoomServiceImpl) ListRooms(ctx context.Context) ([]domain.Room, error) {
	return service.roomRepository.FindAll(ctx)
}

func (service *RoomServiceImpl) GetMessages(ctx context.Context, roomID uuid.UUID, page, limit int) ([]domain.Message, error) {
	return service.messageRepository.GetByRoom(ctx, roomID, page, limit)
}
