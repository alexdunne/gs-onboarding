package database

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {

	type testcase struct {
		name              string
		seed              func(ctx context.Context, c *Client)
		expectedItemCount int
	}

	tests := []testcase{
		{
			name: "no items",
			seed: func(ctx context.Context, c *Client) {
				// no-op
			},
			expectedItemCount: 0,
		},
		// {
		// 	name: "one story",
		// 	seed: func(ctx context.Context) {
		// 		client.Write(ctx, models.Item{
		// 			ID:        123,
		// 			Type:      "story",
		// 			Content:   "Hello, world",
		// 			URL:       "gymshark.com",
		// 			Score:     10,
		// 			Title:     "Intro",
		// 			CreatedAt: time.Now(),
		// 			CreatedBy: "shark boi",
		// 		})
		// 	},
		// 	expectedItemCount: 1,
		// },
		// {
		// 	name: "one story and one job",
		// 	seed: func(ctx context.Context) {
		// 		client.Write(ctx, models.Item{
		// 			ID:        123,
		// 			Type:      "story",
		// 			Content:   "Hello, world",
		// 			URL:       "gymshark.com",
		// 			Score:     10,
		// 			Title:     "Intro",
		// 			CreatedAt: time.Now(),
		// 			CreatedBy: "shark boi",
		// 		})

		// 		client.Write(ctx, models.Item{
		// 			ID:        123,
		// 			Type:      "job",
		// 			Content:   "Work for us",
		// 			URL:       "gymshark.com/careers",
		// 			Score:     10,
		// 			Title:     "Senior Software Engineer",
		// 			CreatedAt: time.Now(),
		// 			CreatedBy: "lava gurl",
		// 		})
		// 	},
		// 	expectedItemCount: 2,
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pool, cleanUp, err := createTestDB()
			if err != nil {
				panic(errors.Wrap(err, "creating test db"))
			}
			defer func() {
				if err := cleanUp(); err != nil {
					panic(errors.Wrap(err, "failed to clean up"))
				}
			}()

			testClient := &Client{
				pool: pool,
			}

			ctx := context.TODO()

			tc.seed(ctx, testClient)

			items, err := testClient.GetAll(ctx)

			assert.Equal(t, tc.expectedItemCount, len(items))
			assert.NoError(t, err)
		})
	}
}
