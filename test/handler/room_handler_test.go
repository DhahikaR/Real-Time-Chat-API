package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"real-time-chat-api/internal/handler"
	"real-time-chat-api/internal/models/domain"
	"real-time-chat-api/internal/models/web"
	"real-time-chat-api/test/mocks"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRoomApp(roomService *mocks.MockRoomService) *fiber.App {
	app := fiber.New()
	roomHandler := handler.NewRoomHandler(roomService)

	// Simulate authenticated routes by setting user_id local
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return c.Next()
	})

	app.Post("/rooms", roomHandler.CreateRoom)
	app.Get("/rooms", roomHandler.ListRooms)
	app.Get("/rooms/:id/message", roomHandler.GetMessages)
	return app
}

func TestRoomHandler_CreateRoom_Success(t *testing.T) {
	mockService := new(mocks.MockRoomService)
	app := setupRoomApp(mockService)

	ownerID := uuid.New()
	reqBody := map[string]string{
		"name": "Test Room",
		"type": "public",
	}
	body, _ := json.Marshal(reqBody)

	expectedRoom := domain.Room{
		ID:      uuid.New(),
		Name:    "Test Room",
		Type:    "public",
		OwnerID: ownerID,
	}

	mockService.On("CreateRoom", mock.Anything, mock.AnythingOfType("web.CreateRoomRequest")).Return(expectedRoom, nil)

	req := httptest.NewRequest(http.MethodPost, "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestRoomHandler_CreateRoom_InvalidBody(t *testing.T) {
	mockService := new(mocks.MockRoomService)
	app := setupRoomApp(mockService)

	req := httptest.NewRequest(http.MethodPost, "/rooms", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRoomHandler_CreateRoom_MissingUserID(t *testing.T) {
	mockService := new(mocks.MockRoomService)
	app := fiber.New()
	roomHandler := handler.NewRoomHandler(mockService)
	app.Post("/rooms", roomHandler.CreateRoom)

	reqBody := map[string]string{
		"name": "Test Room",
		"type": "public",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestRoomHandler_CreateRoom_ServiceError(t *testing.T) {
	mockService := new(mocks.MockRoomService)
	app := setupRoomApp(mockService)

	reqBody := map[string]string{
		"name": "Test Room",
		"type": "public",
	}
	body, _ := json.Marshal(reqBody)

	mockService.On("CreateRoom", mock.Anything, mock.AnythingOfType("web.CreateRoomRequest")).Return(domain.Room{}, errors.New("database error"))

	req := httptest.NewRequest(http.MethodPost, "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestRoomHandler_ListRooms_Success(t *testing.T) {
	mockService := new(mocks.MockRoomService)
	app := setupRoomApp(mockService)

	ownerID := uuid.New()
	expectedRooms := []domain.Room{
		{ID: uuid.New(), Name: "Room 1", Type: "public", OwnerID: ownerID},
		{ID: uuid.New(), Name: "Room 2", Type: "public", OwnerID: ownerID},
	}

	mockService.On("ListRooms", mock.Anything).Return(expectedRooms, nil)

	req := httptest.NewRequest(http.MethodGet, "/rooms", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var rooms []domain.Room
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &rooms)
	assert.Len(t, rooms, 2)
	mockService.AssertExpectations(t)
}

func TestRoomHandler_ListRooms_Error(t *testing.T) {
	mockService := new(mocks.MockRoomService)
	app := setupRoomApp(mockService)

	mockService.On("ListRooms", mock.Anything).Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/rooms", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestRoomHandler_GetMessages_Success(t *testing.T) {
	mockService := new(mocks.MockRoomService)
	app := setupRoomApp(mockService)

	roomID := uuid.New()
	senderID := uuid.New()
	now := time.Now()
	expectedMessages := []domain.Message{
		{ID: uuid.New(), RoomID: roomID, SenderID: senderID, Content: "Hello", SentAt: now},
	}

	mockService.On("GetMessages", mock.Anything, roomID, 1, 20).Return(expectedMessages, nil)

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/message", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestRoomHandler_GetMessages_Error(t *testing.T) {
	mockService := new(mocks.MockRoomService)
	app := setupRoomApp(mockService)

	roomID := uuid.New()

	mockService.On("GetMessages", mock.Anything, roomID, 1, 20).Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/message", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestRoomHandler_GetMessages_WithPagination(t *testing.T) {
	mockService := new(mocks.MockRoomService)
	app := setupRoomApp(mockService)

	roomID := uuid.New()
	senderID := uuid.New()
	now := time.Now()
	expectedMessages := []domain.Message{
		{ID: uuid.New(), RoomID: roomID, SenderID: senderID, Content: "Page 2", SentAt: now},
	}

	mockService.On("GetMessages", mock.Anything, roomID, 2, 10).Return(expectedMessages, nil)

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/message?page=2&limit=10", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestRoomHandler_CreateRoom_InvalidUserIDFormat(t *testing.T) {
	mockService := new(mocks.MockRoomService)
	app := fiber.New()
	roomHandler := handler.NewRoomHandler(mockService)

	// Set invalid user_id format
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "not-a-valid-uuid")
		return c.Next()
	})
	app.Post("/rooms", roomHandler.CreateRoom)

	reqBody := web.CreateRoomRequest{
		Name: "Test Room",
		Type: "public",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
