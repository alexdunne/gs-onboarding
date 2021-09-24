package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/alexdunne/gs-onboarding/internal/consumer"
	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/alexdunne/gs-onboarding/internal/queue"
	"github.com/alexdunne/gs-onboarding/pkg/hn"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	queueName = "items"
)

type Config struct {
	WorkerCount            int
	WorkerIntervalDuration time.Duration
	DatabaseDSN            string
	RabbitMQURL            string
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
		RabbitMQURL: fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			viper.GetString("RABBITMQ_USER"),
			viper.GetString("RABBITMQ_PASSWORD"),
			viper.GetString("RABBITMQ_HOST"),
			viper.GetString("RABBITMQ_PORT"),
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
		logger.Fatal("failed to create db connection", zap.Error(err))
	}
	defer db.Close()

	logger.Info("creating HN client")
	hackerNewsClient := hn.New()

	queueClient, err := queue.New(cfg.RabbitMQURL, queueName, logger)
	if err != nil {
		logger.Fatal("failed to create RabbitMQ connection", zap.Error(err))
	}
	defer queueClient.Close()

	messages, err := queueClient.Consume(ctx)
	if err != nil {
		logger.Fatal("failed to consumer message from RabbitMQ", zap.Error(err))
	}

	w := consumer.NewWorker(logger, db, hackerNewsClient)
	wg := &sync.WaitGroup{}

	for i := 0; i < cfg.WorkerCount; i++ {
		wg.Add(1)
		go w.Run(ctx, messages, wg)
	}

	if err := seed(ctx, hackerNewsClient, queueClient, cfg.WorkerIntervalDuration, logger); err != nil {
		logger.Error("failed to seed ids", zap.Error(err))
	}

	wg.Wait()
}

func seed(ctx context.Context, hackerNewsClient hn.Client, queueClient queue.Queue, interval time.Duration, logger *zap.Logger) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(interval):
			ids, err := hackerNewsClient.FetchTopStories()
			if err != nil {
				return errors.Wrap(err, "fetching top stores")
			}

			logger.Info("fetched top story ids", zap.Int("count", len(ids)))

			for _, id := range ids {
				if err := queueClient.Publish(&queue.Message{ID: id}); err != nil {
					return err
				}
			}

		}
	}
}
