package handlers

import (
	"itpath/internal/business/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pkgz/auth/token"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// GetMe возвращает информацию о текущем пользователе
func (h *AuthHandler) GetMe(c *gin.Context) {
	// Получаем user из контекста (их установит middleware go-pkgz/auth)
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tokenUser, ok := userValue.(token.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token user"})
		return
	}

	// Получаем пользователя из БД по данным токена
	user, err := h.authService.GetUserFromToken(tokenUser)
	if err != nil {
		log.Printf("Error getting user from token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

// GetUserByID возвращает информацию о пользователе по ID
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.authService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

// Logout выполняет выход пользователя
func (h *AuthHandler) Logout(c *gin.Context) {
	// go-pkgz/auth обрабатывает logout автоматически
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "logged out successfully",
	})
}
