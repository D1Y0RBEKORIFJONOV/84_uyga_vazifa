package config

import (
	"os"
	"time"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	HTTPUrl     string
	RedisUrl    string
	Context     struct {
		Timeout string
	}
	Token struct {
		Secret     string
		AccessTTL  time.Duration
		RefreshTTL time.Duration
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
	}
	KafkaUrl        string
	CreateUserTopic string
	VeryFyTopic     string
}

func Token() string {
	c := Config{}
	c.Token.Secret = getEnv("TOKEN_SECRET", "token_secret")
	return c.Token.Secret
}

func New() *Config {
	var config Config

	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "local")
	config.HTTPUrl = getEnv("RPC_PORT", "93_cors_service:7777")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	config.DB.Host = getEnv("MONGO_HOST", "mongo")
	config.DB.Port = getEnv("MONGO_PORT", ":27017")
	config.DB.User = getEnv("MONGO_USER", "")
	config.DB.Password = getEnv("MONGO_PASSWORD", "")
	config.DB.Name = getEnv("MONGO_DATABASE", "cors")

	config.Token.Secret = getEnv("TOKEN_SECRET", "D1YORTOP4EEK")
	accessTTl, err := time.ParseDuration(getEnv("TOKEN_ACCESS_TTL", "1h"))
	if err != nil {
		return nil
	}
	refreshTTL, err := time.ParseDuration(getEnv("TOKEN_REFRESH_TTL", "24h"))

	if err != nil {
		return nil
	}
	config.Token.AccessTTL = accessTTl
	config.Token.RefreshTTL = refreshTTL

	config.KafkaUrl = getEnv("KAFKA_URL", "broker:9092")
	config.CreateUserTopic = getEnv("CREATE_USER_TOPIC", "USER-CREATE")
	config.VeryFyTopic = getEnv("VEYFY_TOPIC", "VFY-VFY")
	config.RedisUrl = getEnv("REDIS_URL", "redis:6379")
	return &config
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
