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
	return args.Get(0).([]models.Item), args.Error(1)
}

func (m *Mock) FetchStories(ctx context.Context) ([]models.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Item), args.Error(1)
}

func (m *Mock) FetchJobs(ctx context.Context) ([]models.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Item), args.Error(1)
}
