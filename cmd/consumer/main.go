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
	"github.com/alexdunne/gs-onboarding/internal/database"
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
		return errors.Wrap(err, "loading config")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return errors.Wrap(err, "creating logger")
	}
	defer logger.Sync()

	logger.Info("opening database connection")
	db, err := database.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		return errors.Wrap(err, "opening store db connection")
	}
	defer db.Close()

	logger.Info("creating HN client")
	hn := hn.New()

	consumer := &Consumer{
		logger:      logger,
		db:          db,
		hn:          hn,
		workerCount: cfg.WorkerCount,
	}

	for {
		logger.Info("running consumer")
		err := consumer.run(ctx)
		if err != nil {
			logger.Error("consumer failed", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(cfg.WorkerIntervalDuration):
		}
	}
}

type DBWriter interface {
	Insert(ctx context.Context, item database.Item) error
}

type Consumer struct {
	logger      *zap.Logger
	db          DBWriter
	hn          *hn.Client
	workerCount int
}

// run fetches the top story ids from Hacker News and passes them along a pipeline to be fetched and inserted
func (w *Consumer) run(ctx context.Context) error {
	ids, err := w.hn.FetchTopStories()
	if err != nil {
		return errors.Wrap(err, "fetching top stories")
	}
	w.logger.Info("fetched top story ids", zap.Int("count", len(ids)))

	// convert the top story ids into a channel of ids
	idStream := generator(ctx, ids)

	// create workers to fetch items from the HN API
	fetchers := make([]<-chan fetchItemResult, w.workerCount)
	for i := 0; i < w.workerCount; i++ {
		fetchers[i] = fetchItem(ctx, idStream, w.hn)
	}

	// fan-in the fetched items and decide what to do with them
	for result := range fanInFetchedItems(ctx, fetchers...) {
		if result.Error != nil {
			w.logger.Error("error fetching item", zap.Error(result.Error))
			continue
		}

		if result.Item.Dead || result.Item.Deleted {
			// ignore dead or deleted items
			continue
		}

		w.logger.Info("inserting item", zap.Int("id", result.Item.ID))
		w.db.Insert(ctx, database.Item{
			ID:        result.Item.ID,
			Type:      string(result.Item.Type),
			Content:   result.Item.Text,
			URL:       result.Item.URL,
			Score:     result.Item.Score,
			Title:     result.Item.Title,
			CreatedAt: result.Item.CreatedAt,
			CreatedBy: result.Item.CreatedBy,
		})
	}

	w.logger.Info("finished inserting items")

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
