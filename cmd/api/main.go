package main

import (
	"net/http"

	"github.com/alexdunne/gs-onboarding/hn"
	"github.com/alexdunne/gs-onboarding/internal/memory"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ItemStore interface {
	GetAll() hn.Items
	GetStories() hn.Items
	GetJobs() hn.Items
}

func main() {
	store := memory.NewItemStore()

	server := NewServer(store)
	server.start("localhost:8000")
}

type server struct {
	router *echo.Echo
	store  ItemStore
}

func NewServer(store ItemStore) *server {
	s := &server{
		store: store,
	}

	s.router = echo.New()
	s.router.Use(
		middleware.Recover(),
		middleware.Logger(),
	)

	s.router.GET("/all", s.handleGetAllItems)
	s.router.GET("/stories", s.handleGetStories)
	s.router.GET("/jobs", s.handleGetJobs)

	return s
}

func (s *server) start(address string) {
	err := s.router.Start(address)
	s.router.Logger.Fatal(err)
}

func (s *server) handleGetAllItems(c echo.Context) error {
	items := s.store.GetAll()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

func (s *server) handleGetStories(c echo.Context) error {
	items := s.store.GetStories()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

func (s *server) handleGetJobs(c echo.Context) error {
	items := s.store.GetJobs()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}
