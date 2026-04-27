package config_test

import (
	"os"
	"real-time-chat-api/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load_DefaultValues(t *testing.T) {
	// Clear any existing env vars to test defaults
	os.Unsetenv("APP_PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("JWT_SECRET")

	cfg := config.Load()

	assert.Equal(t, "8080", cfg.AppPort)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "secret", cfg.DBPassword)
	assert.Equal(t, "chatdb", cfg.DBName)
	assert.Equal(t, "localhost:6379", cfg.RedisAddr)
	assert.Equal(t, "secret", cfg.JWTSecret)
}

func TestConfig_Load_FromEnvironment(t *testing.T) {
	// Set custom env vars
	os.Setenv("APP_PORT", "9090")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASSWORD", "supersecret")
	os.Setenv("DB_NAME", "mydb")
	os.Setenv("REDIS_ADDR", "redis.example.com:6380")
	os.Setenv("JWT_SECRET", "my-jwt-secret")

	defer func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("REDIS_ADDR")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg := config.Load()

	assert.Equal(t, "9090", cfg.AppPort)
	assert.Equal(t, "db.example.com", cfg.DBHost)
	assert.Equal(t, "5433", cfg.DBPort)
	assert.Equal(t, "admin", cfg.DBUser)
	assert.Equal(t, "supersecret", cfg.DBPassword)
	assert.Equal(t, "mydb", cfg.DBName)
	assert.Equal(t, "redis.example.com:6380", cfg.RedisAddr)
	assert.Equal(t, "my-jwt-secret", cfg.JWTSecret)
}

func TestConfig_Load_PartialEnvironment(t *testing.T) {
	// Only set some env vars
	os.Setenv("APP_PORT", "3000")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")

	defer func() {
		os.Unsetenv("APP_PORT")
	}()

	cfg := config.Load()

	assert.Equal(t, "3000", cfg.AppPort)
	assert.Equal(t, "localhost", cfg.DBHost) // default
	assert.Equal(t, "5432", cfg.DBPort)      // default
}

func TestConfig_Struct_Fields(t *testing.T) {
	cfg := config.Config{
		AppPort:    "8080",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "password",
		DBName:     "testdb",
		RedisAddr:  "localhost:6379",
		JWTSecret:  "secret",
	}

	assert.NotEmpty(t, cfg.AppPort)
	assert.NotEmpty(t, cfg.DBHost)
	assert.NotEmpty(t, cfg.DBPort)
	assert.NotEmpty(t, cfg.DBUser)
	assert.NotEmpty(t, cfg.DBPassword)
	assert.NotEmpty(t, cfg.DBName)
	assert.NotEmpty(t, cfg.RedisAddr)
	assert.NotEmpty(t, cfg.JWTSecret)
}
