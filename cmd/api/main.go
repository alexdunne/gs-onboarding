package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alexdunne/gs-onboarding/internal/api"
	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Port        int
	DatabaseDSN string
	RedisURL    string
}

func loadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read env file")
	}

	return &Config{
		Port: viper.GetInt("API_PORT"),
		DatabaseDSN: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			viper.GetString("DATABASE_USER"),
			viper.GetString("DATABASE_PASSWORD"),
			viper.GetString("DATABASE_HOST"),
			viper.GetString("DATABASE_PORT"),
			viper.GetString("DATABASE_DB"),
		),
		RedisURL: viper.GetString("REDIS_URL"),
	}, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "loading config"))
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(errors.Wrap(err, "creating logger"))
	}
	defer logger.Sync()

	ctx := context.Background()

	db, err := database.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(errors.Wrap(err, "opening database connection"))
	}
	defer db.Close()

	cache, err := api.NewCache(ctx, cfg.RedisURL, db, logger)
	if err != nil {
		log.Fatal(errors.Wrap(err, "opening cache connection"))
	}
	defer cache.Close()

	h := api.Handler{
		Cache: cache,
	}

	s := api.NewServer(cfg.Port, logger, h)
	if err := s.Start(); err != nil {
		logger.Fatal("running server", zap.Error(err))
	}
}
