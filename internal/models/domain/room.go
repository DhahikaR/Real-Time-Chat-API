package domain

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4 (); primaryKey"`
	Name      string    `gorm:"not null"`
	Type      string    `gorm:"default:public"`
	OwnerID   uuid.UUID `gorm:"uuid"`
	createdAt time.Time `gorm:"created_at"`
}
