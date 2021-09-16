package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

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
