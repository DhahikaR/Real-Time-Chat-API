package service

import (
	"context"
	"real-time-chat-api/internal/models/web"
)

type AuthService interface {
	Register(ctx context.Context, request web.RegisterRequest) error
	Login(ctx context.Context, request web.LoginRequest) (string, error)
}
