// Package handlers provides HTTP handlers for the Sabakan API.
package handlers

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/sweetfish329/sabakan/backend/internal/games"
	"github.com/sweetfish329/sabakan/backend/internal/middleware"
	"github.com/sweetfish329/sabakan/backend/internal/models"
	"gorm.io/gorm"
)

// slugPattern validates slug format (lowercase letters, numbers, hyphens).
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// CreateGameServerRequest represents the request body for creating a game server.
type CreateGameServerRequest struct {
	Slug        string                  `json:"slug"`
	Name        string                  `json:"name"`
	Game        string                  `json:"game"`
	Description string                  `json:"description,omitempty"`
	Ports       []GameServerPortRequest `json:"ports,omitempty"`
	Envs        []GameServerEnvRequest  `json:"envs,omitempty"`
}

// GameServerPortRequest represents a port mapping in the request.
type GameServerPortRequest struct {
	HostPort      int    `json:"hostPort"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol,omitempty"`
}

// GameServerEnvRequest represents an environment variable in the request.
type GameServerEnvRequest struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret,omitempty"`
}

// UpdateGameServerRequest represents the request body for updating a game server.
type UpdateGameServerRequest struct {
	Name        string                  `json:"name,omitempty"`
	Description string                  `json:"description,omitempty"`
	Ports       []GameServerPortRequest `json:"ports,omitempty"`
	Envs        []GameServerEnvRequest  `json:"envs,omitempty"`
}

// GameServerHandler handles game server-related HTTP requests.
type GameServerHandler struct {
	db *gorm.DB
}

// NewGameServerHandler creates a new game server handler.
func NewGameServerHandler(db *gorm.DB) *GameServerHandler {
	return &GameServerHandler{db: db}
}

// Create handles POST /api/game-servers.
func (h *GameServerHandler) Create(c echo.Context) error {
	var req CreateGameServerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate slug
	if req.Slug == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Slug is required",
		})
	}

	if !slugPattern.MatchString(req.Slug) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Slug must contain only lowercase letters, numbers, and hyphens",
		})
	}

	// Validate name
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Name is required",
		})
	}

	// Check for duplicate slug
	var existing models.GameServer
	if err := h.db.Where("slug = ?", req.Slug).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "conflict",
			Message: "A server with this slug already exists",
		})
	}

	// Get game handler for default image
	image := "unknown:latest"
	if gameHandler, ok := games.Get(req.Game); ok {
		image = gameHandler.GetDefaultImage()
	}

	// Get user ID from context
	userID := middleware.GetUserID(c)

	// Create game server
	server := models.GameServer{
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		Image:       image,
		Status:      models.GameServerStatusStopped,
		OwnerID:     userID,
	}

	if err := h.db.Create(&server).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create game server",
		})
	}

	// Create ports
	for _, p := range req.Ports {
		protocol := p.Protocol
		if protocol == "" {
			protocol = "tcp"
		}
		port := models.GameServerPort{
			GameServerID:  server.ID,
			HostPort:      p.HostPort,
			ContainerPort: p.ContainerPort,
			Protocol:      protocol,
		}
		h.db.Create(&port)
	}

	// Create envs
	for _, e := range req.Envs {
		env := models.GameServerEnv{
			GameServerID: server.ID,
			Key:          e.Key,
			Value:        e.Value,
			IsSecret:     e.IsSecret,
		}
		h.db.Create(&env)
	}

	// Reload with associations
	h.db.Preload("Ports").Preload("Envs").First(&server, server.ID)

	return c.JSON(http.StatusCreated, server)
}

// List handles GET /api/game-servers.
func (h *GameServerHandler) List(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var servers []models.GameServer
	if err := h.db.Where("owner_id = ?", userID).
		Preload("Ports").
		Preload("Envs").
		Find(&servers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to list game servers",
		})
	}

	return c.JSON(http.StatusOK, servers)
}

// Get handles GET /api/game-servers/:slug.
func (h *GameServerHandler) Get(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Slug is required",
		})
	}

	var server models.GameServer
	if err := h.db.Where("slug = ?", slug).
		Preload("Ports").
		Preload("Envs").
		First(&server).Error; err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Game server not found",
		})
	}

	return c.JSON(http.StatusOK, server)
}

// Update handles PUT /api/game-servers/:slug.
func (h *GameServerHandler) Update(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Slug is required",
		})
	}

	var server models.GameServer
	if err := h.db.Where("slug = ?", slug).First(&server).Error; err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Game server not found",
		})
	}

	var req UpdateGameServerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Update fields
	if req.Name != "" {
		server.Name = req.Name
	}
	if req.Description != "" {
		server.Description = req.Description
	}

	if err := h.db.Save(&server).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update game server",
		})
	}

	// Reload with associations
	h.db.Preload("Ports").Preload("Envs").First(&server, server.ID)

	return c.JSON(http.StatusOK, server)
}

// Delete handles DELETE /api/game-servers/:slug.
func (h *GameServerHandler) Delete(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Slug is required",
		})
	}

	var server models.GameServer
	if err := h.db.Where("slug = ?", slug).First(&server).Error; err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Game server not found",
		})
	}

	// Delete associated ports and envs
	h.db.Where("game_server_id = ?", server.ID).Delete(&models.GameServerPort{})
	h.db.Where("game_server_id = ?", server.ID).Delete(&models.GameServerEnv{})

	// Delete the server
	if err := h.db.Delete(&server).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete game server",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
