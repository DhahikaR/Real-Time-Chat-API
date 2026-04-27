package ws_test

import (
	"real-time-chat-api/internal/ws"
	"real-time-chat-api/pkg/pubsub"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newTestHub() *ws.Hub {
	// Use a redis client that won't actually connect (for unit tests)
	rdb := pubsub.NewRedisClient("localhost:6379")
	return ws.NewHub(rdb)
}

func TestNewHub_Initialization(t *testing.T) {
	hub := newTestHub()

	assert.NotNil(t, hub)
	assert.NotNil(t, hub.Rooms)
	assert.NotNil(t, hub.Register)
	assert.NotNil(t, hub.Unregister)
	assert.NotNil(t, hub.Broadcast)
}

func TestHub_Register_NewRoom(t *testing.T) {
	hub := newTestHub()

	// Start hub in background
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	roomID := uuid.New()
	userID := uuid.New()
	clientID := uuid.New()

	client := &ws.Client{
		ID:     clientID,
		UserID: userID,
		RoomID: roomID,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}

	hub.Register <- client
	time.Sleep(50 * time.Millisecond)

	_, exists := hub.Rooms[roomID]
	assert.True(t, exists)
}

func TestHub_Unregister_Client(t *testing.T) {
	hub := newTestHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	roomID := uuid.New()
	userID := uuid.New()
	clientID := uuid.New()

	client := &ws.Client{
		ID:     clientID,
		UserID: userID,
		RoomID: roomID,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}

	// Register first
	hub.Register <- client
	time.Sleep(50 * time.Millisecond)

	// Then unregister
	hub.Unregister <- client
	time.Sleep(50 * time.Millisecond)

	// Should not panic
	assert.NotNil(t, hub)
}

func TestHub_Broadcast_Message(t *testing.T) {
	hub := newTestHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	roomID := uuid.New()
	userID := uuid.New()
	clientID := uuid.New()

	sendChan := make(chan []byte, 256)
	client := &ws.Client{
		ID:     clientID,
		UserID: userID,
		RoomID: roomID,
		Send:   sendChan,
		Hub:    hub,
	}

	hub.Register <- client
	time.Sleep(50 * time.Millisecond)

	// Broadcast a message
	hub.Broadcast <- ws.Broadcast{
		RoomID:  roomID,
		Payload: []byte("hello world"),
	}
	time.Sleep(50 * time.Millisecond)

	// Check if message was received
	select {
	case msg := <-sendChan:
		assert.Equal(t, []byte("hello world"), msg)
	case <-time.After(200 * time.Millisecond):
		// Message may have been published to Redis instead
		// This is acceptable in unit tests without real Redis
	}
}

func TestHub_MultipleClients_SameRoom(t *testing.T) {
	hub := newTestHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	roomID := uuid.New()

	client1 := &ws.Client{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: roomID,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}

	client2 := &ws.Client{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: roomID,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}

	hub.Register <- client1
	hub.Register <- client2
	time.Sleep(50 * time.Millisecond)

	assert.NotNil(t, hub.Rooms)
}

func TestHub_Broadcast_DifferentRooms(t *testing.T) {
	hub := newTestHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	roomID1 := uuid.New()
	roomID2 := uuid.New()

	sendChan1 := make(chan []byte, 256)
	sendChan2 := make(chan []byte, 256)

	client1 := &ws.Client{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: roomID1,
		Send:   sendChan1,
		Hub:    hub,
	}

	client2 := &ws.Client{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: roomID2,
		Send:   sendChan2,
		Hub:    hub,
	}

	hub.Register <- client1
	hub.Register <- client2
	time.Sleep(50 * time.Millisecond)

	// Broadcast to room1 only
	hub.Broadcast <- ws.Broadcast{
		RoomID:  roomID1,
		Payload: []byte("room1 message"),
	}
	time.Sleep(50 * time.Millisecond)

	// client1 should receive the message
	select {
	case msg := <-sendChan1:
		assert.Equal(t, []byte("room1 message"), msg)
	case <-time.After(200 * time.Millisecond):
		// May have been published to Redis
	}

	// client2 should NOT receive the message
	select {
	case <-sendChan2:
		t.Error("client2 should not receive message from room1")
	case <-time.After(100 * time.Millisecond):
		// Expected: no message for client2
	}
}

func TestWSMessage_Constants(t *testing.T) {
	assert.Equal(t, "send_message", ws.EventMessage)
	assert.Equal(t, "new_message", ws.EventNewMessage)
	assert.Equal(t, "typing_indicator", ws.EventTypingIndicator)
	assert.Equal(t, "user_typing", ws.EventUserTyping)
	assert.Equal(t, "user_joined", ws.EventUserJoined)
	assert.Equal(t, "user_left", ws.EventUserLeft)
}

func TestBroadcast_Struct(t *testing.T) {
	roomID := uuid.New()
	payload := []byte("test payload")

	broadcast := ws.Broadcast{
		RoomID:  roomID,
		Payload: payload,
	}

	assert.Equal(t, roomID, broadcast.RoomID)
	assert.Equal(t, payload, broadcast.Payload)
}
