package handlers

import (
	"net/http"

	"github.com/igwedaniel/artizan/internal/services"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService *services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// GET /me (protected)
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	return c.JSON(http.StatusOK, user)
}
