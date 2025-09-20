package data

import (
	"context"
	"itpath/internal/data/entities"
)

// UserRepository определяет контракт работы с данными пользователей
type UserRepository interface {
	// Основные CRUD операции
	GetByID(ctx context.Context, id int64) (*entities.UserEntity, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*entities.UserEntity, error)
	Create(ctx context.Context, req entities.CreateUserRequest) (*entities.UserEntity, error)
	Update(ctx context.Context, id int64, req entities.UpdateUserRequest) (*entities.UserEntity, error)
	Delete(ctx context.Context, id int64) error

	// Специализированные операции
	UpdateTelegramData(ctx context.Context, telegramID int64, username, firstName, lastName, photoURL string) (*entities.UserEntity, error)
	UpdateRole(ctx context.Context, userID int64, role entities.UserRole) (*entities.UserEntity, error)
	UpdateSubscription(ctx context.Context, userID int64, subscription entities.SubscriptionType) (*entities.UserEntity, error)

	// Поиск и фильтрация
	GetByRole(ctx context.Context, role entities.UserRole, limit, offset int) ([]*entities.UserEntity, error)

	// Проверки существования
	EmailExists(ctx context.Context, email string, excludeUserID ...int64) (bool, error)
	UsernameExists(ctx context.Context, username string, excludeUserID ...int64) (bool, error)
}
