package http

import (
	"github.com/igwedaniel/artizan/internal/adapters/http/handlers"
	"github.com/igwedaniel/artizan/internal/adapters/http/middleware"
	"github.com/igwedaniel/artizan/internal/services"
	"github.com/labstack/echo/v4"
)

type Services struct {
	AuthService *services.AuthService
	UserService *services.UserService
	// Add more services here as needed
}

// NewServer creates and configures an Echo server
func NewServer(svcs *Services) *echo.Echo {

	e := echo.New()
	authHandler := handlers.NewAuthHandler(svcs.AuthService)
	userHandler := handlers.NewUserHandler(svcs.UserService)

	// Public routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	e.POST("/auth/nonce", authHandler.GetNonce)
	e.POST("/auth/login", authHandler.Authenticate)
	e.POST("/auth/refresh", authHandler.RefreshToken)

	// Protected routes
	g := e.Group("", middleware.AuthMiddleware(svcs.AuthService))
	g.GET("/me", userHandler.GetCurrentUser)

	return e
}
