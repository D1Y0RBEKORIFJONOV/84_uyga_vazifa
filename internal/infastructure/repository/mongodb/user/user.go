package usermongodb

import (
	"context"
	"cors/internal/config"
	userentity "cors/internal/entity/user"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"net/url"
)

type MongoDB struct {
	mongoClient      *mongo.Client
	db               *mongo.Database
	collectionUser   *mongo.Collection
	collectionStatus *mongo.Collection
	logger           *slog.Logger
}

func NewMongoDB(cfg *config.Config, logger *slog.Logger) (*MongoDB, error) {
	uri := url.URL{
		Scheme: "mongodb",
		Host:   fmt.Sprintf("%s%s", cfg.DB.Host, cfg.DB.Port),
	}

	clientOptions := options.Client().ApplyURI(uri.String())

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("could not connect to MongoDB: %w", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not ping MongoDB: %w", err)
	}

	return &MongoDB{
		mongoClient:      client,
		db:               client.Database(cfg.DB.Name),
		collectionStatus: client.Database(cfg.DB.Name).Collection("status"),
		collectionUser:   client.Database(cfg.DB.Name).Collection("user"),
		logger:           logger,
	}, nil
}

func (m *MongoDB) CreateStatus(ctx context.Context, user *userentity.Status) error {
	const op = "MongoDB.CreateStatus"
	log := m.logger.With(slog.String("method", op))
	log.Info("Starting")
	defer log.Info("Ending")
	_, err := m.collectionUser.InsertOne(ctx, user)
	if err != nil {
		log.Error("err", err)
		return fmt.Errorf("%w", err)
	}
	return nil
}
func (m *MongoDB) SaveUserToMongo(ctx context.Context, user *userentity.User) error {
	const op = "MongoDB.SaveUserToMongo"
	log := m.logger.With(slog.String("method", op))
	log.Info("Starting")
	defer log.Info("Ending")
	_, err := m.collectionUser.InsertOne(ctx, user)
	if err != nil {
		log.Error("err", err)
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (m *MongoDB) UpdateStatus(ctx context.Context, message *userentity.Message, userId string) error {
	const op = "MongoDB.UpdateStatus"
	log := m.logger.With(slog.String("method", op))
	log.Info("Starting")
	defer log.Info("Ending")
	update := bson.M{
		"$push": bson.M{
			"messages": message,
		},
	}
	filter := bson.M{
		"user_id": userId,
	}
	_, err := m.collectionStatus.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error("err", err)
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (m *MongoDB) GetUserOnMongoDb(ctx context.Context, field, value string) (*userentity.User, error) {
	const op = "MongoDB.GetUserOnMongoDb"
	log := m.logger.With(slog.String("method", op))
	log.Info("Starting")
	defer log.Info("Ending")
	filter := bson.M{field: value}
	var user userentity.User
	err := m.collectionUser.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Error("err", err)
		return nil, fmt.Errorf("%w", err)
	}
	return &user, nil
}
