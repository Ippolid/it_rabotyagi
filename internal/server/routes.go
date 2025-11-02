package server

import (
	"it_rabotyagi/api/openapi"
	"it_rabotyagi/internal/business/services"
	"it_rabotyagi/internal/data/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterRoutes регистрирует все маршруты и Swagger
func RegisterRoutes(e *echo.Echo, authService *services.AuthService, repo *repositories.UserRepository, sessionRepo *repositories.SessionRepository) error {
	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Статические файлы OpenAPI (YAML и другие файлы)
	e.Static("/api-docs", "api/openapi")

	// Swagger UI - используем кастомный index.html из api/openapi
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/api-docs/index.html")
	})

	// Создаем реализацию обработчиков
	impl := NewServerImplementation(authService, repo, sessionRepo)

	// Регистрируем обработчики через обертку
	wrapper := openapi.ServerInterfaceWrapper{Handler: impl}

	// Публичные маршруты (без авторизации)
	e.POST("/api/v1/auth/register", wrapper.RegisterUser)
	e.POST("/api/v1/auth/login", wrapper.LoginUser)
	e.POST("/api/v1/auth/refresh", wrapper.RefreshTokens)

	// Защищенные маршруты (требуют авторизации)
	authRequired := e.Group("/api/v1")
	authRequired.Use(AuthMiddleware(authService))
	authRequired.GET("/users/me", wrapper.GetCurrentUser)

	// Маршруты с опциональной авторизацией
	optionalAuth := e.Group("/api/v1")
	optionalAuth.Use(OptionalAuthMiddleware(authService))
	optionalAuth.GET("/mentors", wrapper.ListMentors)

	return nil
}
