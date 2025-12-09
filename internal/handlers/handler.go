package handlers

import (
	"log/slog"

	"github.com/krishna102001/dependecy-injection/internal/services"
)

type Handler struct {
	svc    *services.Service
	logger *slog.Logger
}

func InitHandler(svc *services.Service, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}
