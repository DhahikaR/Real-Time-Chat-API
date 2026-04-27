package service_test

import (
	"context"
	"errors"
	"real-time-chat-api/internal/models/domain"
	"real-time-chat-api/internal/models/web"
	"real-time-chat-api/internal/service"
	"real-time-chat-api/test/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoomService_CreateRoom_Success(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)
	roomService := service.NewRoomService(mockRoomRepo, mockMsgRepo)

	ownerID := uuid.New()
	req := web.CreateRoomRequest{
		Name:    "Test Room",
		Type:    "public",
		OwnerID: ownerID,
	}

	mockRoomRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Room")).Return(nil)

	room, err := roomService.CreateRoom(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, req.Name, room.Name)
	assert.Equal(t, req.Type, room.Type)
	assert.Equal(t, ownerID, room.OwnerID)
	mockRoomRepo.AssertExpectations(t)
}

func TestRoomService_CreateRoom_Error(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)
	roomService := service.NewRoomService(mockRoomRepo, mockMsgRepo)

	ownerID := uuid.New()
	req := web.CreateRoomRequest{
		Name:    "Test Room",
		Type:    "public",
		OwnerID: ownerID,
	}

	mockRoomRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Room")).Return(errors.New("database error"))

	room, err := roomService.CreateRoom(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, domain.Room{}, room)
	mockRoomRepo.AssertExpectations(t)
}

func TestRoomService_ListRooms_Success(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)
	roomService := service.NewRoomService(mockRoomRepo, mockMsgRepo)

	ownerID := uuid.New()
	expectedRooms := []domain.Room{
		{ID: uuid.New(), Name: "Room 1", Type: "public", OwnerID: ownerID},
		{ID: uuid.New(), Name: "Room 2", Type: "public", OwnerID: ownerID},
	}

	mockRoomRepo.On("FindAll", mock.Anything).Return(expectedRooms, nil)

	rooms, err := roomService.ListRooms(context.Background())
	assert.NoError(t, err)
	assert.Len(t, rooms, 2)
	mockRoomRepo.AssertExpectations(t)
}

func TestRoomService_ListRooms_Error(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)
	roomService := service.NewRoomService(mockRoomRepo, mockMsgRepo)

	mockRoomRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

	rooms, err := roomService.ListRooms(context.Background())
	assert.Error(t, err)
	assert.Nil(t, rooms)
	mockRoomRepo.AssertExpectations(t)
}

func TestRoomService_GetMessages_Success(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)
	roomService := service.NewRoomService(mockRoomRepo, mockMsgRepo)

	roomID := uuid.New()
	senderID := uuid.New()
	now := time.Now()
	expectedMessages := []domain.Message{
		{ID: uuid.New(), RoomID: roomID, SenderID: senderID, Content: "Hello", SentAt: now},
		{ID: uuid.New(), RoomID: roomID, SenderID: senderID, Content: "World", SentAt: now},
	}

	mockMsgRepo.On("GetByRoom", mock.Anything, roomID, 1, 20).Return(expectedMessages, nil)

	messages, err := roomService.GetMessages(context.Background(), roomID, 1, 20)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	mockMsgRepo.AssertExpectations(t)
}

func TestRoomService_GetMessages_Error(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)
	roomService := service.NewRoomService(mockRoomRepo, mockMsgRepo)

	roomID := uuid.New()

	mockMsgRepo.On("GetByRoom", mock.Anything, roomID, 1, 20).Return(nil, errors.New("database error"))

	messages, err := roomService.GetMessages(context.Background(), roomID, 1, 20)
	assert.Error(t, err)
	assert.Nil(t, messages)
	mockMsgRepo.AssertExpectations(t)
}

func TestRoomService_GetMessages_Pagination(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)
	roomService := service.NewRoomService(mockRoomRepo, mockMsgRepo)

	roomID := uuid.New()
	senderID := uuid.New()
	now := time.Now()
	expectedMessages := []domain.Message{
		{ID: uuid.New(), RoomID: roomID, SenderID: senderID, Content: "Page 2 Msg", SentAt: now},
	}

	mockMsgRepo.On("GetByRoom", mock.Anything, roomID, 2, 10).Return(expectedMessages, nil)

	messages, err := roomService.GetMessages(context.Background(), roomID, 2, 10)
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	mockMsgRepo.AssertExpectations(t)
}
