package inmem

import "github.com/alexdunne/gs-onboarding/hn"

type InMemoryItemStore struct {
	store map[string]hn.Item
}

func NewInMemoryItemStore() *InMemoryItemStore {
	return &InMemoryItemStore{
		store: map[string]hn.Item{},
	}
}

func (s *InMemoryItemStore) GetAll() hn.Items {
	items := s.find(func(item hn.Item) bool {
		return true
	})

	return items
}

func (s *InMemoryItemStore) GetStories() hn.Items {
	items := s.find(func(item hn.Item) bool {
		return item.Type == hn.StoryType
	})

	return items
}

func (s *InMemoryItemStore) GetJobs() hn.Items {
	items := s.find(func(item hn.Item) bool {
		return item.Type == hn.JobType
	})

	return items
}

type findItemFilter func(item hn.Item) bool

func (s *InMemoryItemStore) find(fn findItemFilter) hn.Items {
	items := hn.Items{}

	for _, v := range s.store {
		if fn(v) {
			items = append(items, v)
		}
	}

	return items
}
