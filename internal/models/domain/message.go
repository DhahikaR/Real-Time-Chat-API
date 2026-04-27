package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4 ();primaryKey"`
	RoomID      uuid.UUID `gorm:"type:uuid;index"`
	SenderID    uuid.UUID `gorm:"type:uuid;index"`
	Content     string    `gorm:"not null"`
	SentAt      time.Time
	IsDeletedAt bool `gorm:"default:false"`
}
