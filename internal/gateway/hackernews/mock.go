package hackernews

import (
	"context"

	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) FetchAll(ctx context.Context) ([]models.Item, error) {
	args := m.Called(ctx)

	itemsArg, ok := args.Get(0).([]models.Item)
	if !ok {
		return nil, nil
	}

	return itemsArg, args.Error(1)
}

func (m *Mock) FetchStories(ctx context.Context) ([]models.Item, error) {
	args := m.Called(ctx)

	itemsArg, ok := args.Get(0).([]models.Item)
	if !ok {
		return nil, nil
	}

	return itemsArg, args.Error(1)
}

func (m *Mock) FetchJobs(ctx context.Context) ([]models.Item, error) {
	args := m.Called(ctx)

	itemsArg, ok := args.Get(0).([]models.Item)
	if !ok {
		return nil, nil
	}

	return itemsArg, args.Error(1)
}
