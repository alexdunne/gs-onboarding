package postgres

import (
	"context"

	"github.com/alexdunne/gs-onboarding/hn"
	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
)

type ItemStore struct {
	conn *pgx.Conn
}

func NewItemStore() *ItemStore {
	store := &ItemStore{}

	return store
}

func (i *ItemStore) Open(ctx context.Context, connStr string) (err error) {
	if i.conn, err = pgx.Connect(ctx, connStr); err != nil {
		return err
	}

	return nil
}

func (i *ItemStore) Close(ctx context.Context) {
	i.conn.Close(ctx)
}

func (i *ItemStore) GetAll(ctx context.Context) (hn.Items, error) {
	var items hn.Items
	err := pgxscan.Select(ctx, i.conn, &items, `SELECT id, type, content, url, score, title, created_at, created_by FROM items`)

	return items, err
}

func (i *ItemStore) GetStories(ctx context.Context) (hn.Items, error) {
	var items hn.Items
	err := pgxscan.Select(
		ctx,
		i.conn,
		&items,
		`SELECT id, type, content, url, score, title, created_at, created_by FROM items WHERE itemType = "story"`,
	)

	return items, err
}

func (i *ItemStore) GetJobs(ctx context.Context) (hn.Items, error) {
	var items hn.Items
	err := pgxscan.Select(
		ctx,
		i.conn,
		&items,
		`SELECT id, type, content, url, score, title, created_at, created_by FROM items WHERE itemType = "job"`,
	)

	return items, err
}
