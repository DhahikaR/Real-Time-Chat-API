package helper_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"real-time-chat-api/internal/helper"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type TestPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestReadFromRequestBody_Success(t *testing.T) {
	app := fiber.New()
	app.Post("/test", func(c *fiber.Ctx) error {
		var payload TestPayload
		err := helper.ReadFromRequestBody(c, &payload)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(payload)
	})

	body := bytes.NewBufferString(`{"name":"John","email":"john@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/test", body)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestReadFromRequestBody_InvalidJSON(t *testing.T) {
	app := fiber.New()
	app.Post("/test", func(c *fiber.Ctx) error {
		var payload TestPayload
		err := helper.ReadFromRequestBody(c, &payload)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(payload)
	})

	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/test", body)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestReadFromRequestBody_EmptyBody(t *testing.T) {
	app := fiber.New()
	app.Post("/test", func(c *fiber.Ctx) error {
		var payload TestPayload
		err := helper.ReadFromRequestBody(c, &payload)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(payload)
	})

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	// Fiber's BodyParser returns an error for empty body with JSON content type
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
