package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/alexdunne/gs-onboarding/pkg/hn"
	"go.uber.org/zap"
)

type Worker struct {
	logger *zap.Logger
	db     database.Database
	hn     hn.Client
}

func NewWorker(logger *zap.Logger, db database.Database, hn hn.Client) *Worker {
	return &Worker{
		logger: logger,
		db:     db,
		hn:     hn,
	}
}

func (w *Worker) run(ctx context.Context, idStream <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case id, ok := <-idStream:
			if !ok {
				return
			}

			item, err := w.hn.FetchItem(id)
			if err != nil {
				w.logger.Error(fmt.Sprintf("fetching item id %d", id), zap.Error(err))
				continue
			}

			log.Println(item.ID, item.Dead, item.Deleted)

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
