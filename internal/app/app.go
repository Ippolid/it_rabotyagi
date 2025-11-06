package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"it_rabotyagi/internal/business/services"
	"it_rabotyagi/internal/config"
	"it_rabotyagi/internal/data/database"
	"it_rabotyagi/internal/data/repositories"
	"it_rabotyagi/internal/logger"
	"it_rabotyagi/internal/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func Run() error {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger.InitLocalLogger(cfg.Logger.Level)

	logger.Info("Creating new App...")
	// Подключаемся к базе данных
	db, err := database.NewPostgresConnection(cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	logger.Info("Connected to database")

	// Инициализируем зависимости слой за слоем
	// DATA LAYER
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	moduleRepo := repositories.NewModuleRepository(db)
	questionRepo := repositories.NewQuestionRepository(db)
	logger.Info("Initializing UserRepo and SessionRepo...")

	// BUSINESS LAYER
	authService := services.NewAuthService(cfg.Auth.Secret, cfg.Auth.TokenDuration, cfg.Auth.RefreshDuration)
	logger.Info("Initializing AuthService...")

	// PRESENTATION LAYER
	e := echo.New()
	e.HideBanner = true
	if err := server.RegisterRoutes(e, authService, userRepo, sessionRepo, courseRepo, moduleRepo, questionRepo); err != nil {
		return fmt.Errorf("failed to register routes: %w", err)
	}
	logger.Info("Routes registered successfully...")

	// HTTP сервер
	addr := cfg.Server.Host + ":" + cfg.Server.Port

	// Запускаем сервер в горутине
	go func() {
		logger.Info("Server starting",
			zap.String("url", fmt.Sprintf("http://localhost:%s", cfg.Server.Port)),
			zap.String("swagger_url", fmt.Sprintf("http://localhost:%s/api-docs/index.html", cfg.Server.Port)),
			zap.String("address", addr))
		if err := e.Start(addr); err != nil && err.Error() != "http: Server closed" {
			logger.Fatal("Failed to start server:", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("Server gracefully stopped")
	return nil
}
