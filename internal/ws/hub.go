package ws

import (
	"context"
	"real-time-chat-api/pkg/pubsub"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Hub struct {
	Rooms      map[uuid.UUID]map[uuid.UUID]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Broadcast
	mu         sync.RWMutex
	rdb        *redis.Client
}

func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		Rooms:      make(map[uuid.UUID]map[uuid.UUID]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Broadcast, 256),
		rdb:        rdb,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if _, ok := h.Rooms[client.RoomID]; !ok {
				h.Rooms[client.RoomID] = make(map[uuid.UUID]*Client)
				go h.subscribeRoom(client.RoomID) // subscribe Redis saat room pertama kali aktif
			}
			h.Rooms[client.RoomID][client.ID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if room, ok := h.Rooms[client.RoomID]; ok {
				delete(room, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			// Kirim ke semua client lokal di room ini
			h.mu.RLock()
			for _, client := range h.Rooms[msg.RoomID] {
				select {
				case client.Send <- msg.Payload:
				default:
					close(client.Send)
				}
			}
			h.mu.RUnlock()

			// Publish ke Redis untuk instance server lain
			go h.publishToRedis(msg)
		}
	}
}

func (h *Hub) publishToRedis(msg Broadcast) {
	channel := pubsub.RoomChannel(msg.RoomID)
	h.rdb.Publish(context.Background(), channel, msg.Payload)
}

// subscribeRoom men-subscribe channel Redis untuk room tertentu.
// Pesan dari Redis diteruskan ke semua client lokal di room tersebut.
func (h *Hub) subscribeRoom(roomID uuid.UUID) {
	channel := pubsub.RoomChannel(roomID)
	sub := h.rdb.Subscribe(context.Background(), channel)
	defer sub.Close()

	for redisMsg := range sub.Channel() {
		payload := []byte(redisMsg.Payload)
		h.mu.RLock()
		for _, client := range h.Rooms[roomID] {
			select {
			case client.Send <- payload:
			default:
			}
		}
		h.mu.RUnlock()
	}
}
