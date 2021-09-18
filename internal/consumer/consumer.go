package consumer

import (
	"context"
	"sync"

	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/alexdunne/gs-onboarding/pkg/hn"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Consumer struct {
	logger      *zap.Logger
	writer      database.ItemWriter
	hn          hn.Client
	workerCount int
}

func New(logger *zap.Logger, writer database.ItemWriter, hn hn.Client, workerCount int) *Consumer {
	return &Consumer{
		logger:      logger,
		writer:      writer,
		hn:          hn,
		workerCount: workerCount,
	}
}

// run fetches the top story ids from Hacker News and passes them along a pipeline to be fetched and inserted
func (c *Consumer) Run(ctx context.Context) error {
	c.logger.Info("running consumer")

	idStream := make(chan int)

	// create workers to fetch and insert the data
	wg := &sync.WaitGroup{}
	for i := 0; i < c.workerCount; i++ {
		wg.Add(1)
		go worker(ctx, c.logger, c.writer, c.hn, idStream, wg)
	}

	ids, err := c.hn.FetchTopStories()
	if err != nil {
		return errors.Wrap(err, "fetching top stories")
	}
	c.logger.Info("fetched top story ids", zap.Int("count", len(ids)))

	for _, id := range ids {
		idStream <- id
	}
	close(idStream)

	wg.Wait()

	c.logger.Info("finished inserting items")

	return nil
}
