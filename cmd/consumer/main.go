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

	"github.com/alexdunne/gs-onboarding/hn"
	"github.com/alexdunne/gs-onboarding/internal/postgres"
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
	if err := run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := loadConfig()
	if err != nil {
		return errors.Wrap(err, "loading config")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return errors.Wrap(err, "creating logger")
	}
	defer logger.Sync()

	go func() {
		// run the consumer the first time when the goroutine starts
		consumer(ctx, cfg, logger)

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.Tick(cfg.WorkerIntervalDuration):
				logger.Info("running consumer")
				consumer(ctx, cfg, logger)
			}
		}
	}()

	// handle interrupts and propagate the changes across the consumer pipeline
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()

	return nil
}

// consumer fetches the top story ids from Hacker News and passes them along a pipeline to be fetched and inserted
func consumer(ctx context.Context, cfg *Config, logger *zap.Logger) error {
	logger.Info("opening database connection")
	store := postgres.NewItemStore()
	if err := store.Open(ctx, cfg.DatabaseDSN); err != nil {
		return errors.Wrap(err, "opening store db connection")
	}
	defer store.Close(ctx)

	logger.Info("creating HN client")
	client := hn.NewClient()

	ids, err := client.FetchTopStories()
	if err != nil {
		return errors.Wrap(err, "fetching top stories")
	}
	logger.Info("fetched top story ids", zap.Int("count", len(ids)))

	// convert the top story ids into a channel of ids
	idStream := generator(ctx, ids)

	// create workers to fetch items from the HN API
	fetchers := make([]<-chan fetchItemResult, cfg.WorkerCount)
	for i := 0; i < cfg.WorkerCount; i++ {
		fetchers[i] = fetchItem(ctx, idStream, client)
	}

	// fan-in the fetched items and decide what to do with them
	for result := range fanInFetchedItems(ctx, fetchers...) {
		if result.Error != nil {
			logger.Error("error fetching item", zap.Error(result.Error))
			continue
		}

		if result.Item.Dead || result.Item.Deleted {
			// ignore dead or deleted items
			break
		}

		logger.Info("inserting item", zap.Int("id", result.Item.ID))
		store.Insert(ctx, result.Item)
	}

	logger.Info("finished inserting items")

	return nil
}

func generator(ctx context.Context, ids []int) <-chan int {
	idStream := make(chan int)

	go func() {
		defer close(idStream)

		for _, id := range ids {
			select {
			case <-ctx.Done():
				return
			case idStream <- id:
			}
		}
	}()

	return idStream
}

type fetchItemResult struct {
	Item  *hn.Item
	Error error
}

func fetchItem(ctx context.Context, idStream <-chan int, client *hn.Client) <-chan fetchItemResult {
	itemStream := make(chan fetchItemResult)

	go func() {
		defer close(itemStream)

		for {
			select {
			case <-ctx.Done():
				return
			case id, ok := <-idStream:
				if !ok {
					return
				}

				item, err := client.FetchItem(id)
				// wrap both the return item and the error so that whatever consumes this channel can deal with the errors more appropriately
				result := fetchItemResult{Item: item, Error: err}

				select {
				case <-ctx.Done():
					return
				case itemStream <- result:
				}

			}

		}

	}()

	return itemStream
}

func fanInFetchedItems(ctx context.Context, channels ...<-chan fetchItemResult) <-chan fetchItemResult {
	aggregateStream := make(chan fetchItemResult)

	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, c := range channels {
		go func(c <-chan fetchItemResult) {
			defer wg.Done()

			for item := range c {
				select {
				case <-ctx.Done():
					return
				case aggregateStream <- item:
				}
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(aggregateStream)
	}()

	return aggregateStream
}
