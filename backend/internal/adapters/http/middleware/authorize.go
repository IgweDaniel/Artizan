package middleware

import (
	"net/http"
	"strings"

	"github.com/igwedaniel/artizan/internal/services"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware returns an echo middleware that verifies JWT tokens using AuthService
func AuthMiddleware(authService *services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"})
			}
			user, err := authService.VerifyToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}
			// Store user in context for handlers to use
			c.Set("user", user)
			return next(c)
		}
	}
}
