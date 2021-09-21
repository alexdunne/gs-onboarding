package consumer

import (
	"context"
	"sync"
	"testing"

	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/alexdunne/gs-onboarding/pkg/hn"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestRun(t *testing.T) {
	type testcase struct {
		name        string
		database    *database.Mock
		hn          *hn.Mock
		ids         []int
		expectMocks func(t *testing.T, dbMock *database.Mock, hnMock *hn.Mock)
	}

	tests := []testcase{
		{
			name:     "one item",
			database: &database.Mock{},
			hn:       &hn.Mock{},
			ids:      []int{1},
			expectMocks: func(t *testing.T, dbMock *database.Mock, hnMock *hn.Mock) {
				hnMock.On("FetchItem", 1).Return(&hn.Item{ID: 1}, nil)
				dbMock.On("Write", context.TODO(), mock.AnythingOfType("models.Item")).Return(nil)
			},
		},
		{
			name:     "three items",
			database: &database.Mock{},
			hn:       &hn.Mock{},
			ids:      []int{1, 2, 3},
			expectMocks: func(t *testing.T, dbMock *database.Mock, hnMock *hn.Mock) {
				hnMock.On("FetchItem", 1).Return(&hn.Item{ID: 1}, nil)
				hnMock.On("FetchItem", 2).Return(&hn.Item{ID: 2}, nil)
				hnMock.On("FetchItem", 3).Return(&hn.Item{ID: 3}, nil)
				dbMock.On("Write", context.TODO(), mock.AnythingOfType("models.Item")).Return(nil).Times(3)
			},
		},
		{
			name:     "ignores dead or deleted items",
			database: &database.Mock{},
			hn:       &hn.Mock{},
			ids:      []int{1, 2, 3},
			expectMocks: func(t *testing.T, dbMock *database.Mock, hnMock *hn.Mock) {
				hnMock.On("FetchItem", 1).Return(&hn.Item{ID: 1, Dead: true}, nil)
				hnMock.On("FetchItem", 2).Return(&hn.Item{ID: 2, Deleted: true}, nil)
				hnMock.On("FetchItem", 3).Return(&hn.Item{ID: 3, Dead: true, Deleted: true}, nil)
				// dbMock.On("Write", context.TODO(), mock.AnythingOfType("models.Item")).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMocks != nil {
				tt.expectMocks(t, tt.database, tt.hn)
			}

			logger, err := zap.NewProduction()
			if err != nil {
				panic(err)
			}

			idStream := make(chan int)
			if len(tt.ids) > 0 {
				go func() {
					for _, id := range tt.ids {
						idStream <- id
					}
					close(idStream)
				}()
			}

			worker := NewWorker(logger, tt.database, tt.hn)
			wg := &sync.WaitGroup{}
			wg.Add(1)

			go worker.run(context.TODO(), idStream, wg)
			wg.Wait()

			if tt.expectMocks != nil {
				tt.database.AssertExpectations(t)
				tt.hn.AssertExpectations(t)
			}
		})
	}
}
