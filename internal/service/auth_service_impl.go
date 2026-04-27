package service

import (
	"context"
	"errors"
	"real-time-chat-api/internal/models/domain"
	"real-time-chat-api/internal/models/web"
	"real-time-chat-api/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepository repository.UserRepository
	jwtSecret      string
}

func NewAuthService(repository repository.UserRepository, secret string) AuthService {
	return &AuthServiceImpl{
		userRepository: repository,
		jwtSecret:      secret,
	}
}

func (service *AuthServiceImpl) Register(ctx context.Context, request web.RegisterRequest) error {
	existing, err := service.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already exist")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(hashed),
	}

	return service.userRepository.Save(ctx, user)
}

func (service *AuthServiceImpl) Login(ctx context.Context, request web.LoginRequest) (string, error) {
	user, err := service.userRepository.FindByEmail(ctx, request.Email)
	if err != nil || user == nil {
		return "", errors.New("wrong email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return "", errors.New("wrong email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(service.jwtSecret))
}
