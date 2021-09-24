package api

import (
	"context"
	"fmt"
	"time"

	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Cache is an interace to expose cache methods
type Cache interface {
	GetAll(ctx context.Context) ([]models.Item, error)
	GetStories(ctx context.Context) ([]models.Item, error)
	GetJobs(ctx context.Context) ([]models.Item, error)
}

type itemCache struct {
	db     database.Database
	cache  *cache.Cache
	ring   *redis.Ring
	ttl    time.Duration
	logger *zap.Logger
}

// CacheOption is an interface for a functional option
type CacheOption func(c *itemCache)

// WithTTL is a functional option to configure the cache TTL
func WithTTL(ttl time.Duration) CacheOption {
	return func(c *itemCache) {
		c.ttl = ttl
	}
}

// NewCache creates a new cache
func NewCache(ctx context.Context, redisAddr string, db database.Database, logger *zap.Logger, opts ...CacheOption) (*itemCache, error) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"leader": redisAddr,
		},
	})

	c := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	_, err := ring.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "pinging with new client")
	}

	ret := &itemCache{
		db:     db,
		cache:  c,
		ring:   ring,
		ttl:    5 * time.Minute,
		logger: logger,
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret, nil
}

// GetAll fetches all items from the cache and falls back to fetching from the database
func (c *itemCache) GetAll(ctx context.Context) ([]models.Item, error) {
	var items []models.Item

	key := "items:all"
	err := c.cache.Once(&cache.Item{
		Key:   key,
		Value: &items,
		TTL:   c.ttl,
		Do: func(*cache.Item) (interface{}, error) {
			c.logger.Info(fmt.Sprintf("%s cache missed. fetching from source", key))
			return c.db.GetAll(ctx)
		},
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetStories fetches all story items from the cache and falls back to fetching from the database
func (c *itemCache) GetStories(ctx context.Context) ([]models.Item, error) {
	var items []models.Item

	key := "items:stories"
	err := c.cache.Once(&cache.Item{
		Key:   key,
		Value: &items,
		TTL:   c.ttl,
		Do: func(*cache.Item) (interface{}, error) {
			c.logger.Info(fmt.Sprintf("%s cache missed. fetching from source", key))
			return c.db.GetStories(ctx)
		},
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetJobs fetches all job items from the cache and falls back to fetching from the database
func (c *itemCache) GetJobs(ctx context.Context) ([]models.Item, error) {
	var items []models.Item

	key := "items:jobs"
	err := c.cache.Once(&cache.Item{
		Key:   key,
		Value: &items,
		TTL:   c.ttl,
		Do: func(*cache.Item) (interface{}, error) {
			c.logger.Info(fmt.Sprintf("%s cache missed. fetching from source", key))
			return c.db.GetJobs(ctx)
		},
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *itemCache) Close() {
	c.ring.Close()
}
