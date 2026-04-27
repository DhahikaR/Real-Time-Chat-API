package ws_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestReadPump_Success_SingleMessage(t *testing.T) {
	// ARRANGE — siapkan koneksi mock dengan 1 pesan di antrian
	conn := &mockWSConn{
		incomingMsgs: [][]byte{[]byte("halo dunia")},
	}
	client, hub := newTestClient(conn)

	// ACT — jalankan ReadPump di goroutine (karena dia loop tak terbatas)
	go client.ReadPump()

	// ASSERT — tunggu sampai pesan masuk ke hub.Broadcast
	select {
	case broadcast := <-hub.Broadcast:
		assert.Equal(t, client.RoomID, broadcast.RoomID,
			"[SUCCESS] RoomID di broadcast harus sesuai milik client")
		assert.Equal(t, []byte("halo dunia"), broadcast.Payload,
			"[SUCCESS] Payload harus sama persis dengan pesan aslinya")
	case <-time.After(time.Second):
		t.Fatal("[SUCCESS] GAGAL: pesan tidak sampai ke hub dalam 1 detik")
	}
}

func TestReadPump_Success_MultipleMessages(t *testing.T) {
	// ARRANGE
	expectedMessages := [][]byte{
		[]byte("pesan-1"),
		[]byte("pesan-2"),
		[]byte("pesan-3"),
	}
	conn := &mockWSConn{incomingMsgs: expectedMessages}
	client, hub := newTestClient(conn)

	// ACT
	go client.ReadPump()

	// ASSERT — verifikasi satu per satu, urutan HARUS terjaga
	for i, expected := range expectedMessages {
		select {
		case broadcast := <-hub.Broadcast:
			assert.Equal(t, client.RoomID, broadcast.RoomID,
				"[SUCCESS] RoomID pesan ke-%d harus benar", i+1)
			assert.Equal(t, expected, broadcast.Payload,
				"[SUCCESS] Isi pesan ke-%d harus sesuai urutan (FIFO)", i+1)
		case <-time.After(time.Second):
			t.Fatalf("[SUCCESS] GAGAL: pesan ke-%d tidak sampai ke hub", i+1)
		}
	}
}

func TestReadPump_Success_PongHandlerRegistered(t *testing.T) {
	// ARRANGE — langsung error agar ReadPump cepat keluar setelah fase setup
	conn := &mockWSConn{
		readErr: fmt.Errorf("force stop"),
	}
	client, _ := newTestClient(conn)

	// ACT
	go client.ReadPump()

	// Tunggu sampai ReadPump menyelesaikan fase inisialisasi (setup handler)
	waitFor(t, func() bool {
		conn.mu.Lock()
		defer conn.mu.Unlock()
		return conn.pongHandler != nil
	}, time.Second, "[SUCCESS] pongHandler harus sudah di-set setelah ReadPump mulai")

	// ASSERT 1 — handler harus sudah di-set
	conn.mu.Lock()
	handler := conn.pongHandler
	conn.mu.Unlock()
	assert.NotNil(t, handler, "[SUCCESS] pongHandler tidak boleh nil")

	// ASSERT 2 — handler harus bisa dipanggil dan return nil (tanpa error)
	// Ini mensimulasikan: server mengirim pong → handler dipanggil
	err := handler("")
	assert.NoError(t, err,
		"[SUCCESS] pongHandler harus return nil karena tugasnya hanya reset deadline")
}

func TestReadPump_Success_ConnectionClosedWhenDone(t *testing.T) {
	// ARRANGE — tidak ada pesan → ReadPump langsung selesai (EOF)
	conn := &mockWSConn{
		incomingMsgs: [][]byte{},
	}
	client, _ := newTestClient(conn)

	// ACT
	go client.ReadPump()

	// ASSERT
	waitFor(t, conn.isClosed, time.Second,
		"[SUCCESS] koneksi harus ditutup via defer setelah ReadPump selesai")

	assert.True(t, conn.isClosed(),
		"[SUCCESS] Close() harus dipanggil untuk membebaskan resource koneksi")
}
func TestReadPump_Success_ClientUnregisteredOnExit(t *testing.T) {
	// ARRANGE — antrian kosong, ReadPump langsung exit
	conn := &mockWSConn{
		incomingMsgs: [][]byte{},
	}
	client, hub := newTestClient(conn)

	// ACT
	go client.ReadPump()

	// ASSERT
	select {
	case unregistered := <-hub.Unregister:
		assert.Equal(t, client, unregistered,
			"[SUCCESS] client yang di-unregister harus sama dengan client yang selesai")
	case <-time.After(time.Second):
		t.Fatal("[SUCCESS] GAGAL: client tidak di-unregister setelah ReadPump selesai")
	}
}

