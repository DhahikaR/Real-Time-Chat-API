package exception

import (
	"real-time-chat-api/internal/models/web"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if fiberErr, ok := err.(*fiber.Error); ok {
		code := fiberErr.Code
		if code == 0 {
			code = fiber.StatusInternalServerError
		}
		statusText := "ERROR"
		if code == fiber.StatusBadRequest {
			statusText = "BAD REQUEST"
		} else if code == fiber.StatusNotFound {
			statusText = "NOT FOUND"
		} else if code == fiber.StatusInternalServerError {
			statusText = "INTERNAL SERVER ERROR"
		}

		return c.Status(code).JSON(web.WebResponse{
			Code:   code,
			Status: statusText,
			Data:   fiberErr.Message,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
		Code:   fiber.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data:   err.Error(),
	})
}
