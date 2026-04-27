package middleware_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"real-time-chat-api/internal/middleware"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-jwt-secret"

func generateTestToken(userID, email string, secret string, expiry time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(expiry).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func setupMiddlewareApp() *fiber.App {
	app := fiber.New()
	protected := app.Group("/protected", middleware.JWTMiddleware(testSecret))
	protected.Get("/resource", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		email := c.Locals("email").(string)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user_id": userID,
			"email":   email,
		})
	})
	return app
}

func TestJWTMiddleware_Success(t *testing.T) {
	app := setupMiddlewareApp()

	userID := uuid.New().String()
	email := "john@example.com"
	token := generateTestToken(userID, email, testSecret, 24*time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/protected/resource", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, userID, result["user_id"])
	assert.Equal(t, email, result["email"])
}

func TestJWTMiddleware_MissingAuthorizationHeader(t *testing.T) {
	app := setupMiddlewareApp()

	req := httptest.NewRequest(http.MethodGet, "/protected/resource", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var result map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, "UNAUTHORIZED", result["status"])
	assert.Contains(t, result["message"], "missing authorized header")
}

func TestJWTMiddleware_InvalidAuthorizationFormat(t *testing.T) {
	app := setupMiddlewareApp()

	req := httptest.NewRequest(http.MethodGet, "/protected/resource", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var result map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Contains(t, result["message"], "invalid authorization format")
}

func TestJWTMiddleware_MissingBearerPrefix(t *testing.T) {
	app := setupMiddlewareApp()

	req := httptest.NewRequest(http.MethodGet, "/protected/resource", nil)
	req.Header.Set("Authorization", "Token sometoken")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	app := setupMiddlewareApp()

	req := httptest.NewRequest(http.MethodGet, "/protected/resource", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var result map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Contains(t, result["message"], "invalid or expired token")
}

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
	app := setupMiddlewareApp()

	userID := uuid.New().String()
	email := "john@example.com"
	// Generate token that expired 1 hour ago
	token := generateTestToken(userID, email, testSecret, -1*time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/protected/resource", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJWTMiddleware_WrongSecret(t *testing.T) {
	app := setupMiddlewareApp()

	userID := uuid.New().String()
	email := "john@example.com"
	// Sign with different secret
	token := generateTestToken(userID, email, "wrong-secret", 24*time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/protected/resource", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJWTMiddleware_TokenMissingUserID(t *testing.T) {
	app := setupMiddlewareApp()

	// Token without user_id claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "john@example.com",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(testSecret))

	req := httptest.NewRequest(http.MethodGet, "/protected/resource", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var result map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Contains(t, result["message"], "invalid user id in token")
}

func TestJWTMiddleware_TokenMissingEmail(t *testing.T) {
	app := setupMiddlewareApp()

	// Token without email claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uuid.New().String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(testSecret))

	req := httptest.NewRequest(http.MethodGet, "/protected/resource", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var result map[string]string
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.Contains(t, result["message"], "invalid email in token")
}