func TestReadPump_Error_UnregisterOnImmediateReadError(t *testing.T) {
	// ARRANGE
	conn := &mockWSConn{
		readErr: fmt.Errorf("websocket: connection reset by peer"),
	}
	client, hub := newTestClient(conn)

	// ACT
	go client.ReadPump()

	// ASSERT
	select {
	case unregistered := <-hub.Unregister:
		assert.Equal(t, client, unregistered,
			"[ERROR] client harus di-unregister meski ReadMessage() langsung error")
	case <-time.After(time.Second):
		t.Fatal("[ERROR] GAGAL: client tidak di-unregister setelah error koneksi")
	}
}

func TestReadPump_Error_ConnectionClosedOnReadError(t *testing.T) {
	// ARRANGE
	conn := &mockWSConn{
		readErr: fmt.Errorf("websocket: read tcp: use of closed network connection"),
	}
	client, _ := newTestClient(conn)

	// ACT
	go client.ReadPump()

	// ASSERT — Close() harus dipanggil bahkan saat error
	waitFor(t, conn.isClosed, time.Second,
		"[ERROR] koneksi harus ditutup via defer meski ReadMessage() error")

	assert.True(t, conn.isClosed(),
		"[ERROR] Close() wajib dipanggil untuk membebaskan resource koneksi")
}

func TestReadPump_Error_NoBroadcastOnImmediateError(t *testing.T) {
	// ARRANGE
	conn := &mockWSConn{
		readErr: fmt.Errorf("network error dari awal"),
	}
	client, hub := newTestClient(conn)

	// ACT
	go client.ReadPump()

	// Tunggu sampai ReadPump benar-benar selesai (koneksi ditutup)
	waitFor(t, conn.isClosed, time.Second, "[ERROR] ReadPump harus sudah selesai")

	// ASSERT — channel Broadcast harus kosong (tidak ada pesan yang berhasil masuk)
	select {
	case strayMsg := <-hub.Broadcast:
		t.Fatalf("[ERROR] TIDAK BOLEH ada pesan di Broadcast saat error, dapat: %v",
			strayMsg.Payload)
	default:
		// Channel kosong = benar ✓
		assert.True(t, true,
			"[ERROR] hub.Broadcast harus kosong saat ReadMessage() langsung error")
	}
}

func TestReadPump_Error_MidwayDisconnect(t *testing.T) {
	// ARRANGE — 2 pesan sukses, lalu EOF (simulasi disconnect)
	conn := &mockWSConn{
		incomingMsgs: [][]byte{
			[]byte("pesan-sebelum-putus-1"),
			[]byte("pesan-sebelum-putus-2"),
		},
		// Setelah 2 pesan habis, ReadMessage() akan return EOF otomatis
	}
	client, hub := newTestClient(conn)

	// ACT
	go client.ReadPump()

	// ASSERT 1 — 2 pesan sukses harus berhasil masuk ke hub
	for i := 1; i <= 2; i++ {
		select {
		case broadcast := <-hub.Broadcast:
			assert.NotEmpty(t, broadcast.Payload,
				"[ERROR/Midway] pesan ke-%d sebelum disconnect harus ada isinya", i)
			assert.Equal(t, client.RoomID, broadcast.RoomID,
				"[ERROR/Midway] RoomID pesan ke-%d harus benar", i)
		case <-time.After(time.Second):
			t.Fatalf("[ERROR/Midway] pesan ke-%d tidak sampai ke hub", i)
		}
	}

	// ASSERT 2 — setelah semua pesan habis (EOF), client harus di-unregister
	select {
	case unregistered := <-hub.Unregister:
		assert.Equal(t, client, unregistered,
			"[ERROR/Midway] client harus di-unregister setelah koneksi terputus")
	case <-time.After(time.Second):
		t.Fatal("[ERROR/Midway] client tidak di-unregister setelah EOF")
	}
}

