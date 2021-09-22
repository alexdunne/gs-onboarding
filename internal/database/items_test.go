package database

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var testDB *TestDatabase

func TestMain(m *testing.M) {
	db, err := createTestDB()
	if err != nil {
		panic(errors.Wrap(err, "creating test db"))
	}

	testDB = db

	code := m.Run()

	if err := db.cleanUp(); err != nil {
		panic(errors.Wrap(err, "failed to clean up"))
	}

	os.Exit(code)
}

func TestGetAll(t *testing.T) {
	client := &Client{
		pool: testDB.pool,
	}

	type testcase struct {
		name              string
		seed              func(ctx context.Context)
		expectedItemCount int
	}

	tests := []testcase{
		{
			name: "no items",
			seed: func(ctx context.Context) {
				// no-op
			},
			expectedItemCount: 0,
		},
		{
			name: "one story",
			seed: func(ctx context.Context) {
				client.Write(ctx, models.Item{
					ID:        1,
					Type:      "story",
					Content:   "Hello, world",
					URL:       "gymshark.com",
					Score:     10,
					Title:     "Intro",
					CreatedAt: time.Now(),
					CreatedBy: "shark boi",
				})
			},
			expectedItemCount: 1,
		},
		{
			name: "one story and one job",
			seed: func(ctx context.Context) {
				client.Write(ctx, models.Item{
					ID:        1,
					Type:      "story",
					Content:   "Hello, world",
					URL:       "gymshark.com",
					Score:     10,
					Title:     "Intro",
					CreatedAt: time.Now(),
					CreatedBy: "shark boi",
				})

				client.Write(ctx, models.Item{
					ID:        2,
					Type:      "job",
					Content:   "Work for us",
					URL:       "gymshark.com/careers",
					Score:     10,
					Title:     "Senior Software Engineer",
					CreatedAt: time.Now(),
					CreatedBy: "lava gurl",
				})
			},
			expectedItemCount: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := testDB.reset()
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.TODO()
			tc.seed(ctx)
			items, err := client.GetAll(ctx)

			assert.Equal(t, tc.expectedItemCount, len(items))
			assert.NoError(t, err)
		})
	}
}

func TestGetStories(t *testing.T) {
	client := &Client{
		pool: testDB.pool,
	}

	type testcase struct {
		name              string
		seed              func(ctx context.Context)
		expectedItemCount int
	}

	tests := []testcase{
		{
			name: "no items",
			seed: func(ctx context.Context) {
				// no-op
			},
			expectedItemCount: 0,
		},
		{
			name: "one story",
			seed: func(ctx context.Context) {
				client.Write(ctx, models.Item{
					ID:        1,
					Type:      "story",
					Content:   "Hello, world",
					URL:       "gymshark.com",
					Score:     10,
					Title:     "Intro",
					CreatedAt: time.Now(),
					CreatedBy: "shark boi",
				})
			},
			expectedItemCount: 1,
		},
		{
			name: "one job",
			seed: func(ctx context.Context) {
				client.Write(ctx, models.Item{
					ID:        2,
					Type:      "job",
					Content:   "Work for us",
					URL:       "gymshark.com/careers",
					Score:     10,
					Title:     "Senior Software Engineer",
					CreatedAt: time.Now(),
					CreatedBy: "lava gurl",
				})
			},
			expectedItemCount: 0,
		},
		{
			name: "one story and one job",
			seed: func(ctx context.Context) {
				client.Write(ctx, models.Item{
					ID:        1,
					Type:      "story",
					Content:   "Hello, world",
					URL:       "gymshark.com",
					Score:     10,
					Title:     "Intro",
					CreatedAt: time.Now(),
					CreatedBy: "shark boi",
				})

				client.Write(ctx, models.Item{
					ID:        2,
					Type:      "job",
					Content:   "Work for us",
					URL:       "gymshark.com/careers",
					Score:     10,
					Title:     "Senior Software Engineer",
					CreatedAt: time.Now(),
					CreatedBy: "lava gurl",
				})
			},
			expectedItemCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := testDB.reset()
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.TODO()
			tc.seed(ctx)

			items, err := client.GetStories(ctx)

			assert.Equal(t, tc.expectedItemCount, len(items))
			assert.NoError(t, err)
		})
	}
}

func TestGetJobs(t *testing.T) {
	client := &Client{
		pool: testDB.pool,
	}

	type testcase struct {
		name              string
		seed              func(ctx context.Context)
		expectedItemCount int
	}

	tests := []testcase{
		{
			name: "no items",
			seed: func(ctx context.Context) {
				// no-op
			},
			expectedItemCount: 0,
		},
		{
			name: "one story",
			seed: func(ctx context.Context) {
				client.Write(ctx, models.Item{
					ID:        1,
					Type:      "story",
					Content:   "Hello, world",
					URL:       "gymshark.com",
					Score:     10,
					Title:     "Intro",
					CreatedAt: time.Now(),
					CreatedBy: "shark boi",
				})
			},
			expectedItemCount: 0,
		},
		{
			name: "one job",
			seed: func(ctx context.Context) {

				client.Write(ctx, models.Item{
					ID:        2,
					Type:      "job",
					Content:   "Work for us",
					URL:       "gymshark.com/careers",
					Score:     10,
					Title:     "Senior Software Engineer",
					CreatedAt: time.Now(),
					CreatedBy: "lava gurl",
				})
			},
			expectedItemCount: 1,
		},
		{
			name: "one story and one job",
			seed: func(ctx context.Context) {
				client.Write(ctx, models.Item{
					ID:        1,
					Type:      "story",
					Content:   "Hello, world",
					URL:       "gymshark.com",
					Score:     10,
					Title:     "Intro",
					CreatedAt: time.Now(),
					CreatedBy: "shark boi",
				})

				client.Write(ctx, models.Item{
					ID:        2,
					Type:      "job",
					Content:   "Work for us",
					URL:       "gymshark.com/careers",
					Score:     10,
					Title:     "Senior Software Engineer",
					CreatedAt: time.Now(),
					CreatedBy: "lava gurl",
				})
			},
			expectedItemCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := testDB.reset()
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.TODO()
			tc.seed(ctx)

			items, err := client.GetJobs(ctx)

			assert.Equal(t, tc.expectedItemCount, len(items))
			assert.NoError(t, err)
		})
	}
}
