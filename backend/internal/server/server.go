// Package server provides the HTTP server for Sabakan.
package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sweetfish329/sabakan/backend/internal/container"
	"github.com/sweetfish329/sabakan/backend/internal/handlers"
)

// New creates a new Echo server with all middleware and routes configured.
func New(containerService *container.Service) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Sabakan!")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Container API routes
	api := e.Group("/api")
	containerHandler := handlers.NewContainerHandler(containerService)
	api.GET("/containers", containerHandler.List)
	api.GET("/containers/:id", containerHandler.Get)
	api.POST("/containers/:id/start", containerHandler.Start)
	api.POST("/containers/:id/stop", containerHandler.Stop)
	api.GET("/containers/:id/logs", containerHandler.Logs)

	return e
}
