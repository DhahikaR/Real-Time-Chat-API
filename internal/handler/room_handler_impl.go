package handler

import (
	"real-time-chat-api/internal/helper"
	"real-time-chat-api/internal/models/web"
	"real-time-chat-api/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RoomHandlerImpl represents room handler implementation
type RoomHandlerImpl struct {
	roomService service.RoomService
}

// NewRoomHandler creates new room handler
// @Summary Create new room
// @Description Create a new chat room
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param CreateRoomRequest body web.CreateRoomRequest true "Create room request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /rooms [post]
func NewRoomHandler(roomService service.RoomService) RoomHandler {
	return &RoomHandlerImpl{roomService: roomService}
}

// CreateRoom handles room creation
// @Summary Create new room
// @Description Create a new chat room
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param CreateRoomRequest body web.CreateRoomRequest true "Create room request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /rooms [post]
func (handler *RoomHandlerImpl) CreateRoom(c *fiber.Ctx) error {
	roomCreateRequest := web.CreateRoomRequest{}
	if err := helper.ReadFromRequestBody(c, &roomCreateRequest); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return helper.Forbidden(c, "unauthorized: missing or invalid user_id")
	}

	ownerID, err := uuid.Parse(userID)
	if err != nil {
		return helper.BadRequest(c, "invalid user_id format")
	}
	roomCreateRequest.OwnerID = ownerID

	room, err := handler.roomService.CreateRoom(c.Context(), roomCreateRequest)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(room)
}

// ListRooms handles listing all rooms
// @Summary List all rooms
// @Description Get all available chat rooms
// @Tags rooms
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /rooms [get]
func (handler *RoomHandlerImpl) ListRooms(c *fiber.Ctx) error {
	rooms, err := handler.roomService.ListRooms(c.Context())
	if err != nil {
		return helper.InternalServerError(c, "Internal Server Error")
	}
	return c.Status(fiber.StatusOK).JSON(rooms)
}

// GetMessages handles getting room messages
// @Summary Get room messages
// @Description Get messages from a specific room with pagination
// @Tags rooms
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Messages per page (default 20)"
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /rooms/{id}/message [get]
func (handler *RoomHandlerImpl) GetMessages(c *fiber.Ctx) error {
	roomID, _ := uuid.Parse(c.Params("id"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	messages, err := handler.roomService.GetMessages(c.Context(), roomID, page, limit)

	if err != nil {
		return helper.InternalServerError(c, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}
