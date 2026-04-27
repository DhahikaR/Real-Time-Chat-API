package config_test

import (
	"real-time-chat-api/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB_InvalidConnection(t *testing.T) {
	cfg := config.Config{
		DBHost:     "invalid-host-that-does-not-exist",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "password",
		DBName:     "testdb",
	}

	db, err := config.NewDB(cfg)

	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, db)
	} else {
		assert.NotNil(t, db)
	}
}

func TestNewDB_DSNFormat(t *testing.T) {
	cfg := config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "secret",
		DBName:     "chatdb",
	}

	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "secret", cfg.DBPassword)
	assert.Equal(t, "chatdb", cfg.DBName)
}

func TestNewDB_EmptyConfig(t *testing.T) {
	cfg := config.Config{}

	db, err := config.NewDB(cfg)
	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, db)
	}
}
