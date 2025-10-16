package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"itpath/internal/business/services"
	"itpath/internal/config"
	"itpath/internal/data/database"
	"itpath/internal/data/repositories"
	"itpath/internal/logger"
	"itpath/internal/presentation/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	logger.Info("Initializing UserRepo...")
	
	// BUSINESS LAYER
	authService := services.NewAuthService(userRepo, cfg.Auth.Secret)
	logger.Info("Initializing AuthService...")
	
	// OAuth Service
	oauthService := services.NewOAuthService(
		cfg.Auth.GitHub.ClientID,
		cfg.Auth.GitHub.ClientSecret,
		cfg.Auth.Google.ClientID,
		cfg.Auth.Google.ClientSecret,
		cfg.Server.PublicURL,
	)
	logger.Info("Initializing OAuthService...")
	
	// PRESENTATION LAYER
	router := routes.SetupRoutes(cfg, authService, oauthService)
	logger.Info("Initializing Routes...")

	// HTTP сервер
	server := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		logger.Info("Server starting",
			zap.String("url", fmt.Sprintf("http://%s", server.Addr)),
			zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("Server gracefully stopped")
	return nil
}
