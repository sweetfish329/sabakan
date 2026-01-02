// Package handlers provides HTTP handlers for the Sabakan API.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sweetfish329/sabakan/backend/internal/container"
)

// ContainerHandler handles container-related HTTP requests.
type ContainerHandler struct {
	service *container.Service
}

// NewContainerHandler creates a new container handler.
func NewContainerHandler(service *container.Service) *ContainerHandler {
	return &ContainerHandler{
		service: service,
	}
}

// List handles GET /api/containers.
func (h *ContainerHandler) List(c echo.Context) error {
	containers, err := h.service.List(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, containers)
}

// Get handles GET /api/containers/:id.
func (h *ContainerHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "container ID is required")
	}

	container, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, container)
}

// Start handles POST /api/containers/:id/start.
func (h *ContainerHandler) Start(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "container ID is required")
	}

	if err := h.service.Start(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// StopRequest represents the request body for stopping a container.
type StopRequest struct {
	Timeout uint `json:"timeout"`
}

// Stop handles POST /api/containers/:id/stop.
func (h *ContainerHandler) Stop(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "container ID is required")
	}

	// Parse optional timeout from query or body
	timeout := uint(10) // Default 10 seconds
	if t := c.QueryParam("timeout"); t != "" {
		if parsed, err := strconv.ParseUint(t, 10, 32); err == nil {
			timeout = uint(parsed)
		}
	}

	if err := h.service.Stop(c.Request().Context(), id, timeout); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// Logs handles GET /api/containers/:id/logs.
func (h *ContainerHandler) Logs(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "container ID is required")
	}

	// Parse optional lines parameter
	lines := 100 // Default 100 lines
	if l := c.QueryParam("lines"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			lines = parsed
		}
	}

	logs, err := h.service.Logs(c.Request().Context(), id, lines)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, logs)
}
