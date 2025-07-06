package handlers

import (
	"net/http"

	"github.com/igwedaniel/artizan/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) GetNonce(c echo.Context) error {
	var req struct {
		WalletAddress string `json:"wallet_address"`
	}
	if err := c.Bind(&req); err != nil || req.WalletAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	msg, err := h.AuthService.GetNonceMessage(req.WalletAddress)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": msg})
}

func (h *AuthHandler) Authenticate(c echo.Context) error {
	var req struct {
		WalletAddress string `json:"wallet_address"`
		Signature     string `json:"signature"`
	}
	if err := c.Bind(&req); err != nil || req.WalletAddress == "" || req.Signature == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	tokens, user, err := h.AuthService.Authenticate(req.WalletAddress, req.Signature)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"tokens": tokens, "user": user})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind(&req); err != nil || req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	tokens, user, err := h.AuthService.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"tokens": tokens, "user": user})
}
