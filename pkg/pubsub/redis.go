package pubsub

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func RoomChannel(roomID uuid.UUID) string {
	return fmt.Sprintf("room:%s", roomID)
}
