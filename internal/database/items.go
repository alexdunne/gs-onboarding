package database

import (
	"context"
	"fmt"

	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
)

func (c *Client) GetAll(ctx context.Context) ([]models.Item, error) {
	var items []models.Item
	err := pgxscan.Select(ctx, c.pool, &items, `SELECT id, type, content, url, score, title, created_at, created_by FROM items`)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Client) GetStories(ctx context.Context) ([]models.Item, error) {
	var items []models.Item
	err := pgxscan.Select(
		ctx,
		c.pool,
		&items,
		`SELECT id, type, content, url, score, title, created_at, created_by FROM items WHERE type = 'story'`,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Client) GetJobs(ctx context.Context) ([]models.Item, error) {
	var items []models.Item
	err := pgxscan.Select(
		ctx,
		c.pool,
		&items,
		`SELECT id, type, content, url, score, title, created_at, created_by FROM items WHERE type = 'job'`,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Client) Insert(ctx context.Context, item models.Item) error {
	sql := `
	INSERT INTO items (id, type, content, url, score, title, created_by, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (id) DO NOTHING
	`

	if _, err := c.pool.Exec(
		ctx, sql, item.ID, item.Type, item.Content, item.URL,
		item.Score, item.Title, item.CreatedBy, item.CreatedAt,
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("inserting item (id: %d)", item.ID))
	}

	return nil
}
