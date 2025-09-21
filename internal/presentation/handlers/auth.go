package handlers

import (
	"errors"
	"itpath/internal/business"
	"itpath/internal/data/entities"
	"itpath/internal/pkg/response"
	"itpath/internal/presentation/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService business.AuthService
}

func NewAuthHandler(authService business.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// POST /api/v1/auth/telegram
func (h *AuthHandler) TelegramLogin(c *gin.Context) {
	var req dto.TelegramAuthRequest

	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	authData := req.ToBusinessModel()

	authResult, err := h.authService.AuthenticateWithTelegram(c.Request.Context(), authData)
	if err != nil {
		if errors.Is(err, entities.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "User not found", err)
			return
		}
		response.Error(c, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	authResponse := dto.NewAuthResponse(authResult)
	response.Success(c, authResponse, "Successfully authenticated")
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	authResult, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, entities.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "User for token not found", err)
			return
		}
		response.Error(c, http.StatusUnauthorized, "Token refresh failed", err)
		return
	}

	authResponse := dto.NewAuthResponse(authResult)
	response.Success(c, authResponse, "Token refreshed successfully")
}

// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("telegram_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID.(int64)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Logout failed", err)
		return
	}

	response.Success(c, nil, "Successfully logged out")
}
