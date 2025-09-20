package app

import (
	"context"
	"fmt"
	"itpath/internal/business/services"
	"itpath/internal/data/repositories"
	"itpath/internal/pkg/jwt"
	"itpath/internal/presentation/routes"
	"log"
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

	// Подключаемся к базе данных
	db, err := database.NewPostgresConnection(cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Инициализируем зависимости слой за слоем

	// DATA LAYER
	userRepo := repositories.NewUserRepository(db)

	// BUSINESS LAYER
	jwtManager := jwt.NewManager(cfg.JWT.Secret)
	telegramService := services.NewTelegramService(cfg.Telegram.BotToken)
	authService := services.NewAuthService(userRepo, telegramService, jwtManager)

	// PRESENTATION LAYER
	router := routes.SetupRoutes(authService, jwtManager)

	// HTTP сервер
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("✅ Server exited")
	return nil
}
