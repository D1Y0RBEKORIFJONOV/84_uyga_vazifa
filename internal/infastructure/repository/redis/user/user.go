package userredis

import (
	"context"
	"cors/internal/config"
	userentity "cors/internal/entity/user"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisUserRepository struct {
	redisClient *redis.Client
}

func NewRedis(cfg *config.Config) *RedisUserRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUrl,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return &RedisUserRepository{
		redisClient: client,
	}
}

func (r *RedisUserRepository) SaveUserToRedis(ctx context.Context, user *userentity.User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := fmt.Sprintf(":%s", user.Email)
	err = r.redisClient.Set(ctx, key, string(userJson), time.Hour*5).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisUserRepository) GetUserOnRedis(ctx context.Context, email string) (*userentity.User, error) {
	key := fmt.Sprintf(":%s", email)
	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	user := &userentity.User{}
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
