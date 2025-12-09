package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/krishna102001/dependecy-injection/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongo(logger *slog.Logger) (*mongo.Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Warn("config is not loaded", "error", err)
		return nil, fmt.Errorf("config is not loaded")
	}
	if cfg.Mongo == nil || cfg.Mongo.URI == "" || cfg.Mongo.DBName == "" || cfg.Mongo.Collection == "" {
		logger.Error("failed to get the config value of mongo")
		return nil, fmt.Errorf("config value is empty")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		logger.Error("failed to connect the mongodb", "error", err)
		return nil, fmt.Errorf("failed to connect to mongo")
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		logger.Error("Failed to connect with mongodb")
		return nil, fmt.Errorf("failed to ping the mongo")
	}

	logger.Info("Mongodb connected successfull")

	return client, nil
}
