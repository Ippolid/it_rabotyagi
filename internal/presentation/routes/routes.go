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
	// r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)

	// API routes
	r.Static("/web", "./web")
	api := r.Group("/api/v1")

	// Public routes
	{
		// Authentication
		auth := api.Group("/auth")
		auth.POST("/telegram", authHandler.TelegramLogin)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.Auth(jwtManager))
	{
		// Auth endpoints
		protected.POST("/auth/logout", authHandler.Logout)

		// User endpoints
		protected.GET("/me", userHandler.GetMe)
		protected.PUT("/profile", userHandler.UpdateProfile)
	}

	// Role-based routes
	mentorOnly := protected.Group("")
	mentorOnly.Use(middleware.RequireRole("mentor"))
	{
		// Mentor-only endpoints can be added here
	}

	return r
}
