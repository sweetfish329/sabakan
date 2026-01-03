// Package server provides the HTTP server for Sabakan.
package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sweetfish329/sabakan/backend/internal/auth"
	"github.com/sweetfish329/sabakan/backend/internal/config"
	"github.com/sweetfish329/sabakan/backend/internal/container"
	"github.com/sweetfish329/sabakan/backend/internal/handlers"
	"github.com/sweetfish329/sabakan/backend/internal/middleware"
	"github.com/sweetfish329/sabakan/backend/internal/redis"
	"gorm.io/gorm"
)

// Dependencies holds all the dependencies needed by the server.
type Dependencies struct {
	ContainerService *container.Service
	DB               *gorm.DB
	Config           *config.SystemConfig
	SessionStore     redis.SessionStore
}

// New creates a new Echo server with all middleware and routes configured.
func New(deps *Dependencies) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Sabakan!")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Create JWT manager
	jwtManager := auth.NewJWTManager(
		deps.Config.JWT.Secret,
		time.Duration(deps.Config.JWT.AccessTokenExpiry)*time.Minute,
		time.Duration(deps.Config.JWT.RefreshTokenExpiry)*24*time.Hour,
	)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager, deps.SessionStore)
	permMiddleware := middleware.NewPermissionMiddleware(deps.DB)

	// Auth routes (public)
	authHandler := handlers.NewAuthHandler(deps.DB, jwtManager, deps.SessionStore)
	authGroup := e.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/refresh", authHandler.Refresh)
	authGroup.POST("/logout", authMiddleware.Authenticate(authHandler.Logout))

	// OAuth routes (public)
	oauthHandler := handlers.NewOAuthHandler(
		deps.DB,
		jwtManager,
		deps.SessionStore,
		&deps.Config.OAuth,
		"http://localhost:4200", // Frontend URL
	)
	authGroup.GET("/oauth/:provider", oauthHandler.Authorize)
	authGroup.GET("/oauth/:provider/callback", oauthHandler.Callback)

	// API routes (protected)
	api := e.Group("/api")
	api.Use(authMiddleware.Authenticate)

	// Container routes
	containerHandler := handlers.NewContainerHandler(deps.ContainerService)
	containers := api.Group("/containers")
	containers.GET("", containerHandler.List, permMiddleware.RequirePermission("game_server", "read"))
	containers.GET("/:id", containerHandler.Get, permMiddleware.RequirePermission("game_server", "read"))
	containers.POST("/:id/start", containerHandler.Start, permMiddleware.RequirePermission("game_server", "start"))
	containers.POST("/:id/stop", containerHandler.Stop, permMiddleware.RequirePermission("game_server", "stop"))
	containers.GET("/:id/logs", containerHandler.Logs, permMiddleware.RequirePermission("game_server", "read"))

	// Mod routes
	modHandler := handlers.NewModHandler(deps.DB)
	mods := api.Group("/mods")
	mods.GET("", modHandler.List, permMiddleware.RequirePermission("mod", "read"))
	mods.GET("/:id", modHandler.Get, permMiddleware.RequirePermission("mod", "read"))
	mods.POST("", modHandler.Create, permMiddleware.RequirePermission("mod", "create"))
	mods.PUT("/:id", modHandler.Update, permMiddleware.RequirePermission("mod", "update"))
	mods.DELETE("/:id", modHandler.Delete, permMiddleware.RequirePermission("mod", "delete"))

	return e
}
