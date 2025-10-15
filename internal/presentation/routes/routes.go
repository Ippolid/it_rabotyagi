package routes

import (
	"itpath/internal/business/services"
	"itpath/internal/config"
	"itpath/internal/presentation/handlers"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/token"
)

// SimpleLogger реализует интерфейс logger.L для go-pkgz/auth
type SimpleLogger struct{}

func (l SimpleLogger) Logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func SetupRoutes(cfg *config.Config, authService *services.AuthService) (*gin.Engine, *auth.Service) {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-JWT")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Статические файлы
	router.Static("/web", "./web")
	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// Настройка go-pkgz/auth
	authOptions := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) {
			return cfg.Auth.Secret, nil
		}),
		ClaimsUpd:       token.ClaimsUpdFunc(authService.ClaimsUpdater), // Подключаем обновление claims для сохранения в БД
		TokenDuration:   time.Minute * time.Duration(cfg.Auth.TokenDuration),
		CookieDuration:  time.Hour * 24 * 7, // 7 дней
		Issuer:          "itpath",
		URL:             cfg.Server.PublicURL,
		AvatarStore:     avatar.NewNoOp(),
		Logger:          SimpleLogger{},
		DisableXSRF:     true,                     // Отключаем XSRF для упрощения
		SecureCookies:   false,                    // Для HTTP (в продакшене должно быть true)
		SameSiteCookie:  http.SameSiteDefaultMode, // Default для совместимости
		JWTCookieDomain: "",                       // Пустой домен для работы на любом хосте
	}

	authSvc := auth.NewService(authOptions)

	// Добавляем провайдеры аутентификации
	// Telegram
	if cfg.Auth.Telegram.Token != "" {
		authSvc.AddProvider("telegram", cfg.Auth.Telegram.Token, "")
		log.Println("✅ Telegram provider enabled")
	}

	// Google
	if cfg.Auth.Google.ClientID != "" && cfg.Auth.Google.ClientSecret != "" {
		authSvc.AddProvider("google", cfg.Auth.Google.ClientID, cfg.Auth.Google.ClientSecret)
		log.Println("✅ Google provider enabled")
	}

	// GitHub
	if cfg.Auth.GitHub.ClientID != "" && cfg.Auth.GitHub.ClientSecret != "" {
		authSvc.AddProvider("github", cfg.Auth.GitHub.ClientID, cfg.Auth.GitHub.ClientSecret)
		log.Println("✅ GitHub provider enabled")
	}

	// Монтируем auth routes от go-pkgz/auth
	authRoutes, avaRoutes := authSvc.Handlers()

	// Auth handlers
	router.GET("/auth/:provider/login", gin.WrapH(authRoutes))
	router.POST("/auth/:provider/login", gin.WrapH(authRoutes))
	router.GET("/auth/:provider/callback", gin.WrapH(authRoutes))
	router.GET("/auth/logout", gin.WrapH(authRoutes))

	// Avatar handlers
	router.GET("/avatar/:avatar", gin.WrapH(avaRoutes))

	// Инициализируем хендлеры
	authHandler := handlers.NewAuthHandler(authService)

	// API routes
	api := router.Group("/api/v1")
	{
		// Публичные эндпоинты
		api.GET("/users/:id", authHandler.GetUserByID)

		// Защищенные эндпоинты (требуют аутентификации)
		authorized := api.Group("")
		authorized.Use(AuthMiddleware(authSvc))
		{
			authorized.GET("/me", authHandler.GetMe)
			authorized.POST("/logout", authHandler.Logout)
		}
	}

	return router, authSvc
}

// AuthMiddleware проверяет JWT токен через go-pkgz/auth
func AuthMiddleware(authSvc *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем middleware от auth service
		middleware := authSvc.Middleware()

		// Создаем обработчик для проверки аутентификации
		var authenticated bool
		var userInfo token.User

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем user из контекста
			user, err := token.GetUserInfo(r)
			if err != nil {
				authenticated = false
				return
			}

			userInfo = user
			authenticated = true
		})

		// Применяем middleware
		wrappedHandler := middleware.Auth(testHandler)
		wrappedHandler.ServeHTTP(c.Writer, c.Request)

		if !authenticated {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Сохраняем user в контекст Gin
		c.Set("user", userInfo)
		c.Next()
	}
}
