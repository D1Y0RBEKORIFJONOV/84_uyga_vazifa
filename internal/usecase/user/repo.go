package userusecase

import (
	"context"
	userentity "cors/internal/entity/user"
)

type (
	Saver interface {
		CreateStatus(ctx context.Context, user *userentity.Status) error
		SaveUserToMongoDB(ctx context.Context, user *userentity.User) error
		SaveUserToRedis(ctx context.Context, user *userentity.User) error
	}
	Updater interface {
		UpdateStatus(ctx context.Context, message *userentity.Message,userId string) error
	}
	Broker interface {
		Publish(value []byte, topicKey string) error
	}
	Provider interface {
		GetUserOnMongoDb(ctx context.Context, field, value string) (*userentity.User, error)
		GetUserOnRedis(ctx context.Context, email string) (*userentity.User, error)
	}
)

type Repo struct {
	saver    Saver
	broker   Broker
	updater  Updater
	provider Provider
}

func NewRepo(saver Saver, broker Broker, updater Updater, provider Provider) *Repo {
	return &Repo{
		saver:    saver,
		broker:   broker,
		updater:  updater,
		provider: provider,
	}
}

func (r *Repo) GetUserOnRedis(ctx context.Context, email string) (*userentity.User, error){
	return r.provider.GetUserOnRedis(ctx, email)
}

func (r *Repo) CreateStatus(ctx context.Context, user *userentity.Status) error {
	return r.saver.CreateStatus(ctx, user)
}

func (r *Repo) UpdateStatus(ctx context.Context,message *userentity.Message,userId string) error {
	return r.updater.UpdateStatus(ctx, message,userId)
}

func (r *Repo) Publish(value []byte, topicKey string) error {
	return r.broker.Publish(value, topicKey)
}

func (r *Repo) GetUserOnMongoDb(ctx context.Context, field, value string) (*userentity.User, error) {
	return r.provider.GetUserOnMongoDb(ctx, field, value)
}


func (r *Repo)SaveUserToMongoDB(ctx context.Context, user *userentity.User) error {
	return r.saver.SaveUserToMongoDB(ctx, user)
}

func (r *Repo)SaveUserToRedis(ctx context.Context, user *userentity.User) error {
	return r.saver.SaveUserToRedis(ctx, user)
}