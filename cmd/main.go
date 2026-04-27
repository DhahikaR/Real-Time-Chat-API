package main

import (
	"log"
	"real-time-chat-api/internal/config"
	"real-time-chat-api/internal/exception"
	"real-time-chat-api/internal/handler"
	"real-time-chat-api/internal/middleware"
	"real-time-chat-api/internal/repository"
	"real-time-chat-api/internal/service"
	"real-time-chat-api/internal/ws"
	"real-time-chat-api/pkg/pubsub"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// @title           Real Time Chat API + WebSocket
// @version         1.0
// @description     Documentation API for Real Time Chat using Fiber + WebSocket
// @host            localhost:8080
// @BasePath        /
// @termsOfService  http://swagger.io/terms/

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, use env OS")
	}

	cfg := config.Load()

	db, err := config.NewDB(cfg)
	if err != nil {
		log.Fatal("Database connection failed", err)
	}

	rdb := pubsub.NewRedisClient(cfg.RedisAddr)

	userRepository := repository.NewUserRepository(db)
	roomRepository := repository.NewRoomRepository(db)
	messageRepository := repository.NewMessageRepository(db)

	authService := service.NewAuthService(userRepository, cfg.JWTSecret)
	roomService := service.NewRoomService(roomRepository, messageRepository)

	hub := ws.NewHub(rdb)
	go hub.Run()

	authHandler := handler.NewAuthHandler(authService)
	roomHandler := handler.NewRoomHandler(roomService)
	wsHandler := handler.NewWSHandler(hub)

	router := fiber.New(fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	})

	// Public routes
	router.Post("/auth/register", authHandler.Register)
	router.Post("/auth/login", authHandler.Login)

	// Protected routes with JWT middleware
	auth := router.Group("/", middleware.JWTMiddleware(cfg.JWTSecret))
	auth.Get("/ws/connect", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}, wsHandler.Connect())
	auth.Post("/rooms", roomHandler.CreateRoom)
	auth.Get("/rooms", roomHandler.ListRooms)
	auth.Get("/rooms/:id/message", roomHandler.GetMessages)

	log.Println("Server Running on Port", cfg.AppPort)

	if err := router.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
