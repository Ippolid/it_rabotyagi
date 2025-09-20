package routes

import (
	"github.com/gin-gonic/gin"
	"itpath/internal/business"
	"itpath/internal/pkg/jwt"
	"itpath/internal/pkg/middleware"
	"itpath/internal/presentation/handlers"
)

func SetupRoutes(
	authService business.AuthService,
	jwtManager *jwt.TokenManager,
) *gin.Engine {
	r := gin.Default()

	// Global middleware
	//r.Use(middleware.CORS())
	//r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)
	//healthHandler := handlers.NewHealthHandler()

	// API routes
	api := r.Group("/api/v1")

	// Public routes
	{
		// Health check
		//api.GET("/health", healthHandler.Health)

		// Authentication
		auth := api.Group("/auth")
		auth.POST("/telegram", authHandler.TelegramLogin)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.Auth(jwtManager))
	{
		// User endpoints
		protected.GET("/me", userHandler.GetMe)
		protected.PUT("/profile", userHandler.UpdateProfile)

		// Auth endpoints
		protected.POST("/auth/logout", authHandler.Logout)
	}

	// Role-based routes
	mentorOnly := protected.Group("")
	mentorOnly.Use(middleware.RequireRole("mentor"))
	{
		// Mentor-only endpoints
	}

	return r
}
