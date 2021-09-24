package gateway

import (
	"net/http"

	"github.com/alexdunne/gs-onboarding/internal/gateway/hackernews"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	HNClient hackernews.Client
}

// GetAllItems handles requests to GET /all
func (h *Handler) GetAllItems(c echo.Context) error {
	items, err := h.HNClient.FetchAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

// GetAllItems handles requests to GET /stories
func (h *Handler) GetStories(c echo.Context) error {
	items, err := h.HNClient.FetchStories(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

// GetAllItems handles requests to GET /jobs
func (h *Handler) GetJobs(c echo.Context) error {
	items, err := h.HNClient.FetchJobs(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}
