package web

import "github.com/google/uuid"

type CreateRoomRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required,oneof=public private"`
	OwnerID uuid.UUID
}
