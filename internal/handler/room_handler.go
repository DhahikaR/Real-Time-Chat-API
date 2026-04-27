package handler

import "github.com/gofiber/fiber/v2"

type RoomHandler interface {
	CreateRoom(c *fiber.Ctx) error
	ListRooms(c *fiber.Ctx) error
	GetMessages(c *fiber.Ctx) error
}
