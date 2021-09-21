package hn

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

func (m *Mock) FetchTopStories() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}

func (m *Mock) FetchItem(id int) (*Item, error) {
	args := m.Called(id)
	return args.Get(0).(*Item), args.Error(1)
}
