package handlers

import (
	"itpath/internal/business/services"
	"itpath/internal/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService  *services.AuthService
	oauthService *services.OAuthService
}

func NewAuthHandler(authService *services.AuthService, oauthService *services.OAuthService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		oauthService: oauthService,
	}
}

// ============================================================================
// Базовые хендлеры для работы с пользователями
// ============================================================================

// GetMe возвращает информацию о текущем пользователе
func (h *AuthHandler) GetMe(c *gin.Context) {
	// Получаем claims из контекста (их установит middleware)
	claimsValue, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := claimsValue.(*services.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	// Получаем пользователя из БД
	user, err := h.authService.GetUserByID(claims.UserID)
	if err != nil {
		logger.Error("Error getting user", zap.Error(err))
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

// UpdateProfile обновляет профиль текущего пользователя
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Получаем claims из контекста
	claimsValue, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := claimsValue.(*services.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	user, err := h.authService.GetUserByID(claims.UserID)
	if err != nil {
		logger.Error("Error getting user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user data"})
		return
	}

	// Парсим данные для обновления
	var updateData struct {
		Name        *string `json:"name"`
		Username    *string `json:"username"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	// Обновляем поля
	if updateData.Name != nil {
		user.Name = *updateData.Name
	}
	if updateData.Username != nil {
		user.Username = updateData.Username
	}
	if updateData.Description != nil {
		user.Description = updateData.Description
	}

	// Сохраняем изменения
	if err := h.authService.UpdateUser(user); err != nil {
		logger.Error("Error updating user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

// Logout выполняет выход пользователя
func (h *AuthHandler) Logout(c *gin.Context) {
	// Удаляем JWT cookie
	c.SetCookie("jwt", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "logged out successfully",
	})
}

// ============================================================================
// OAuth хендлеры
// ============================================================================

// GitHubLogin начинает процесс авторизации через GitHub
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	authURL, err := h.oauthService.GetGitHubAuthURL()
	if err != nil {
		logger.Error("Error getting GitHub auth URL", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate auth URL"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GitHubCallback обрабатывает callback от GitHub
func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code parameter is missing"})
		return
	}

	// Получаем данные пользователя от GitHub
	githubUser, err := h.oauthService.HandleGitHubCallback(c.Request.Context(), state, code)

	logger.Debug("user", zap.Any("user", githubUser))

	if err != nil {
		logger.Error("Error handling GitHub callback", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate with GitHub"})
		return
	}

	// Получаем или создаем пользователя в БД
	user, err := h.authService.GetOrCreateUserFromGitHub(githubUser)
	if err != nil {
		logger.Error("Error creating/getting user from GitHub", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Генерируем JWT токен
	token, err := h.authService.GenerateToken(user, "github")
	if err != nil {
		logger.Error("Error generating token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Устанавливаем cookie
	c.SetCookie("jwt", token, 60*60*24*7, "/", "", false, true) // 7 дней

	// Редиректим на фронтенд
	c.Redirect(http.StatusTemporaryRedirect, "/?auth=success")
}

// GoogleLogin начинает процесс авторизации через Google
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	authURL, err := h.oauthService.GetGoogleAuthURL()
	if err != nil {
		logger.Error("Error getting Google auth URL", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate auth URL"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback обрабатывает callback от Google
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code parameter is missing"})
		return
	}

	// Получаем данные пользователя от Google
	googleUser, err := h.oauthService.HandleGoogleCallback(c.Request.Context(), state, code)
	if err != nil {
		logger.Error("Error handling Google callback", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate with Google"})
		return
	}

	// Получаем или создаем пользователя в БД
	user, err := h.authService.GetOrCreateUserFromGoogle(googleUser)
	if err != nil {
		logger.Error("Error creating/getting user from Google", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Генерируем JWT токен
	token, err := h.authService.GenerateToken(user, "google")
	if err != nil {
		logger.Error("Error generating token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Устанавливаем cookie
	c.SetCookie("jwt", token, 60*60*24*7, "/", "", false, true) // 7 дней

	// Редиректим на фронтенд
	c.Redirect(http.StatusTemporaryRedirect, "/?auth=success")
}

// ============================================================================
// Хендлеры для линковки аккаунтов
// ============================================================================

// LinkGitHub связывает GitHub аккаунт с текущим пользователем
func (h *AuthHandler) LinkGitHub(c *gin.Context) {
	// TODO: Implement account linking flow
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented yet"})
}

// UnlinkGitHub отвязывает GitHub аккаунт от текущего пользователя
func (h *AuthHandler) UnlinkGitHub(c *gin.Context) {
	claimsValue, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := claimsValue.(*services.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	if err := h.authService.UnlinkGitHubAccount(claims.UserID); err != nil {
		logger.Error("Error unlinking GitHub account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unlink GitHub account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "GitHub account unlinked successfully",
	})
}

// LinkGoogle связывает Google аккаунт с текущим пользователем
func (h *AuthHandler) LinkGoogle(c *gin.Context) {
	// TODO: Implement account linking flow
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented yet"})
}

// UnlinkGoogle отвязывает Google аккаунт от текущего пользователя
func (h *AuthHandler) UnlinkGoogle(c *gin.Context) {
	claimsValue, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := claimsValue.(*services.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	if err := h.authService.UnlinkGoogleAccount(claims.UserID); err != nil {
		logger.Error("Error unlinking Google account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unlink Google account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Google account unlinked successfully",
	})
}

// LinkTelegram связывает Telegram аккаунт с текущим пользователем
func (h *AuthHandler) LinkTelegram(c *gin.Context) {
	claimsValue, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := claimsValue.(*services.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	// Парсим данные Telegram
	var linkData struct {
		TelegramID string                 `json:"telegram_id" binding:"required"`
		Data       map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&linkData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	// Линкуем аккаунт
	if err := h.authService.LinkTelegramAccount(claims.UserID, linkData.TelegramID, linkData.Data); err != nil {
		logger.Error("Error linking Telegram account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func (h *AuthHandler) UnlinkTelegram(context *gin.Context) {

}
