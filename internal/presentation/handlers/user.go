package handlers

import (
	"errors"
	"fmt"
	"itpath/internal/business"
	"itpath/internal/data/entities"
	"itpath/internal/pkg/response"
	"itpath/internal/presentation/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	authService business.AuthService
}

func NewUserHandler(authService business.AuthService) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

// GET /api/v1/me
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("telegram_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	fmt.Println(userID, c)
	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(int64))
	fmt.Println(user)
	if err != nil {
		if errors.Is(err, entities.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "User not found", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve user", err)
		return
	}

	response.Success(c, user.ToPublic(), "User retrieved successfully")
}

// PUT /api/v1/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("telegram_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	updateReq := req.ToBusinessModel()

	user, err := h.authService.UpdateProfile(c.Request.Context(), userID.(int64), updateReq)
	if err != nil {
		// Здесь можно добавить более специфичную обработку ошибок, например, для дубликатов
		response.Error(c, http.StatusBadRequest, "Profile update failed", err)
		return
	}

	response.Success(c, user.ToPublic(), "Profile updated successfully")
}
