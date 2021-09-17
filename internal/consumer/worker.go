package consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/alexdunne/gs-onboarding/pkg/hn"
	"go.uber.org/zap"
)

func worker(ctx context.Context, logger *zap.Logger, db database.Database, hn *hn.Client, idStream <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case id, ok := <-idStream:
			if !ok {
				return
			}

			item, err := hn.FetchItem(id)
			if err != nil {
				logger.Error(fmt.Sprintf("fetching item id %d", id), zap.Error(err))
				continue
			}

			if item.Dead || item.Deleted {
				// ignore dead or deleted items
				continue
			}

			logger.Info("inserting item", zap.Int("id", item.ID))
			db.Insert(ctx, database.Item{
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
