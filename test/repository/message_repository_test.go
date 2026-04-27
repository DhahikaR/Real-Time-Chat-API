package repository_test

import (
	"context"
	"errors"
	"real-time-chat-api/internal/models/domain"
	"real-time-chat-api/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMessageTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	return gormDB, mock
}

func TestMessageRepository_Save_Success(t *testing.T) {
	gormDB, mock := setupMessageTestDB(t)
	repo := repository.NewMessageRepository(gormDB)

	msgID := uuid.New()
	roomID := uuid.New()
	senderID := uuid.New()
	message := &domain.Message{
		ID:          msgID,
		RoomID:      roomID,
		SenderID:    senderID,
		Content:     "Hello World",
		SentAt:      time.Now(),
		IsDeletedAt: false,
	}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"messages\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(message.ID))
	mock.ExpectCommit()

	err := repo.Save(context.Background(), message)
	assert.NoError(t, err)
}

func TestMessageRepository_Save_Error(t *testing.T) {
	gormDB, mock := setupMessageTestDB(t)
	repo := repository.NewMessageRepository(gormDB)

	msgID := uuid.New()
	roomID := uuid.New()
	senderID := uuid.New()
	message := &domain.Message{
		ID:       msgID,
		RoomID:   roomID,
		SenderID: senderID,
		Content:  "Hello World",
		SentAt:   time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"messages\"").
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	err := repo.Save(context.Background(), message)
	assert.Error(t, err)
}

func TestMessageRepository_GetByRoom_Success(t *testing.T) {
	gormDB, mock := setupMessageTestDB(t)
	repo := repository.NewMessageRepository(gormDB)

	roomID := uuid.New()
	msgID1 := uuid.New()
	msgID2 := uuid.New()
	senderID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "room_id", "sender_id", "content", "sent_at", "is_deleted_at"}).
		AddRow(msgID1, roomID, senderID, "Message 1", now, false).
		AddRow(msgID2, roomID, senderID, "Message 2", now, false)

	// Use Any() to match the query pattern
	mock.ExpectQuery("SELECT \\* FROM \"messages\"").
		WillReturnRows(rows)

	messages, err := repo.GetByRoom(context.Background(), roomID, 1, 20)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
}

func TestMessageRepository_GetByRoom_Empty(t *testing.T) {
	gormDB, mock := setupMessageTestDB(t)
	repo := repository.NewMessageRepository(gormDB)

	roomID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "room_id", "sender_id", "content", "sent_at", "is_deleted_at"})

	mock.ExpectQuery("SELECT \\* FROM \"messages\"").
		WillReturnRows(rows)

	messages, err := repo.GetByRoom(context.Background(), roomID, 1, 20)
	assert.NoError(t, err)
	assert.Empty(t, messages)
}

func TestMessageRepository_GetByRoom_Error(t *testing.T) {
	gormDB, mock := setupMessageTestDB(t)
	repo := repository.NewMessageRepository(gormDB)

	roomID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM \"messages\"").
		WillReturnError(errors.New("database error"))

	messages, err := repo.GetByRoom(context.Background(), roomID, 1, 20)
	assert.Error(t, err)
	assert.Nil(t, messages)
}
