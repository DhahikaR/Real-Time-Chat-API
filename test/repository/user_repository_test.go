package repository_test

import (
	"context"
	"errors"
	"real-time-chat-api/internal/models/domain"
	"real-time-chat-api/internal/repository"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupUserTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestUserRepository_Save_Success(t *testing.T) {
	gormDB, mock := setupUserTestDB(t)
	repo := repository.NewUserRepository(gormDB)

	userID := uuid.New()
	user := &domain.User{
		ID:           userID,
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(
			sqlmock.AnyArg(), // name
			sqlmock.AnyArg(), // email
			sqlmock.AnyArg(), // password_hash
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // id
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))
	mock.ExpectCommit()

	err := repo.Save(context.Background(), user)
	assert.NoError(t, err)
}

func TestUserRepository_Save_Error(t *testing.T) {
	gormDB, mock := setupUserTestDB(t)
	repo := repository.NewUserRepository(gormDB)

	userID := uuid.New()
	user := &domain.User{
		ID:           userID,
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))
	mock.ExpectRollback()

	err := repo.Save(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
}

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	gormDB, mock := setupUserTestDB(t)
	repo := repository.NewUserRepository(gormDB)

	userID := uuid.New()
	email := "john@example.com"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "created_at"}).
		AddRow(userID, "John Doe", email, "hashedpassword", now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	user, err := repo.FindByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	gormDB, mock := setupUserTestDB(t)
	repo := repository.NewUserRepository(gormDB)

	email := "notfound@example.com"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
		WithArgs(email, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByEmail(context.Background(), email)
	assert.Nil(t, err)
	assert.Nil(t, user)

}

func TestUserRepository_FindById_Success(t *testing.T) {
	gormDB, mock := setupUserTestDB(t)
	repo := repository.NewUserRepository(gormDB)

	userID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "created_at"}).
		AddRow(userID, "John Doe", "john@example.com", "hashedpassword", now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
		WithArgs(userID, 1).
		WillReturnRows(rows)

	user, err := repo.FindById(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
}

func TestUserRepository_FindById_NotFound(t *testing.T) {
	gormDB, mock := setupUserTestDB(t)
	repo := repository.NewUserRepository(gormDB)

	userID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
		WithArgs(userID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindById(context.Background(), userID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	_ = user
}
