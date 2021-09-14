package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexdunne/gs-onboarding/hn"
	"github.com/alexdunne/gs-onboarding/internal/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Addr        string
	DatabaseDSN string
}

func loadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read env file")
	}

	return &Config{
		Addr: viper.GetString("ADDR"),
		DatabaseDSN: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			viper.GetString("DATABASE_USER"),
			viper.GetString("DATABASE_PASSWORD"),
			viper.GetString("DATABASE_HOST"),
			viper.GetString("DATABASE_PORT"),
			viper.GetString("DATABASE_DB"),
		),
	}, nil
}

type ItemStore interface {
	GetAll(ctx context.Context) (hn.Items, error)
	GetStories(ctx context.Context) (hn.Items, error)
	GetJobs(ctx context.Context) (hn.Items, error)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	ctx := context.Background()

	store := postgres.NewItemStore()
	if err := store.Open(ctx, cfg.DatabaseDSN); err != nil {
		return err
	}
	defer store.Close(ctx)

	server := NewServer(store)
	return server.start(cfg.Addr)
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

func (s *server) start(addr string) error {
	return s.router.Start(addr)
}

func (s *server) handleGetAllItems(c echo.Context) error {
	items, err := s.store.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

func (s *server) handleGetStories(c echo.Context) error {
	items, err := s.store.GetStories(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

func (s *server) handleGetJobs(c echo.Context) error {
	items, err := s.store.GetJobs(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}
