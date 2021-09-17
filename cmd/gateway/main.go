package main

import (
	"log"

	"github.com/alexdunne/gs-onboarding/internal/gateway"
	"github.com/alexdunne/gs-onboarding/internal/gateway/hackernews"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Addr           string
	GRPCServerAddr string
}

func loadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read env file")
	}

	return &Config{
		Addr:           viper.GetString("GATEWAY_ADDR"),
		GRPCServerAddr: viper.GetString("GRPC_SERVER_ADDR"),
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

	client, err := hackernews.New(cfg.GRPCServerAddr)
	if err != nil {
		log.Fatal(errors.Wrap(err, "creating grpc client"))
	}
	defer client.Close()

	router := echo.New()
	router.HideBanner = true
	router.Use(
		middleware.Recover(),
		middleware.Logger(),
	)

	h := gateway.Handler{
		HNClient: client,
	}

	router.GET("/all", h.GetAllItems)
	router.GET("/stories", h.GetStories)
	router.GET("/jobs", h.GetJobs)

	if err := router.Start(cfg.Addr); err != nil {
		logger.Fatal("starting server", zap.Error(err))
	}
}
