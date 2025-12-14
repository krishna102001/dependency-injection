package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/krishna102001/dependecy-injection/internal/models"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertUser(ctx context.Context, user models.User) (string, error)
	// GetUserByPhone(ctx context.Context, countryCode, phone string) (*models.User, error)
}

type Service struct {
	db     AuthRepository
	logger *slog.Logger
	router *chi.Mux
}

func NewService(authDB AuthRepository, logger *slog.Logger) *Service {
	r := chi.NewMux()
	return &Service{
		db:     authDB,
		logger: logger,
		router: r,
	}
}

func (s *Service) RegisterUser(ctx context.Context, reqBody models.User) (string, error) {
	if reqBody.Email == "" || reqBody.Phone == "" {
		return "", fmt.Errorf("email/phone is empty")
	}
	// db call check user exist or not ?
	_, err := s.db.GetUserByEmail(ctx, reqBody.Email)
	if err != nil && !errors.Is(err, models.ErrNoDataFound) {
		return "", fmt.Errorf("failed to get the user from db %v", err)
	}
	if err == nil {
		return "", fmt.Errorf("user already exist")
	}
	// if not present then make the db insert operations

	userId, err := s.db.InsertUser(ctx, reqBody)
	if err != nil {
		return "", fmt.Errorf("error in the database %v", err)
	}

	// return userId
	return userId, nil
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get the user %v", err)
	}
	if user.Password != password {
		return nil, fmt.Errorf("unauthorized access")
	}

	return user, nil
}

func (s *Service) GetServiceMux() *chi.Mux {
	return s.router
}
