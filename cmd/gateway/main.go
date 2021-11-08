package main

import (
	"fmt"
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
	viper.AutomaticEnv()

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
		logger.Fatal("creating grpc client", zap.Error(err))
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

	logger.Info(fmt.Sprintf("starting server at %s", cfg.Addr))
	if err := router.Start(cfg.Addr); err != nil {
		logger.Fatal("starting server", zap.Error(err))
	}
}