func TestReadPump_Error_BroadcastPreservedBeforeError(t *testing.T) {
	// ARRANGE — 1 pesan sukses, lalu error
	roomID := uuid.New() // gunakan RoomID spesifik agar bisa diverifikasi
	conn := &mockWSConn{
		incomingMsgs: [][]byte{[]byte("pesan-valid")},
		// Setelah 1 pesan, ReadMessage() akan EOF
	}
	client, hub := newTestClient(conn)
	client.RoomID = roomID // override RoomID dengan nilai spesifik

	// ACT
	go client.ReadPump()

	// ASSERT — pesan yang sempat masuk harus tetap valid
	select {
	case broadcast := <-hub.Broadcast:
		assert.Equal(t, roomID, broadcast.RoomID,
			"[ERROR] RoomID pesan yang masuk sebelum error harus tetap benar")
		assert.Equal(t, []byte("pesan-valid"), broadcast.Payload,
			"[ERROR] Payload pesan yang masuk sebelum error tidak boleh berubah")
	case <-time.After(time.Second):
		t.Fatal("[ERROR] GAGAL: pesan sebelum error tidak sampai ke hub")
	}
}

func TestWritePump_Success_SingleMessage(t *testing.T) {
	// ARRANGE
	conn := &mockWSConn{}
	client, _ := newTestClient(conn)

	// ACT — jalankan WritePump di goroutine (karena dia loop tak terbatas)
	go client.WritePump()

	// Kirim pesan ke channel Send — ini mensimulasikan Hub mengirim ke client
	client.Send <- []byte("kirim ke browser")

	// ASSERT — tunggu sampai pesan benar-benar ditulis ke koneksi
	waitFor(t, func() bool {
		for _, m := range conn.getWritten() {
			if m.msgType == websocket.TextMessage &&
				string(m.data) == "kirim ke browser" {
				return true
			}
		}
		return false
	}, time.Second, "[SUCCESS] pesan harus terkirim sebagai TextMessage")

	// Verifikasi detail pesan lebih spesifik
	written := conn.getWritten()
	assert.Equal(t, websocket.TextMessage, written[0].msgType,
		"[SUCCESS] tipe pesan harus TextMessage (bukan Ping, Close, atau Binary)")
	assert.Equal(t, []byte("kirim ke browser"), written[0].data,
		"[SUCCESS] isi pesan tidak boleh berubah saat diteruskan")
}

func TestWritePump_Success_MultipleMessagesInOrder(t *testing.T) {
	// ARRANGE
	conn := &mockWSConn{}
	client, _ := newTestClient(conn)
	go client.WritePump()

	// ACT — kirim 4 pesan berurutan ke channel Send
	messages := []string{"pertama", "kedua", "ketiga", "keempat"}
	for _, msg := range messages {
		client.Send <- []byte(msg)
	}

	// ASSERT — tunggu semua pesan terkirim
	waitFor(t, func() bool {
		return len(conn.getWritten()) >= len(messages)
	}, time.Second, "[SUCCESS] semua 4 pesan harus terkirim ke koneksi")

	written := conn.getWritten()
	for i, expected := range messages {
		assert.Equal(t, websocket.TextMessage, written[i].msgType,
			"[SUCCESS] pesan ke-%d harus bertipe TextMessage", i+1)
		assert.Equal(t, []byte(expected), written[i].data,
			"[SUCCESS] pesan ke-%d harus sesuai urutan pengiriman (FIFO)", i+1)
	}
}
func TestWritePump_Success_CloseMessageOnChannelClose(t *testing.T) {
	// ARRANGE
	conn := &mockWSConn{}
	client, _ := newTestClient(conn)

	// ACT
	go client.WritePump()
	close(client.Send) // simulasi: Hub menutup channel karena shutdown/kick

	// ASSERT — CloseMessage harus dikirim
	waitFor(t, func() bool {
		for _, m := range conn.getWritten() {
			if m.msgType == websocket.CloseMessage {
				return true
			}
		}
		return false
	}, time.Second, "[SUCCESS] CloseMessage harus dikirim saat channel Send ditutup")

	// Verifikasi jumlah CloseMessage: harus tepat 1, tidak lebih
	var closeCount int
	for _, m := range conn.getWritten() {
		if m.msgType == websocket.CloseMessage {
			closeCount++
		}
	}
	assert.Equal(t, 1, closeCount,
		"[SUCCESS] harus ada tepat 1 CloseMessage, tidak boleh lebih")
}

