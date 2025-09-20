package business

import (
	"context"
	"itpath/internal/business/models"
)

// AuthService - интерфейс сервиса авторизации
type AuthService interface {
	// Авторизация
	AuthenticateWithTelegram(ctx context.Context, data models.TelegramAuthData) (*models.AuthResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResult, error)
	Logout(ctx context.Context, userID int64) error

	// Управление пользователями
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	UpdateProfile(ctx context.Context, userID int64, req models.UpdateProfileRequest) (*models.User, error)

	// Проверки доступа
	ValidateAccess(ctx context.Context, userID int64, requiredRole string) (bool, error)
	ValidateSubscription(ctx context.Context, userID int64, requiredSubscription string) (bool, error)
}

// TelegramService - интерфейс для работы с Telegram
type TelegramService interface {
	ValidateAuthData(data models.TelegramAuthData) error
}
