package services

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
)

type AuthRepository interface {
}

type Service struct {
	db     AuthRepository
	logger *slog.Logger
	router *chi.Mux
}

func NewService(conn AuthRepository, logger *slog.Logger) *Service {
	r := chi.NewMux()
	return &Service{
		db:     conn,
		logger: logger,
		router: r,
	}
}

func (s *Service) GetServiceMux() *chi.Mux {
	return s.router
}
