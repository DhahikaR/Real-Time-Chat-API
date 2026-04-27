package config

import (
	"log"
	"os"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisAddr  string
	JWTSecret  string
}

func Load() Config {
	return Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: requireENV("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME", "chatdb"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:  requireENV("JWT_SECRET"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func requireENV(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("❌ Environment variable %q is required but not set", key)
	}
	return v
}
