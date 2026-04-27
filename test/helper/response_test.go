package helper_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"real-time-chat-api/internal/helper"
	"real-time-chat-api/internal/models/web"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestBadRequest(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return helper.BadRequest(c, "bad request message")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result web.WebResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, "BAD REQUEST", result.Status)
	assert.Equal(t, fiber.StatusBadRequest, result.Code)
	assert.Equal(t, "bad request message", result.Data)
}

func TestResponseSuccess(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return helper.ResponseSuccess(c, "operation successful")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result web.WebResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, "SUCCESS", result.Status)
	assert.Equal(t, fiber.StatusOK, result.Code)
	assert.Equal(t, "operation successful", result.Data)
}

func TestInternalServerError(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return helper.InternalServerError(c, "internal error occurred")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var result web.WebResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, "INTERNAL SEVER ERROR", result.Status)
	assert.Equal(t, fiber.StatusInternalServerError, result.Code)
}

func TestForbidden(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return helper.Forbidden(c, "access denied")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	var result web.WebResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, "FORBIDDEN", result.Status)
	assert.Equal(t, fiber.StatusForbidden, result.Code)
	assert.Equal(t, "access denied", result.Data)
}

func TestUnauthorized(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return helper.Unauthorized(c, "unauthorized access")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var result map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, "UNAUTHORIZED", result["status"])
	assert.Equal(t, "unauthorized access", result["message"])
}
