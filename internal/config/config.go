package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Logger LoggerConfig
}

type ServerConfig struct {
	Port      string
	Host      string
	PublicURL string
}

type LoggerConfig struct {
	Level string
}

type DatabaseConfig struct {
	URL string
}

type AuthConfig struct {
	Secret          string
	TokenDuration   int // в минутах
	RefreshDuration int // в минутах
	Telegram        TelegramConfig
	Google          GoogleConfig
	GitHub          GitHubConfig
}

type TelegramConfig struct {
	Token string
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
}

type GitHubConfig struct {
	ClientID     string
	ClientSecret string
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:      getEnv("SERVER_PORT", "8080"),
			Host:      getEnv("HTTP_HOST", "0.0.0.0"),
			PublicURL: getEnv("PUBLIC_URL", "http://localhost:8080"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		Auth: AuthConfig{
			Secret:          getEnv("JWT_SECRET", ""),
			TokenDuration:   60,   // 1 час
			RefreshDuration: 4320, // 3 дня
			Telegram: TelegramConfig{
				Token: getEnv("TELEGRAM_BOT_TOKEN", ""),
			},
			Google: GoogleConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			},
			GitHub: GitHubConfig{
				ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
				ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
			},
		},
		Logger: LoggerConfig{
			Level: getEnv("LOGGER_LEVEL", "info"),
		},
	}

	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.Auth.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
