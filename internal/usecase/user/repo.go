package userusecase

import (
	"context"
	userentity "cors/internal/entity/user"
)

type (
	MongoSaver interface {
		CreateStatus(ctx context.Context, user *userentity.Status) error
		SaveUserToMongo(ctx context.Context, user *userentity.User) error
	}
	MongoUpdater interface {
		UpdateStatus(ctx context.Context, message *userentity.Message, userId string) error
	}
	Redis interface {
		SaveUserToRedis(ctx context.Context, user *userentity.User) error
		GetUserOnRedis(ctx context.Context, email string) (*userentity.User, error)
	}
	Broker interface {
		Publish(value []byte, topicKey string) error
	}
	Provider interface {
		GetUserOnMongoDb(ctx context.Context, field, value string) (*userentity.User, error)
	}
)

type Repo struct {
	saver    MongoSaver
	broker   Broker
	updater  MongoUpdater
	provider Provider
	redis    Redis
}

func NewRepo(saver MongoSaver, broker Broker, updater MongoUpdater, provider Provider, redis Redis) *Repo {
	return &Repo{
		saver:    saver,
		broker:   broker,
		updater:  updater,
		provider: provider,
		redis:    redis,
	}
}

func (r *Repo) GetUserOnRedis(ctx context.Context, email string) (*userentity.User, error) {
	return r.redis.GetUserOnRedis(ctx, email)
}

func (r *Repo) CreateStatus(ctx context.Context, user *userentity.Status) error {
	return r.saver.CreateStatus(ctx, user)
}

func (r *Repo) UpdateStatus(ctx context.Context, message *userentity.Message, userId string) error {
	return r.updater.UpdateStatus(ctx, message, userId)
}

func (r *Repo) Publish(value []byte, topicKey string) error {
	return r.broker.Publish(value, topicKey)
}

func (r *Repo) GetUserOnMongoDb(ctx context.Context, field, value string) (*userentity.User, error) {
	return r.provider.GetUserOnMongoDb(ctx, field, value)
}

func (r *Repo) SaveUserToMongo(ctx context.Context, user *userentity.User) error {
	return r.saver.SaveUserToMongo(ctx, user)
}

func (r *Repo) SaveUserToRedis(ctx context.Context, user *userentity.User) error {
	return r.redis.SaveUserToRedis(ctx, user)
}
