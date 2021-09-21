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
	db          database.Database
	hn          hn.Client
	workerCount int
}

func New(logger *zap.Logger, db database.Database, hn hn.Client, workerCount int) *Consumer {
	return &Consumer{
		logger:      logger,
		db:          db,
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
	worker := &Worker{
		logger: c.logger,
		db:     c.db,
		hn:     c.hn,
	}
	for i := 0; i < c.workerCount; i++ {
		wg.Add(1)
		go worker.run(ctx, idStream, wg)
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
