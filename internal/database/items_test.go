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

var testClient *Client

func TestMain(m *testing.M) {
	pool, cleanUp, err := createTestDB()
	if err != nil {
		panic(errors.Wrap(err, "creating test db"))
	}

	testClient = &Client{
		pool: pool,
	}

	code := m.Run()

	if err := cleanUp(); err != nil {
		panic(errors.Wrap(err, "failed to clean up"))
	}

	os.Exit(code)
}

func TestGetAll(t *testing.T) {
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
				testClient.Write(ctx, models.Item{
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
				testClient.Write(ctx, models.Item{
					ID:        1,
					Type:      "story",
					Content:   "Hello, world",
					URL:       "gymshark.com",
					Score:     10,
					Title:     "Intro",
					CreatedAt: time.Now(),
					CreatedBy: "shark boi",
				})

				testClient.Write(ctx, models.Item{
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
			ctx := context.TODO()

			tc.seed(ctx)

			items, err := testClient.GetAll(ctx)

			assert.Equal(t, tc.expectedItemCount, len(items))
			assert.NoError(t, err)
		})
	}
}
