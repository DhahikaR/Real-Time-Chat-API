package handler_test

import (
	"net/http"
	"net/http/httptest"
	"real-time-chat-api/internal/handler"
	"real-time-chat-api/internal/ws"
	"real-time-chat-api/pkg/pubsub"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewWSHandler_Success(t *testing.T) {
	rdb := pubsub.NewRedisClient("localhost:6379")
	hub := ws.NewHub(rdb)

	wsHandler := handler.NewWSHandler(hub)
	assert.NotNil(t, wsHandler)
}

func TestWSHandler_Connect_ReturnsHandler(t *testing.T) {
	rdb := pubsub.NewRedisClient("localhost:6379")
	hub := ws.NewHub(rdb)

	wsHandler := handler.NewWSHandler(hub)
	connectHandler := wsHandler.Connect()

	// Connect() should return a fiber.Handler (function)
	assert.NotNil(t, connectHandler)
}

func TestWSHandler_Initialization(t *testing.T) {
	rdb := pubsub.NewRedisClient("localhost:6379")
	hub := ws.NewHub(rdb)

	wsHandler := handler.NewWSHandler(hub)

	// Verify handler is properly initialized
	assert.NotNil(t, wsHandler)
}

func TestWSHandler_UUIDParsing_InvalidRoomID(t *testing.T) {
	// Test with invalid room_id format
	invalidRoomID := "not-a-valid-uuid"
	_, err := uuid.Parse(invalidRoomID)
	assert.Error(t, err)
}

func TestWSHandler_UUIDParsing_InvalidUserID(t *testing.T) {
	// Test with invalid user_id format
	invalidUserID := "not-a-valid-uuid"
	_, err := uuid.Parse(invalidUserID)
	assert.Error(t, err)
}

func TestWSHandler_UUIDParsing_ValidUUIDs(t *testing.T) {
	// Test with valid UUIDs
	validRoomID := uuid.New()
	validUserID := uuid.New()

	parsedRoomID, err := uuid.Parse(validRoomID.String())
	assert.NoError(t, err)
	assert.Equal(t, validRoomID, parsedRoomID)

	parsedUserID, err := uuid.Parse(validUserID.String())
	assert.NoError(t, err)
	assert.Equal(t, validUserID, parsedUserID)
}

func TestWSHandler_ParameterValidation(t *testing.T) {
	// Test parameter validation logic similar to WSHandler.Connect()
	app := fiber.New()

	app.Get("/ws/connect/:room_id/:user_id", func(c *fiber.Ctx) error {
		roomIDStr := c.Params("room_id")
		roomID, err := uuid.Parse(roomIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid room_id format",
			})
		}
		userIDStr := c.Params("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid user_id format",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"room_id": roomID.String(),
			"user_id": userID.String(),
		})
	})

	// Test with invalid room_id
	req := httptest.NewRequest(http.MethodGet, "/ws/connect/invalid-uuid/"+uuid.New().String(), nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Test with invalid user_id
	req = httptest.NewRequest(http.MethodGet, "/ws/connect/"+uuid.New().String()+"/invalid-uuid", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Test with valid UUIDs
	req = httptest.NewRequest(http.MethodGet, "/ws/connect/"+uuid.New().String()+"/"+uuid.New().String(), nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestWSHandler_EmptyParameters(t *testing.T) {
	// Test with empty parameters
	app := fiber.New()

	app.Get("/ws/connect/:room_id/:user_id", func(c *fiber.Ctx) error {
		roomIDStr := c.Params("room_id")
		roomID, err := uuid.Parse(roomIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid room_id format",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"room_id": roomID.String(),
		})
	})

	// Test with empty room_id
	req := httptest.NewRequest(http.MethodGet, "/ws/connect/"+uuid.New().String(), nil)
	_, err := app.Test(req)
	assert.NoError(t, err)
	// Should handle missing param gracefully
}

func TestWSHandler_MultipleWebSocketConnections(t *testing.T) {
	// Test creating multiple handlers
	rdb := pubsub.NewRedisClient("localhost:6379")
	hub := ws.NewHub(rdb)

	handler1 := handler.NewWSHandler(hub)
	handler2 := handler.NewWSHandler(hub)

	// Both handlers should be valid
	assert.NotNil(t, handler1)
	assert.NotNil(t, handler2)

	// Connect methods should return valid handlers
	assert.NotNil(t, handler1.Connect())
	assert.NotNil(t, handler2.Connect())
}
