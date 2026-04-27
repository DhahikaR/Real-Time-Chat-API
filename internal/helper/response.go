package helper

import (
	"real-time-chat-api/internal/models/web"

	"github.com/gofiber/fiber/v2"
)

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
		Code:   fiber.StatusBadRequest,
		Status: "BAD REQUEST",
		Data:   message,
	})
}

func ResponseSuccess(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:   fiber.StatusOK,
		Status: "SUCCESS",
		Data:   message,
	})
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
		Code:   fiber.StatusInternalServerError,
		Status: "INTERNAL SEVER ERROR",
		Data:   message,
	})
}

func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(web.WebResponse{
		Code:   fiber.StatusForbidden,
		Status: "FORBIDDEN",
		Data:   message,
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status":  "UNAUTHORIZED",
		"message": message,
	})
}
