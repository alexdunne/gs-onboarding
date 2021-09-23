package database

import (
	"context"

	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

// Database is a interface to expose methods to fetch and store items
type Database interface {
	GetAll(ctx context.Context) ([]models.Item, error)
	GetStories(ctx context.Context) ([]models.Item, error)
	GetJobs(ctx context.Context) ([]models.Item, error)
	Write(ctx context.Context, item models.Item) error
}

// Client for database
type Client struct {
	pool *pgxpool.Pool
}

// New starts db connection
func New(ctx context.Context, connStr string) (*Client, error) {
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed connecting to the database")
	}

	return &Client{pool: pool}, nil
}

// Close closes db connection
func (c *Client) Close() {
	c.pool.Close()
}
