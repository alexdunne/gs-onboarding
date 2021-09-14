package memory

import "github.com/alexdunne/gs-onboarding/hn"

type ItemStore struct {
	store map[string]hn.Item
}

func NewItemStore() *ItemStore {
	return &ItemStore{
		store: map[string]hn.Item{},
	}
}

func (s *ItemStore) GetAll() hn.Items {
	items := s.find(func(item hn.Item) bool {
		return true
	})

	return items
}

func (s *ItemStore) GetStories() hn.Items {
	items := s.find(func(item hn.Item) bool {
		return item.Type == hn.StoryType
	})

	return items
}

func (s *ItemStore) GetJobs() hn.Items {
	items := s.find(func(item hn.Item) bool {
		return item.Type == hn.JobType
	})

	return items
}

type findItemFilter func(item hn.Item) bool

func (s *ItemStore) find(fn findItemFilter) hn.Items {
	items := hn.Items{}

	for _, v := range s.store {
		if fn(v) {
			items = append(items, v)
		}
	}

	return items
}
