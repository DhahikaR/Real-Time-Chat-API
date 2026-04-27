package config

import (
	"fmt"
	"real-time-chat-api/internal/models/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // log semua query
	})
	if err != nil {
		return nil, fmt.Errorf("gagal koneksi database: %w", err)
	}

	// Aktifkan UUID extension di PostgreSQL
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	// AutoMigrate: buat/update tabel otomatis sesuai struct model
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Room{},
		&domain.Message{},
	); err != nil {
		return nil, fmt.Errorf("AutoMigrate gagal: %w", err)
	}

	return db, nil
}
