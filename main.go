package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/krishna102001/dependecy-injection/config"
	"github.com/krishna102001/dependecy-injection/internal/database"
	"github.com/krishna102001/dependecy-injection/internal/handlers"
	"github.com/krishna102001/dependecy-injection/internal/logger"
	"github.com/krishna102001/dependecy-injection/internal/services"
)

func main() {
	cfg, err := config.LoadConfigLocal()
	if err != nil {
		log.Fatalf(err.Error())
	}

	logger := logger.Initlogger("debug")
	svc := setupService(cfg, logger)

	setupServer(cfg, svc, logger)

}

func setupService(cfg *config.Config, logger *slog.Logger) *services.Service {
	conn, err := database.InitMongo(logger)
	if err != nil {
		logger.Error("failed to intialized the database")
		log.Fatal("failed to init mongo")
	}
	service := services.NewService(conn, logger)
	return service
}

func setupServer(cfg *config.Config, service *services.Service, logger *slog.Logger) *http.Server {
	if cfg == nil || cfg.Server.Port == "" {
		log.Fatalf("config file not loaded")
		return nil
	}

	handler := handlers.InitHandler(service, logger)

}
