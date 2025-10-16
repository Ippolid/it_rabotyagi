package business

import (
	"itpath/internal/business/models"

	"github.com/go-pkgz/auth/token"
)

// AuthService - интерфейс сервиса авторизации и управления пользователями
type AuthService interface {
	// ============================================================================
	// Основные методы работы с пользователями
	// ============================================================================

	// GetUserByID получает пользователя по ID из базы данных
	GetUserByID(id int64) (*models.User, error)

	// UpdateUser обновляет данные пользователя в базе данных
	UpdateUser(user *models.User) error

	// ============================================================================
	// OAuth авторизация (GitHub, Google, Telegram)
	// ============================================================================

	// ClaimsUpdater обновляет claims после успешной OAuth авторизации
	// Вызывается go-pkgz/auth для обработки OAuth токенов
	ClaimsUpdater(claims token.Claims) token.Claims

	// GetOrCreateUserFromOAuth получает или создает пользователя на основе OAuth данных
	GetOrCreateUserFromOAuth(oauthUser token.User) (*models.User, error)

	// GetUserFromToken получает пользователя из token.User
	GetUserFromToken(tokenUser token.User) (*models.User, error)

	// ============================================================================
	// Аутентификация через GitHub
	// ============================================================================

	// AuthenticateWithGitHub аутентифицирует пользователя через GitHub
	AuthenticateWithGitHub(githubID string) (*models.User, error)

	// LinkGitHubAccount связывает GitHub аккаунт с существующим пользователем
	LinkGitHubAccount(userID int64, githubID string, githubData map[string]interface{}) error

	// UnlinkGitHubAccount отвязывает GitHub аккаунт от пользователя
	UnlinkGitHubAccount(userID int64) error

	// ============================================================================
	// Аутентификация через Google
	// ============================================================================

	// AuthenticateWithGoogle аутентифицирует пользователя через Google
	AuthenticateWithGoogle(googleID string) (*models.User, error)

	// LinkGoogleAccount связывает Google аккаунт с существующим пользователем
	LinkGoogleAccount(userID int64, googleID string, googleData map[string]interface{}) error

	// UnlinkGoogleAccount отвязывает Google аккаунт от пользователя
	UnlinkGoogleAccount(userID int64) error

	// ============================================================================
	// Аутентификация через Telegram
	// ============================================================================

	// AuthenticateWithTelegram аутентифицирует пользователя через Telegram
	AuthenticateWithTelegram(telegramID string) (*models.User, error)

	// LinkTelegramAccount связывает Telegram аккаунт с существующим пользователем
	LinkTelegramAccount(userID int64, telegramID string, telegramData map[string]interface{}) error

	// UnlinkTelegramAccount отвязывает Telegram аккаунт от пользователя
	UnlinkTelegramAccount(userID int64) error
}
