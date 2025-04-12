package main

import (
	"os"

	"github.com/sunnymotiani/PackTrack/server/models"
)

type config struct {
	PSQL  models.PostgresConfig
	Redis models.RedisConfig
	CSRF  struct {
		Key    string
		Secure bool
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		Database: os.Getenv("PSQL_DATABASE"),
		Password: os.Getenv("PSQL_PASSWORD"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
		User:     os.Getenv("PSQL_USER"),
	}

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Redis = models.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}
	return cfg, nil
}
