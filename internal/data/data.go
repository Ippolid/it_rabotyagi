package data

import (
	"itpath/internal/data/entities"
)

// UserRepository определяет контракт работы с данными пользователей
type UserRepository interface {
	FindUserByTelegramID(telegramID string) (*entities.UserEntity, error)
	FindUserByGoogleID(googleID string) (*entities.UserEntity, error)
	FindUserByGitHubID(githubID string) (*entities.UserEntity, error)
	FindUserByEmail(email string) (*entities.UserEntity, error)
	FindUserByID(id int64) (*entities.UserEntity, error)
	CreateUser(user *entities.UserEntity) error
	UpdateUser(user *entities.UserEntity) error
}
