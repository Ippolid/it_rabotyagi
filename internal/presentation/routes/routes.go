package routes

import (
	"itpath/internal/business/services"
	"itpath/internal/config"
	"itpath/internal/presentation/handlers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает и возвращает роутер
func SetupRoutes(cfg *config.Config, authService *services.AuthService, oauthService *services.OAuthService) *gin.Engine {
	router := gin.Default()

	// Настройка middleware
	setupMiddleware(router)

	// Настройка статических файлов
	setupStaticFiles(router)

	// Инициализируем хендлеры
	authHandler := handlers.NewAuthHandler(authService, oauthService)

	// Настройка OAuth роутов
	setupOAuthRoutes(router, authHandler)

	// Настройка API роутов
	setupAPIRoutes(router, authService, authHandler)

	return router
}

// setupMiddleware настраивает middleware для роутера
func setupMiddleware(router *gin.Engine) {
	// CORS middleware
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:8080"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// setupStaticFiles настраивает раздачу статических файлов
func setupStaticFiles(router *gin.Engine) {
	router.Static("/web", "./web")
	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})
}

// setupOAuthRoutes монтирует OAuth роуты
func setupOAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	{
		// GitHub OAuth
		auth.GET("/github/login", authHandler.GitHubLogin)
		auth.GET("/github/callback", authHandler.GitHubCallback)

		// Google OAuth
		auth.GET("/google/login", authHandler.GoogleLogin)
		auth.GET("/google/callback", authHandler.GoogleCallback)

		// Logout
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/logout", authHandler.Logout)
	}
}

// setupAPIRoutes настраивает API эндпоинты
func setupAPIRoutes(router *gin.Engine, authService *services.AuthService, authHandler *handlers.AuthHandler) {
	api := router.Group("/api/v1")
	{
		// ============================================================================
		// Публичные эндпоинты (не требуют авторизации)
		// ============================================================================
		public := api.Group("/public")
		{
			public.GET("/users/:id", authHandler.GetUserByID)
		}

		// ============================================================================
		// Защищенные эндпоинты (требуют авторизации)
		// ============================================================================
		authorized := api.Group("")
		authorized.Use(AuthMiddleware(authService))
		{
			// Профиль текущего пользователя
			authorized.GET("/me", authHandler.GetMe)
			authorized.PUT("/me", authHandler.UpdateProfile)
			authorized.POST("/logout", authHandler.Logout)

			// Управление связанными аккаунтами
			accounts := authorized.Group("/accounts")
			{
				// GitHub
				accounts.POST("/github/link", authHandler.LinkGitHub)
				accounts.DELETE("/github/unlink", authHandler.UnlinkGitHub)

				// Google
				accounts.POST("/google/link", authHandler.LinkGoogle)
				accounts.DELETE("/google/unlink", authHandler.UnlinkGoogle)

				// Telegram
				accounts.POST("/telegram/link", authHandler.LinkTelegram)
				accounts.DELETE("/telegram/unlink", authHandler.UnlinkTelegram)
			}
		}
	}
}

// AuthMiddleware проверяет JWT токен
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пытаемся получить токен из cookie
		tokenString, err := c.Cookie("jwt")

		// Если в cookie нет, пробуем Authorization header
		if err != nil || tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				// Ожидаем формат: "Bearer <token>"
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Authentication required",
			})
			return
		}

		// Валидируем токен
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
			})
			return
		}

		// Сохраняем claims в контекст
		c.Set("claims", claims)
		c.Next()
	}
}
