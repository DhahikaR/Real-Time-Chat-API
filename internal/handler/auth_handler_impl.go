package handler

import (
	"real-time-chat-api/internal/helper"
	"real-time-chat-api/internal/models/web"
	"real-time-chat-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

// AuthHandlerImpl represents auth handler implementation
type AuthHandlerImpl struct {
	authService service.AuthService
}

// NewAuthHandler creates new auth handler
// @Summary Register new user
// @Description Register a new user to the system
// @Tags auth
// @Accept json
// @Produce json
// @Param RegisterRequest body web.RegisterRequest true "Register request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &AuthHandlerImpl{authService: authService}
}

// Register handles user registration
// @Summary Register new user
// @Description Register a new user to the system
// @Tags auth
// @Accept json
// @Produce json
// @Param RegisterRequest body web.RegisterRequest true "Register request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func (handler *AuthHandlerImpl) Register(c *fiber.Ctx) error {
	authRegisterRequest := web.RegisterRequest{}
	if err := helper.ReadFromRequestBody(c, &authRegisterRequest); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	if err := handler.authService.Register(c.Context(), authRegisterRequest); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	return helper.ResponseSuccess(c, "User registered successfully")
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param LoginRequest body web.LoginRequest true "Login request"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (handler *AuthHandlerImpl) Login(c *fiber.Ctx) error {
	authLoginRequest := web.LoginRequest{}
	if err := helper.ReadFromRequestBody(c, &authLoginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	token, err := handler.authService.Login(c.Context(), authLoginRequest)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
