package database

import (
	"context"

	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) GetAll(ctx context.Context) ([]models.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Item), args.Error(1)
}

func (m *Mock) GetStories(ctx context.Context) ([]models.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Item), args.Error(1)
}

func (m *Mock) GetJobs(ctx context.Context) ([]models.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Item), args.Error(1)
}

func (m *Mock) Write(ctx context.Context, item models.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)

}
