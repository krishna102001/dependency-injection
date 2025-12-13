package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/krishna102001/dependecy-injection/config"
	"github.com/krishna102001/dependecy-injection/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AuthDB struct {
	dbName      string
	mongoURI    string
	usersColl   string
	mongoClient *mongo.Client
	logger      *slog.Logger
}

func GetMongoDB(logger *slog.Logger) (*AuthDB, error) {
	cfg, err := config.GetConfig()
	if err != nil || cfg.Mongo == nil {
		logger.Warn("config is not loaded", "error", err)
		return nil, fmt.Errorf("config is not loaded")
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

	return &AuthDB{
		dbName:      cfg.Mongo.DBName,
		mongoURI:    cfg.Mongo.URI,
		usersColl:   cfg.Mongo.Collection,
		mongoClient: client,
		logger:      logger,
	}, nil
}

func (db *AuthDB) GetUserCollection() *mongo.Collection {
	return db.mongoClient.Database(db.dbName).Collection(db.usersColl)
}

func (db *AuthDB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	coll := db.GetUserCollection()

	result := coll.FindOne(ctx, bson.M{"email": email})
	var user models.User
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrNoDataFound
		}
		return nil, fmt.Errorf("failed to decode the user data %w", err)
	}
	log.Printf("user %+v", user)
	return &user, nil
}

func (db *AuthDB) InsertUser(ctx context.Context, user models.User) (string, error) {
	coll := db.GetUserCollection()

	result, err := coll.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to insert user in db %v", err)
	}
	return result.InsertedID.(bson.ObjectID).Hex(), nil
}
