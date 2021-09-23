package hn

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

func (m *Mock) FetchTopStories() ([]int, error) {
	args := m.Called()

	idsArg, ok := args.Get(0).([]int)
	if !ok {
		return nil, nil
	}

	return idsArg, args.Error(1)
}

func (m *Mock) FetchItem(id int) (*Item, error) {
	args := m.Called(id)

	itemArg, ok := args.Get(0).(*Item)
	if !ok {
		return nil, nil
	}

	return itemArg, args.Error(1)
}
