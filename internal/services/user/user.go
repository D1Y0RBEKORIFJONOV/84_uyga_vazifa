package userservice

import (
	"context"
	"cors/internal/config"
	userentity "cors/internal/entity/user"
	userusecase "cors/internal/usecase/user"
	"cors/pkg/email"
	tokens "cors/pkg/token"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"time"
)

type User struct {
	logger *slog.Logger
	user   *userusecase.Repo
	cfg    *config.Config
}

func NewUser(logger *slog.Logger, user *userusecase.Repo, cfg *config.Config) *User {
	return &User{
		logger: logger,
		user:   user,
		cfg:    cfg,
	}
}

func (u *User) CreateUser(ctx context.Context, req *userentity.CreateUser) (*userentity.Status, error) {
	const op = "service.User.CreateUser"
	log := u.logger.With(
		slog.String(op, "creating user"))
	log.Info("Starting CreateUser")

	usr, err := u.user.GetUserOnMongoDb(ctx, "email", req.Email)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Error("err", err)
			return nil, err
		}
		log.Info("User does not exist")
	}
	if usr != nil {
		return nil, errors.New("user already exists")
	}

	defer log.Info("Ending CreateUser")
	secredCode, err := email.SenSecretCode([]string{req.Email})
	if err != nil {
		return nil, err
	}
	id := uuid.NewString()
	reqByte, err := json.Marshal(userentity.User{
		Email:      req.Email,
		Password:   req.Password,
		Role:       "user",
		LastName:   req.LastName,
		FirstName:  req.FirstName,
		UserID:     id,
		SecretCode: secredCode,
	})
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	err = u.user.Publish(reqByte, u.cfg.CreateUserTopic)
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	status := &userentity.Status{
		UserID:   id,
		Messages: []*userentity.Message{},
	}
	status.Messages = append(status.Messages, &userentity.Message{
		CreateAt: time.Now().Format("2006-01-02 15:04:05"),
		Status:   "User created sending",
	})
	err = u.user.CreateStatus(ctx, status)
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	return status, nil
}

func (u *User) VeryFyUser(ctx context.Context, req *userentity.VerifyRequest) (*userentity.Status, error) {
	const op = "service.User.VeryFyUser"
	log := u.logger.With(
		slog.String(op, "verifying user"))

	log.Info("Starting VerifyUser")
	defer log.Info("Ending VerifyUser")
	usr, err := u.user.GetUserOnRedis(ctx, req.Email)
	if err != nil {
		log.Error("err", err)

		return nil, err
	}
	if usr == nil {
		return nil, errors.New("secret code or email invalid")
	}
	if req.Secret != usr.SecretCode {
		return nil, errors.New("secret code error")
	}

	reqByte, err := json.Marshal(usr)
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	err = u.user.Publish(reqByte, u.cfg.VeryFyTopic)
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	err = u.user.UpdateStatus(ctx, &userentity.Message{
		CreateAt: time.Now().Format("2006-01-02 15:04:05"),
		Status:   "User verified sending",
	}, usr.UserID)
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	status := &userentity.Status{
		UserID:   usr.UserID,
		Messages: []*userentity.Message{},
	}
	status.Messages = append(status.Messages, &userentity.Message{
		CreateAt: time.Now().Format("2006-01-02 15:04:05"),
		Status:   "User verified sending",
	})

	return status, nil
}

func (u *User) LoginUser(ctx context.Context, req *userentity.LoginRequest) (*userentity.Token, error) {
	const op = "service.User.LoginUser"
	log := u.logger.With(
		slog.String(op, "logging in"))
	log.Info("Starting LoginUser")
	defer log.Info("Ending LoginUser")
	usr, err := u.user.GetUserOnMongoDb(ctx, "email", req.Email)
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	if usr == nil {
		return nil, errors.New("user does not exist")
	}
	if usr.Password != req.Password {
		log.Info("err", "Invalid password")
		return nil, errors.New("invalid password")
	}
	var token userentity.Token
	token.RefreshToken, err = tokens.NewRefreshToken(usr)
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	token.AccessToken, err = tokens.NewAccessToken(usr)
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	return &token, nil
}
