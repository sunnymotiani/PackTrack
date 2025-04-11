package models

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func (cfg RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

// OpenRedis initializes a Redis client
func OpenRedis(cfg RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: cfg.Password, // Empty if no password
		DB:       cfg.DB,       // Default DB is 0
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}
	return client, nil
}
