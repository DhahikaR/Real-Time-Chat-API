package middleware

import (
	"real-time-chat-api/internal/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHandler := c.Get("Authorization")
		if authHandler == "" {
			return helper.Unauthorized(c, "missing authorized header")
		}

		fields := strings.Split(authHandler, " ")
		if len(fields) != 2 || fields[0] != "Bearer" {
			return helper.Unauthorized(c, "invalid authorization format")
		}

		tokenString := fields[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token signature")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return helper.Unauthorized(c, "invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return helper.Unauthorized(c, "invalid token claims")
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			return helper.Unauthorized(c, "invalid user id in token")
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			return helper.Unauthorized(c, "invalid email in token")
		}

		c.Locals("user_id", userId)
		c.Locals("email", email)

		return c.Next()
	}
}
