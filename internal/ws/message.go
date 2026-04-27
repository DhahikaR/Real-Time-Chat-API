package ws

import (
	"encoding/json"

	"github.com/google/uuid"
)

type WSMessage struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type Broadcast struct {
	RoomID  uuid.UUID
	Payload []byte
}

const (
	EventMessage         = "send_message"
	EventNewMessage      = "new_message"
	EventTypingIndicator = "typing_indicator"
	EventUserTyping      = "user_typing"
	EventUserJoined      = "user_joined"
	EventUserLeft        = "user_left"
)
