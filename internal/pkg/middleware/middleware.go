package middleware

import (
	"itpath/internal/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(jwtManager *jwt.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Добавляем все данные пользователя в контекст
		c.Set("user_id", claims.UserID)
		c.Set("telegram_id", claims.TelegramID)
		c.Set("username", claims.Username)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("role", claims.Role)
		c.Set("subscription", claims.Subscription)

		c.Next()
	}
}

// Middleware для проверки роли
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		if userRole.(string) != role && userRole.(string) != "mentor" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Middleware для проверки подписки
func RequireSubscription(subscription string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userSub, exists := c.Get("subscription")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User subscription not found"})
			c.Abort()
			return
		}

		if subscription == "trial" {
			// trial доступен всем
			c.Next()
			return
		}

		if userSub.(string) != subscription {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Subscription required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
