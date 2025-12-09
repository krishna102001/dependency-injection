package services

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Service struct {
	db     *mongo.Client
	logger *slog.Logger
}

func NewService(conn *mongo.Client, logger *slog.Logger) *Service {

	return &Service{
		db:     conn,
		logger: logger,
	}
}
