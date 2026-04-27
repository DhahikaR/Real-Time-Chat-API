package client_test

import (
	"real-time-chat-api/pkg/pubsub"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisClient_ReturnsClient(t *testing.T) {
	// NewRedisClient should return a non-nil client even with invalid address
	// (redis client is lazy - it doesn't connect until first command)
	client := pubsub.NewRedisClient("localhost:6379")
	assert.NotNil(t, client)
}

func TestNewRedisClient_CustomAddress(t *testing.T) {
	client := pubsub.NewRedisClient("redis.example.com:6380")
	assert.NotNil(t, client)
}

func TestNewRedisClient_EmptyAddress(t *testing.T) {
	// Even with empty address, client should be created (lazy connection)
	client := pubsub.NewRedisClient("")
	assert.NotNil(t, client)
}

func TestRoomChannel_Format(t *testing.T) {
	roomID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	channel := pubsub.RoomChannel(roomID)

	assert.Equal(t, "room:550e8400-e29b-41d4-a716-446655440000", channel)
}

func TestRoomChannel_UniquePerRoom(t *testing.T) {
	roomID1 := uuid.New()
	roomID2 := uuid.New()

	channel1 := pubsub.RoomChannel(roomID1)
	channel2 := pubsub.RoomChannel(roomID2)

	assert.NotEqual(t, channel1, channel2)
}

func TestRoomChannel_ContainsRoomPrefix(t *testing.T) {
	roomID := uuid.New()
	channel := pubsub.RoomChannel(roomID)

	assert.Contains(t, channel, "room:")
	assert.Contains(t, channel, roomID.String())
}

func TestNewRedisClient_ConnectionOptions(t *testing.T) {
	addr := "localhost:6379"
	client := pubsub.NewRedisClient(addr)
	assert.NotNil(t, client)

	// Verify the client options are set correctly
	opts := client.Options()
	assert.Equal(t, addr, opts.Addr)
}
