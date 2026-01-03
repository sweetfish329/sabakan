package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sweetfish329/sabakan/backend/internal/models"
	"gorm.io/gorm"
)

// ModHandler handles mod-related HTTP requests.
type ModHandler struct {
	db *gorm.DB
}

// NewModHandler creates a new mod handler.
func NewModHandler(db *gorm.DB) *ModHandler {
	return &ModHandler{db: db}
}

// CreateModRequest represents the request body for creating a mod.
type CreateModRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	SourceURL   string `json:"sourceUrl"`
	Version     string `json:"version"`
}

// UpdateModRequest represents the request body for updating a mod.
type UpdateModRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	SourceURL   *string `json:"sourceUrl"`
	Version     *string `json:"version"`
}

// List handles GET /api/mods.
func (h *ModHandler) List(c echo.Context) error {
	var mods []models.Mod
	if err := h.db.Find(&mods).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch mods")
	}
	return c.JSON(http.StatusOK, mods)
}

// Get handles GET /api/mods/:id.
func (h *ModHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid mod ID")
	}

	var mod models.Mod
	if err := h.db.First(&mod, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Mod not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch mod")
	}

	return c.JSON(http.StatusOK, mod)
}

// Create handles POST /api/mods.
func (h *ModHandler) Create(c echo.Context) error {
	var req CreateModRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate required fields
	if req.Name == "" || req.Slug == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name and slug are required")
	}

	mod := models.Mod{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SourceURL:   req.SourceURL,
		Version:     req.Version,
	}

	if err := h.db.Create(&mod).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create mod")
	}

	return c.JSON(http.StatusCreated, mod)
}

// Update handles PUT /api/mods/:id.
func (h *ModHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid mod ID")
	}

	var mod models.Mod
	if err := h.db.First(&mod, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Mod not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch mod")
	}

	var req UpdateModRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Update only provided fields
	if req.Name != nil {
		mod.Name = *req.Name
	}
	if req.Slug != nil {
		mod.Slug = *req.Slug
	}
	if req.Description != nil {
		mod.Description = *req.Description
	}
	if req.SourceURL != nil {
		mod.SourceURL = *req.SourceURL
	}
	if req.Version != nil {
		mod.Version = *req.Version
	}

	if err := h.db.Save(&mod).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update mod")
	}

	return c.JSON(http.StatusOK, mod)
}

// Delete handles DELETE /api/mods/:id.
func (h *ModHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid mod ID")
	}

	var mod models.Mod
	if err := h.db.First(&mod, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Mod not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch mod")
	}

	if err := h.db.Delete(&mod).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete mod")
	}

	return c.NoContent(http.StatusNoContent)
}
