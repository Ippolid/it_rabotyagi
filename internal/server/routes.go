package server

import (
	"fmt"
	"it_rabotyagi/api/openapi"
	"it_rabotyagi/internal/business/services"
	"it_rabotyagi/internal/data/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterRoutes регистрирует все маршруты и Swagger
<<<<<<< Updated upstream
func RegisterRoutes(e *echo.Echo, authService *services.AuthService, repo *repositories.UserRepository, sessionRepo *repositories.SessionRepository) error {
=======
func RegisterRoutes(
	e *echo.Echo,
	authService *services.AuthService,
	repo *repositories.UserRepository,
	sessionRepo *repositories.SessionRepository,
	courseRepo *repositories.CourseRepository,
	moduleRepo *repositories.ModuleRepository,
	questionRepo *repositories.QuestionRepository,
) error {
>>>>>>> Stashed changes
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

<<<<<<< Updated upstream
=======
	// Курсы (пока вне openapi, простой список опубликованных)
	e.GET("/api/v1/courses", func(c echo.Context) error {
		courses, err := courseRepo.ListPublished(c.Request().Context(), 12, 0)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to load courses"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"items": courses})
	})

	// Детали курса по id
	e.GET("/api/v1/courses/:id", func(c echo.Context) error {
		idParam := c.Param("id")
		var id int64
		if _, err := fmt.Sscan(idParam, &id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid course id"})
		}
		course, err := courseRepo.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "course not found"})
		}
		return c.JSON(http.StatusOK, course)
	})

	// Модули курса
	e.GET("/api/v1/courses/:id/modules", func(c echo.Context) error {
		idParam := c.Param("id")
		var id int64
		if _, err := fmt.Sscan(idParam, &id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid course id"})
		}
		modules, err := moduleRepo.ListByCourse(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to load modules"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"items": modules})
	})

	// Детали модуля
	e.GET("/api/v1/modules/:id", func(c echo.Context) error {
		idParam := c.Param("id")
		var id int64
		if _, err := fmt.Sscan(idParam, &id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid module id"})
		}
		m, err := moduleRepo.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "module not found"})
		}
		return c.JSON(http.StatusOK, m)
	})

	// Вопросы модуля
	e.GET("/api/v1/modules/:id/questions", func(c echo.Context) error {
		idParam := c.Param("id")
		var id int64
		if _, err := fmt.Sscan(idParam, &id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid module id"})
		}
		questions, err := questionRepo.ListByModule(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to load questions"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"items": questions})
	})

>>>>>>> Stashed changes
	return nil
}
