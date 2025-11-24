package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"it_rabotyagi/api/openapi"
	"it_rabotyagi/internal/business/services"
	"it_rabotyagi/internal/data/repositories"
)

// RegisterRoutes регистрирует все маршруты и Swagger
func RegisterRoutes(e *echo.Echo, authService *services.AuthService, repo *repositories.UserRepository, sessionRepo *repositories.SessionRepository, questionRepo *repositories.QuestionRepository, courseRepo *repositories.CourseRepository, mentorRepo *repositories.MentorRepository, permissionRepo *repositories.PermissionRepository) error {
	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Специальный обработчик для openapi.yaml с отключенным кешем
	e.GET("/api-docs/openapi.yaml", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Response().Header().Set("Pragma", "no-cache")
		c.Response().Header().Set("Expires", "0")
		return c.File("api/openapi/openapi.yaml")
	})

	// Статические файлы OpenAPI (остальные файлы)
	e.Static("/api-docs", "api/openapi")

	// Swagger UI - используем кастомный index.html из api/openapi
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/api-docs/index.html")
	})

	// Создаем реализацию обработчиков
	impl := NewServerImplementation(authService, repo, sessionRepo, questionRepo, courseRepo, mentorRepo, permissionRepo)

	api := e.Group("/api/v1")
	wrapper := openapi.ServerInterfaceWrapper{Handler: impl}

	// Публичные маршруты (без авторизации)
	api.POST("/auth/register", wrapper.RegisterUser)
	api.POST("/auth/login", wrapper.LoginUser)
	api.POST("/auth/refresh", wrapper.RefreshTokens)
	api.GET("/questions", wrapper.ListQuestions)
	api.GET("/questions/:id", wrapper.GetQuestionById)
	api.GET("/courses", wrapper.ListCourses)
	api.GET("/courses/:id", wrapper.GetCourseById)
	api.GET("/courses/:id/modules", wrapper.ListModules)
	api.GET("/courses/:id/modules/:moduleId", wrapper.GetModuleById)
	api.GET("/mentors", wrapper.ListMentors)
	api.GET("/mentors/:id", wrapper.GetMentorById)

	// Маршруты с обязательной авторизацией
	authRequired := api.Group("")
	authRequired.Use(AuthMiddleware(authService))
	authRequired.GET("/users/me", wrapper.GetCurrentUser)
	authRequired.POST("/questions", wrapper.CreateQuestion)
	authRequired.PATCH("/questions/:id", wrapper.UpdateQuestion)
	authRequired.DELETE("/questions/:id", wrapper.DeleteQuestion)

	authRequired.POST("/courses", wrapper.CreateCourse)
	authRequired.PATCH("/courses/:id", wrapper.UpdateCourse)
	authRequired.DELETE("/courses/:id", wrapper.DeleteCourse)
	authRequired.POST("/courses/:id/modules", wrapper.CreateModule)
	authRequired.PATCH("/courses/:id/modules/:moduleId", wrapper.UpdateModule)
	authRequired.DELETE("/courses/:id/modules/:moduleId", wrapper.DeleteModule)
	authRequired.POST("/courses/:id/enroll", wrapper.EnrollCourse)
	authRequired.GET("/courses/:id/progress", wrapper.GetCourseProgress)
	authRequired.POST("/courses/:courseId/modules/:moduleId/complete", wrapper.CompleteModule)

	return nil
}
