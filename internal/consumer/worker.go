package consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/alexdunne/gs-onboarding/internal/queue"
	"github.com/alexdunne/gs-onboarding/pkg/hn"
	"go.uber.org/zap"
)

// Worker is responsible for fetching items and inserting the items in the database
type Worker struct {
	logger *zap.Logger
	db     database.Database
	hn     hn.Client
}

// NewWorker creates a new worker
func NewWorker(logger *zap.Logger, db database.Database, hn hn.Client) *Worker {
	return &Worker{
		logger: logger,
		db:     db,
		hn:     hn,
	}
}

// Run is responsible for processing messages
func (w *Worker) Run(ctx context.Context, message <-chan *queue.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-message:
			if !ok {
				return
			}

			w.logger.Info("processing message", zap.Int("id", msg.ID))

			item, err := w.hn.FetchItem(msg.ID)
			if err != nil {
				w.logger.Error(fmt.Sprintf("fetching item id %d", msg.ID), zap.Error(err))
				continue
			}

			if item.Dead || item.Deleted {
				// ignore dead or deleted items
				continue
			}

			w.logger.Info("inserting item", zap.Int("id", item.ID))
			w.db.Write(ctx, models.Item{
				ID:        item.ID,
				Type:      string(item.Type),
				Content:   item.Text,
				URL:       item.URL,
				Score:     item.Score,
				Title:     item.Title,
				CreatedAt: item.CreatedAt,
				CreatedBy: item.CreatedBy,
			})
		}
	}
}
