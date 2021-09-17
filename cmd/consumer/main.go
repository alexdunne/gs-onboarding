package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/alexdunne/gs-onboarding/internal/consumer"
	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/alexdunne/gs-onboarding/pkg/hn"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	WorkerCount            int
	WorkerIntervalDuration time.Duration
	DatabaseDSN            string
}

func loadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read env file")
	}

	c := &Config{
		WorkerCount:            runtime.NumCPU(),
		WorkerIntervalDuration: 300 * time.Second,
		DatabaseDSN: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			viper.GetString("DATABASE_USER"),
			viper.GetString("DATABASE_PASSWORD"),
			viper.GetString("DATABASE_HOST"),
			viper.GetString("DATABASE_PORT"),
			viper.GetString("DATABASE_DB"),
		),
	}

	intervalSeconds := viper.GetInt("WORKER_INTERVAL_SECONDS")
	if intervalSeconds != 0 {
		c.WorkerIntervalDuration = time.Duration(intervalSeconds) * time.Second
	}

	return c, nil
}

func main() {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	go func() {
		// handle interrupts and propagate the changes across the consumer pipeline
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancelFn()
	}()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "loading config"))
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(errors.Wrap(err, "creating logger"))
	}
	defer logger.Sync()

	db, err := database.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(errors.Wrap(err, "opening store db connection"))
	}
	defer db.Close()

	logger.Info("creating HN client")
	hn := hn.New()

	c := consumer.New(logger, db, hn, cfg.WorkerCount)

	for {
		err := c.Run(ctx)
		if err != nil {
			logger.Error("consumer failed", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return
		case <-time.Tick(cfg.WorkerIntervalDuration):
		}
	}
}
