package gateway

import (
	"net/http"

	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB database.Database
}

func (h *Handler) HandleGetAllItems(c echo.Context) error {
	items, err := h.DB.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

func (h *Handler) HandleGetStories(c echo.Context) error {
	items, err := h.DB.GetStories(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

func (h *Handler) HandleGetJobs(c echo.Context) error {
	items, err := h.DB.GetJobs(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": items,
	})
}
