package ws_test

import (
	"real-time-chat-api/internal/ws"
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTestClient(conn *mockWSConn) (*ws.Client, *ws.Hub) {
	hub := &ws.Hub{
		Broadcast:  make(chan ws.Broadcast, 10), // buffer 10 pesan
		Register:   make(chan *ws.Client, 5),
		Unregister: make(chan *ws.Client, 5),
	}
	client := &ws.Client{
		ID:     uuid.New(),
		UserID: uuid.New(),
		RoomID: uuid.New(),
		Conn:   conn, // pakai mock, bukan *websocket.Conn sungguhan
		Send:   make(chan []byte, 10),
		Hub:    hub,
	}
	return client, hub
}

func waitFor(t *testing.T, condition func() bool, timeout time.Duration, failMsg string) {
	t.Helper() // membuat error message menunjuk ke baris pemanggil, bukan ke sini
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timeout setelah %v — %s", timeout, failMsg)
}
