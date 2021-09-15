package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexdunne/gs-onboarding/hn"
)

type StubItemStore struct {
	items            map[int]*hn.Item
	getAllCalled     int
	getStoriesCalled int
	getJobsCalled    int
}

func (s *StubItemStore) GetAll(ctx context.Context) (hn.Items, error) {
	s.getAllCalled++

	items := hn.Items{}
	for _, v := range s.items {
		items = append(items, v)
	}

	return items, nil
}

func (s *StubItemStore) GetStories(ctx context.Context) (hn.Items, error) {
	s.getStoriesCalled++

	items := hn.Items{}
	for _, v := range s.items {
		if v.Type == hn.StoryType {
			items = append(items, v)
		}
	}

	return items, nil
}

func (s *StubItemStore) GetJobs(ctx context.Context) (hn.Items, error) {
	s.getJobsCalled++

	items := hn.Items{}
	for _, v := range s.items {
		if v.Type == hn.JobType {
			items = append(items, v)
		}
	}

	return items, nil
}

func TestGetAllItems(t *testing.T) {
	store := &StubItemStore{
		items: map[int]*hn.Item{
			1: {
				ID:        1,
				Type:      "story",
				Text:      "Hello, world!",
				URL:       "gymshark.com",
				Score:     128,
				Title:     "Intro",
				CreatedAt: time.Now(),
				CreatedBy: "Some rando",
			},
			2: {
				ID:        2,
				Type:      "story",
				Text:      "Hello Reloaded",
				URL:       "gymshark.com",
				Score:     256,
				Title:     "I'll be back",
				CreatedAt: time.Now(),
				CreatedBy: "Some rando",
			},
			3: {
				ID:        3,
				Type:      "job",
				Text:      "Software Engineer role",
				URL:       "gymshark.com/careers",
				Score:     512,
				Title:     "Software Engineer",
				CreatedAt: time.Now(),
				CreatedBy: "Shark Boi",
			},
		},
	}
	server := NewServer(store)

	req := httptest.NewRequest(http.MethodGet, "/all", nil)
	rec := httptest.NewRecorder()
	c := server.router.NewContext(req, rec)

	server.handleGetAllItems(c)

	if store.getAllCalled != 1 {
		t.Errorf("expected GetAll to be called once, got %v", store.getAllCalled)
	}

	assertStatusCode(t, rec.Code, http.StatusOK)

	res := decodeResponse(t, rec.Body)
	if len(res.Items) != 3 {
		t.Errorf("received the wrong number of items. got %v, want %v", len(res.Items), 3)
	}
}

func TestGetStories(t *testing.T) {
	store := &StubItemStore{
		items: map[int]*hn.Item{
			1: {
				ID:        1,
				Type:      "story",
				Text:      "Hello, world!",
				URL:       "gymshark.com",
				Score:     128,
				Title:     "Intro",
				CreatedAt: time.Now(),
				CreatedBy: "Some rando",
			},
			3: {
				ID:        3,
				Type:      "job",
				Text:      "Software Engineer role",
				URL:       "gymshark.com/careers",
				Score:     512,
				Title:     "Software Engineer",
				CreatedAt: time.Now(),
				CreatedBy: "Shark Boi",
			},
		},
	}
	server := NewServer(store)

	req := httptest.NewRequest(http.MethodGet, "/stories", nil)
	rec := httptest.NewRecorder()
	c := server.router.NewContext(req, rec)

	server.handleGetStories(c)

	if store.getStoriesCalled != 1 {
		t.Errorf("expected GetStories to be called once, got %v", store.getStoriesCalled)
	}

	assertStatusCode(t, rec.Code, http.StatusOK)

	res := decodeResponse(t, rec.Body)
	if len(res.Items) != 1 {
		t.Errorf("received the wrong number of items. got %v, want %v", len(res.Items), 1)
	}

	if res.Items[0].ID != 1 {
		t.Errorf("expected first returned story to have id %v, got: %v", 1, res.Items[0].ID)
	}
}

func TestGetJobs(t *testing.T) {
	store := &StubItemStore{
		items: map[int]*hn.Item{
			1: {
				ID:        1,
				Type:      "story",
				Text:      "Hello, world!",
				URL:       "gymshark.com",
				Score:     128,
				Title:     "Intro",
				CreatedAt: time.Now(),
				CreatedBy: "Some rando",
			},
			3: {
				ID:        3,
				Type:      "job",
				Text:      "Software Engineer role",
				URL:       "gymshark.com/careers",
				Score:     512,
				Title:     "Software Engineer",
				CreatedAt: time.Now(),
				CreatedBy: "Shark Boi",
			},
		},
	}
	server := NewServer(store)

	req := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	rec := httptest.NewRecorder()
	c := server.router.NewContext(req, rec)

	server.handleGetJobs(c)

	if store.getJobsCalled != 1 {
		t.Errorf("expected GetJobs to be called once, got %v", store.getJobsCalled)
	}

	assertStatusCode(t, rec.Code, http.StatusOK)

	res := decodeResponse(t, rec.Body)
	if len(res.Items) != 1 {
		t.Errorf("received the wrong number of items. got %v, want %v", len(res.Items), 1)
	}

	if res.Items[0].ID != 3 {
		t.Errorf("expected first returned story to have id %v, got: %v", 3, res.Items[0].ID)
	}
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("received the wrong status code. got %v, want %v", got, want)
	}
}

type successResponse struct {
	Items []hn.Item `json:"items"`
}

func decodeResponse(t testing.TB, r io.Reader) successResponse {
	t.Helper()

	var res successResponse
	err := json.NewDecoder(r).Decode(&res)
	if err != nil {
		t.Fatalf("unable to decode response body. body: %v, error: %v", r, err)
	}

	return res
}
