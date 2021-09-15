package postgres

import (
	"context"
	"fmt"

	"github.com/alexdunne/gs-onboarding/hn"
	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
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
	var items []Item
	err := pgxscan.Select(ctx, i.conn, &items, `SELECT id, type, content, url, score, title, created_at, created_by FROM items`)

	return dbItemListToModelList(items), err
}

func (i *ItemStore) GetStories(ctx context.Context) (hn.Items, error) {
	var items []Item
	err := pgxscan.Select(
		ctx,
		i.conn,
		&items,
		`SELECT id, type, content, url, score, title, created_at, created_by FROM items WHERE itemType = "story"`,
	)

	return dbItemListToModelList(items), err
}

func (i *ItemStore) GetJobs(ctx context.Context) (hn.Items, error) {
	var items []Item
	err := pgxscan.Select(
		ctx,
		i.conn,
		&items,
		`SELECT id, type, content, url, score, title, created_at, created_by FROM items WHERE itemType = "job"`,
	)

	return dbItemListToModelList(items), err
}

func (i *ItemStore) Insert(ctx context.Context, item *hn.Item) error {
	sql := `
	INSERT INTO items (id, type, content, url, score, title, created_by, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (id) DO NOTHING
	`

	if _, err := i.conn.Exec(
		ctx, sql, item.ID, item.Type, item.Text, item.URL,
		item.Score, item.Title, item.CreatedBy, item.CreatedAt,
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("inserting item (id: %d)", item.ID))
	}

	return nil
}

func dbItemListToModelList(items []Item) []*hn.Item {
	ret := make([]*hn.Item, len(items))
	for _, i := range items {
		ret = append(ret, dbItemToModel(i))
	}

	return ret
}

func dbItemToModel(item Item) *hn.Item {
	return &hn.Item{
		ID:        item.ID,
		Type:      hn.ItemType(item.Type),
		Text:      item.Content,
		URL:       item.URL,
		Score:     item.Score,
		Title:     item.Title,
		CreatedAt: item.CreatedAt,
		CreatedBy: item.CreatedBy,
	}
}
