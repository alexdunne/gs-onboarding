package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alexdunne/gs-onboarding/internal/api"
	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Addr        string
	DatabaseDSN string
}

func loadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read env file")
	}

	return &Config{
		Addr: viper.GetString("ADDR"),
		DatabaseDSN: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			viper.GetString("DATABASE_USER"),
			viper.GetString("DATABASE_PASSWORD"),
			viper.GetString("DATABASE_HOST"),
			viper.GetString("DATABASE_PORT"),
			viper.GetString("DATABASE_DB"),
		),
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
		log.Fatal(errors.Wrap(err, "opening store db connection"))
	}
	defer db.Close()

	router := echo.New()
	router.HideBanner = true
	router.Use(
		middleware.Recover(),
		middleware.Logger(),
	)

	h := api.Handler{
		DB: db,
	}

	router.GET("/all", h.HandleGetAllItems)
	router.GET("/stories", h.HandleGetStories)
	router.GET("/jobs", h.HandleGetJobs)

	if err := router.Start(cfg.Addr); err != nil {
		logger.Fatal("starting server", zap.Error(err))
	}
}
