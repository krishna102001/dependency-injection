package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/krishna102001/dependecy-injection/config"
	"github.com/krishna102001/dependecy-injection/internal/database"
	"github.com/krishna102001/dependecy-injection/internal/handlers"
	"github.com/krishna102001/dependecy-injection/internal/logger"
	"github.com/krishna102001/dependecy-injection/internal/services"
	"github.com/krishna102001/dependecy-injection/routes"
	"github.com/rs/cors"
)

func main() {
	loadEnv()

	cfg, err := config.GetConfig()
	if err != nil || cfg.Server == nil {
		log.Fatalf("server has yaml error %v", err)
	}

	logger := logger.Initlogger("debug")
	svc := setupService(cfg, logger)

	server := setupServer(cfg, svc, logger)

	if err := server.ListenAndServe(); err == nil {
		log.Println("Successfully started the server")
	}

}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warn: No env file found")
	}
	cfg, err := config.LoadConfigLocal()
	if err != nil || cfg == nil {
		log.Fatalf("failed to load the config file %v", err.Error())
	}

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

	authHandler := handlers.InitHandler(service, logger)

	handler := routes.SetupRoutes(authHandler, logger, service.GetServiceMux())
	handlerWithCors := setupCors(handler)
	server := http.Server{
		Handler: handlerWithCors,
		Addr:    ":" + cfg.Server.Port,
	}
	return &server
}

func setupCors(handler http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"*",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})
	return c.Handler(handler)
}
