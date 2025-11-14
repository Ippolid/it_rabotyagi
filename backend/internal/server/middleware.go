package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"it_rabotyagi/internal/business/services"
)

const (
	UserIDKey   = "user_id"
	NicknameKey = "nickname"
	RoleKey     = "role"
)

// AuthMiddleware создает middleware для проверки JWT токенов
func AuthMiddleware(authService *services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Извлекаем токен из заголовка Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Authorization header required",
					"code":    "UNAUTHORIZED",
				})
			}

			// Проверяем формат "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Invalid authorization header format",
					"code":    "INVALID_AUTH_HEADER",
				})
			}

			token := parts[1]

			// Валидируем токен
			claims, err := authService.ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Invalid or expired token",
					"code":    "INVALID_TOKEN",
				})
			}

			// Сохраняем данные пользователя в контексте
			c.Set(UserIDKey, claims.UserID)
			c.Set(NicknameKey, claims.Nickname)
			c.Set(RoleKey, claims.Role)

			return next(c)
		}
	}
}

// OptionalAuthMiddleware - middleware для опциональной аутентификации
// Если токен предоставлен и валиден, данные пользователя добавляются в контекст
// Если токена нет или он невалиден, запрос продолжается без данных пользователя
func OptionalAuthMiddleware(authService *services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token := parts[1]
					claims, err := authService.ValidateToken(token)
					if err == nil {
						c.Set(UserIDKey, claims.UserID)
						c.Set(NicknameKey, claims.Nickname)
						c.Set(RoleKey, claims.Role)
					}
				}
			}
			return next(c)
		}
	}
}

// GetUserID извлекает ID пользователя из контекста
func GetUserID(c echo.Context) (int, bool) {
	userID, ok := c.Get(UserIDKey).(int)
	return userID, ok
}

// GetNickname извлекает никнейм пользователя из контекста
func GetNickname(c echo.Context) (string, bool) {
	nickname, ok := c.Get(NicknameKey).(string)
	return nickname, ok
}

// GetRole извлекает роль пользователя из контекста
func GetRole(c echo.Context) (string, bool) {
	role, ok := c.Get(RoleKey).(string)
	return role, ok
}

// RoleMiddleware проверяет роль пользователя
func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := GetRole(c)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"message": "Role information not found",
					"code":    "FORBIDDEN",
				})
			}

			// Проверяем, есть ли роль пользователя в списке разрешенных
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "Insufficient permissions",
				"code":    "FORBIDDEN",
			})
		}
	}
}
