package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"real-time-chat-api/internal/handler"
	"real-time-chat-api/internal/models/web"
	"real-time-chat-api/test/mocks"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAuthApp(authService *mocks.MockAuthService) *fiber.App {
	app := fiber.New()
	authHandler := handler.NewAuthHandler(authService)
	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)
	return app
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	app := setupAuthApp(mockService)

	reqBody := web.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	mockService.On("Register", mock.Anything, reqBody).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result web.WebResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &result)
	assert.Equal(t, "SUCCESS", result.Status)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Register_ServiceError(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	app := setupAuthApp(mockService)

	reqBody := web.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	mockService.On("Register", mock.Anything, reqBody).Return(errors.New("email already exist"))

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	app := setupAuthApp(mockService)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	app := setupAuthApp(mockService)

	reqBody := web.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	mockService.On("Login", mock.Anything, reqBody).Return("jwt-token-string", nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &result)
	assert.Equal(t, "jwt-token-string", result["token"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_Unauthorized(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	app := setupAuthApp(mockService)

	reqBody := web.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	mockService.On("Login", mock.Anything, reqBody).Return("", errors.New("wrong email or password"))

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	app := setupAuthApp(mockService)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
