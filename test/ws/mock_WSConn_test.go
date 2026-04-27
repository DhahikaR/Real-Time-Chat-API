package ws_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

// mockWSConn mensimulasikan perilaku *websocket.Conn tanpa koneksi nyata.
// Implements interface ws.WSConn sehingga bisa dipakai sebagai Conn di Client.
type mockWSConn struct {
	mu sync.Mutex // mencegah race condition saat goroutine mengakses field bersamaan

	// --- Untuk ReadMessage() ---
	incomingMsgs [][]byte // antrian pesan yang akan "diterima" (dibaca dari koneksi)
	readIdx      int      // posisi pesan berikutnya yang akan dibaca
	readErr      error    // jika di-set, ReadMessage() langsung return error ini

	// --- Untuk WriteMessage() ---
	writtenMsgs []writtenMsg // semua pesan yang sudah "dikirim", untuk dicek di test
	writeErr    error        // jika di-set, WriteMessage() langsung return error ini

	// --- Status & Handler ---
	closed      bool               // apakah Close() sudah dipanggil?
	pongHandler func(string) error // handler yang di-set via SetPongHandler()
}

// writtenMsg menyimpan detail satu pesan yang ditulis ke koneksi
type writtenMsg struct {
	msgType int
	data    []byte
}

func (m *mockWSConn) ReadMessage() (int, []byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Jika error paksa di-set, langsung kembalikan error itu
	if m.readErr != nil {
		return 0, nil, m.readErr
	}
	// Masih ada pesan di antrian → kembalikan satu per satu
	if m.readIdx < len(m.incomingMsgs) {
		msg := m.incomingMsgs[m.readIdx]
		m.readIdx++
		return websocket.TextMessage, msg, nil
	}
	// Antrian habis → simulasi koneksi ditutup (EOF), ReadPump akan break
	return 0, nil, fmt.Errorf("EOF: no more messages")
}

func (m *mockWSConn) WriteMessage(msgType int, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.writeErr != nil {
		return m.writeErr
	}
	// Simpan pesan yang "dikirim" agar bisa diverifikasi di test
	m.writtenMsgs = append(m.writtenMsgs, writtenMsg{
		msgType: msgType,
		data:    data,
	})
	return nil
}

// Method lain tidak perlu logika khusus untuk test dasar
func (m *mockWSConn) SetReadLimit(limit int64)                    {}
func (m *mockWSConn) SetReadDeadline(t time.Time) error           { return nil }
func (m *mockWSConn) SetWriteDeadline(t time.Time) error          { return nil }
func (m *mockWSConn) SetPongHandler(h func(appData string) error) { m.pongHandler = h }
func (m *mockWSConn) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

// getWritten adalah helper thread-safe untuk membaca salinan writtenMsgs.
// Kenapa tidak langsung baca m.writtenMsgs? Karena goroutine lain bisa
// sedang menulis ke slice itu bersamaan → race condition → data corrupt.
func (m *mockWSConn) getWritten() []writtenMsg {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]writtenMsg, len(m.writtenMsgs))
	copy(result, m.writtenMsgs)
	return result
}

// isClosed adalah helper thread-safe untuk cek status koneksi
func (m *mockWSConn) isClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}