func TestWritePump_Success_ConnectionClosedAfterChannelClose(t *testing.T) {
	// ARRANGE
	conn := &mockWSConn{}
	client, _ := newTestClient(conn)

	// ACT
	go client.WritePump()
	close(client.Send)

	// ASSERT — Close() harus dipanggil via defer setelah WritePump selesai
	waitFor(t, conn.isClosed, time.Second,
		"[SUCCESS] koneksi harus ditutup via defer setelah WritePump return")

	assert.True(t, conn.isClosed(),
		"[SUCCESS] Close() wajib dipanggil untuk mencegah connection leak")
}

func TestWritePump_Success_MessagesBeforeChannelClose(t *testing.T) {
	// ARRANGE
	conn := &mockWSConn{}
	client, _ := newTestClient(conn)
	go client.WritePump()

	// ACT — kirim 2 pesan normal, lalu tutup channel
	client.Send <- []byte("pesan-terakhir-1")
	client.Send <- []byte("pesan-terakhir-2")
	close(client.Send)

	// Tunggu koneksi tertutup (tanda WritePump sudah selesai)
	waitFor(t, conn.isClosed, time.Second, "[SUCCESS] WritePump harus selesai")

	// ASSERT — verifikasi urutan: TextMessage → TextMessage → CloseMessage
	written := conn.getWritten()
	assert.GreaterOrEqual(t, len(written), 3,
		"[SUCCESS] harus ada minimal 3 pesan: 2 TextMessage + 1 CloseMessage")

	// Dua pesan pertama harus TextMessage
	assert.Equal(t, websocket.TextMessage, written[0].msgType,
		"[SUCCESS] pesan ke-1 harus TextMessage")
	assert.Equal(t, []byte("pesan-terakhir-1"), written[0].data)

	assert.Equal(t, websocket.TextMessage, written[1].msgType,
		"[SUCCESS] pesan ke-2 harus TextMessage")
	assert.Equal(t, []byte("pesan-terakhir-2"), written[1].data)

	// Pesan terakhir harus CloseMessage
	lastMsg := written[len(written)-1]
	assert.Equal(t, websocket.CloseMessage, lastMsg.msgType,
		"[SUCCESS] pesan terakhir harus CloseMessage setelah channel ditutup")
}

func TestWritePump_Error_WriteMessageFails(t *testing.T) {
	// ARRANGE — semua WriteMessage() akan gagal
	conn := &mockWSConn{
		writeErr: fmt.Errorf("broken pipe: koneksi sudah putus di sisi network"),
	}
	client, _ := newTestClient(conn)

	// ACT
	go client.WritePump()
	client.Send <- []byte("pesan yang akan gagal")

	// Beri waktu WritePump memproses
	time.Sleep(100 * time.Millisecond)

	// ASSERT — tidak ada pesan yang berhasil tersimpan karena selalu error
	assert.Empty(t, conn.getWritten(),
		"[ERROR] writtenMsgs harus kosong karena WriteMessage() selalu return error")
}
func TestWritePump_Error_NoMessageSentAfterClose(t *testing.T) {
	// ARRANGE
	conn := &mockWSConn{}
	client, _ := newTestClient(conn)

	// ACT
	go client.WritePump()
	close(client.Send) // tutup channel tanpa mengirim pesan apapun

	// Tunggu WritePump selesai
	waitFor(t, conn.isClosed, time.Second,
		"[ERROR] WritePump harus selesai setelah channel ditutup")

	// ASSERT — hanya CloseMessage yang boleh ada
	written := conn.getWritten()
	for _, m := range written {
		assert.Equal(t, websocket.CloseMessage, m.msgType,
			"[ERROR] setelah channel ditutup, HANYA CloseMessage yang boleh terkirim, dapat: %d",
			m.msgType)
	}
}
func TestWritePump_Error_TickerStopsAfterReturn(t *testing.T) {
	// ARRANGE
	conn := &mockWSConn{}
	client, _ := newTestClient(conn)

	// ACT
	go client.WritePump()
	close(client.Send) // trigger WritePump untuk return

	// ASSERT — koneksi harus ditutup (defer berjalan = ticker.Stop() juga berjalan)
	waitFor(t, conn.isClosed, time.Second,
		"[ERROR] defer harus berjalan (ticker.Stop + conn.Close) setelah WritePump return")

	assert.True(t, conn.isClosed(),
		"[ERROR] Close() via defer harus dipanggil, ini bukti defer block berjalan dengan benar")

}
