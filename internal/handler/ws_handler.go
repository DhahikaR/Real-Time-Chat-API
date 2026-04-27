package handler

import (
	"real-time-chat-api/internal/ws"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WSHandler struct{ hub *ws.Hub }

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

// Connect godoc
// @Summary      Connect to WebSocket
// @Description  Upgrade connection to WebSocket for real-time chat
// @Tags         websocket
// @Param        room_id  query  string  true  "Room ID (UUID)"
// @Success      101
// @Failure      401  {object}  map[string]string
// @Failure      426
// @Security     BearerAuth
// @Router       /ws/connect [get]

func (handler *WSHandler) Connect() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {

		userIDStr, ok := c.Locals("user_id").(string)
		if !ok || userIDStr == "" {
			c.Close()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.Close()
			return
		}

		roomIDStr := c.Query("room_id")
		if roomIDStr == "" {
			c.Close()
			return
		}

		roomID, err := uuid.Parse(roomIDStr)
		if err != nil {
			c.Close()
			return
		}

		client := &ws.Client{
			ID:     uuid.New(),
			UserID: userID,
			RoomID: roomID,
			Conn:   c,
			Send:   make(chan []byte, 256),
			Hub:    handler.hub,
		}

		handler.hub.Register <- client

		go client.WritePump()
		client.ReadPump()
	})

}
