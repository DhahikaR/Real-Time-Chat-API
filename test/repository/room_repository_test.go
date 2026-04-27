package repository_test

import (
	"context"
	"errors"
	"real-time-chat-api/internal/models/domain"
	"real-time-chat-api/internal/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupRoomTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestRoomRepository_Save_Success(t *testing.T) {
	gormDB, mock := setupRoomTestDB(t)
	repo := repository.NewRoomRepository(gormDB)

	roomID := uuid.New()
	ownerID := uuid.New()
	room := &domain.Room{
		ID:      roomID,
		Name:    "Test Room",
		Type:    "public",
		OwnerID: ownerID,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rooms"`)).
		WithArgs(
			sqlmock.AnyArg(), // name
			sqlmock.AnyArg(), // type
			sqlmock.AnyArg(), // owner_id
			sqlmock.AnyArg(), // id
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(room.ID))
	mock.ExpectCommit()

	err := repo.Save(context.Background(), room)
	assert.NoError(t, err)
}

func TestRoomRepository_Save_Error(t *testing.T) {
	gormDB, mock := setupRoomTestDB(t)
	repo := repository.NewRoomRepository(gormDB)

	roomID := uuid.New()
	ownerID := uuid.New()
	room := &domain.Room{
		ID:      roomID,
		Name:    "Test Room",
		Type:    "public",
		OwnerID: ownerID,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rooms"`)).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	err := repo.Save(context.Background(), room)
	assert.Error(t, err)
}

func TestRoomRepository_FindAll_Success(t *testing.T) {
	gormDB, mock := setupRoomTestDB(t)
	repo := repository.NewRoomRepository(gormDB)

	roomID1 := uuid.New()
	roomID2 := uuid.New()
	ownerID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "name", "type", "owner_id"}).
		AddRow(roomID1, "Room 1", "public", ownerID).
		AddRow(roomID2, "Room 2", "public", ownerID)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE type = $1`)).
		WithArgs("public").
		WillReturnRows(rows)

	rooms, err := repo.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, rooms, 2)
}

func TestRoomRepository_FindAll_Empty(t *testing.T) {
	gormDB, mock := setupRoomTestDB(t)
	repo := repository.NewRoomRepository(gormDB)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "owner_id"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE type = $1`)).
		WithArgs("public").
		WillReturnRows(rows)

	rooms, err := repo.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, rooms)
}

func TestRoomRepository_FindAll_Error(t *testing.T) {
	gormDB, mock := setupRoomTestDB(t)
	repo := repository.NewRoomRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE type = $1`)).
		WithArgs("public").
		WillReturnError(errors.New("database error"))

	rooms, err := repo.FindAll(context.Background())
	assert.Error(t, err)
	assert.Nil(t, rooms)
}

func TestRoomRepository_FindById_Success(t *testing.T) {
	gormDB, mock := setupRoomTestDB(t)
	repo := repository.NewRoomRepository(gormDB)

	roomID := uuid.New()
	ownerID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "name", "type", "owner_id"}).
		AddRow(roomID, "Test Room", "public", ownerID)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE id = $1`)).
		WithArgs(roomID, 1).
		WillReturnRows(rows)

	room, err := repo.FindById(context.Background(), roomID)
	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, roomID, room.ID)
}

func TestRoomRepository_FindById_NotFound(t *testing.T) {
	gormDB, mock := setupRoomTestDB(t)
	repo := repository.NewRoomRepository(gormDB)

	roomID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE id = $1`)).
		WithArgs(roomID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	room, err := repo.FindById(context.Background(), roomID)
	assert.NoError(t, err)
	assert.Nil(t, room)
}
