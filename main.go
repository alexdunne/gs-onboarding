package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	server := newServer()
	server.start("localhost:8000")
}

type server struct {
	router *echo.Echo
}

func newServer() *server {
	s := &server{}

	s.router = echo.New()
	s.router.Use(
		middleware.Recover(),
		middleware.Logger(),
	)

	s.router.GET("/", hello)

	return s
}

func (s *server) start(address string) {
	err := s.router.Start(address)
	s.router.Logger.Fatal(err)
}

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"hello": "world",
	})
}
