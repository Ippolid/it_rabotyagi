package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config содержит всю конфигурацию приложения.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Telegram TelegramConfig
}

// ServerConfig содержит конфигурацию HTTP-сервера.
type ServerConfig struct {
	Port string `envconfig:"SERVER_PORT" default:"8080"`
}

// DatabaseConfig содержит конфигурацию подключения к базе данных.
type DatabaseConfig struct {
	URL string `envconfig:"DATABASE_URL" required:"true"`
}

// JWTConfig содержит конфигурацию для JWT.
type JWTConfig struct {
	Secret string `envconfig:"JWT_SECRET" required:"true"`
}

// TelegramConfig содержит конфигурацию для Telegram API.
type TelegramConfig struct {
	BotToken string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
}

// Load загружает конфигурацию из файла deploy/.env и переменных окружения.
func Load() (*Config, error) {
	// Загружаем переменные из файла .env.
	// Ошибку можно игнорировать, если файл не найден,
	// так как переменные могут быть установлены напрямую в среде выполнения.
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found at 'deploy/.env': %v", err)
	}

	var cfg Config
	// Парсим переменные окружения в структуру Config.
	err = envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	return &cfg, nil
}
