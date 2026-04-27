package exception_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"real-time-chat-api/internal/exception"
	"real-time-chat-api/internal/models/web"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupExceptionApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	})
}

func TestErrorHandler_FiberError_BadRequest(t *testing.T) {
	app := setupExceptionApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "bad request error")
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
	assert.Equal(t, "bad request error", result.Data)
}

func TestErrorHandler_FiberError_NotFound(t *testing.T) {
	app := setupExceptionApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "resource not found")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var result web.WebResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, "NOT FOUND", result.Status)
}

func TestErrorHandler_FiberError_InternalServerError(t *testing.T) {
	app := setupExceptionApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var result web.WebResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, "INTERNAL SERVER ERROR", result.Status)
}

func TestErrorHandler_GenericError(t *testing.T) {
	app := setupExceptionApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return errors.New("generic error occurred")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var result web.WebResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, "INTERNAL SERVER ERROR", result.Status)
	assert.Equal(t, fiber.StatusInternalServerError, result.Code)
	assert.Equal(t, "generic error occurred", result.Data)
}

func TestErrorHandler_FiberError_ZeroCode(t *testing.T) {
	app := setupExceptionApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return &fiber.Error{Code: 0, Message: "zero code error"}
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	// Zero code should default to 500
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
