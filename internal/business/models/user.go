package models

import (
	"itpath/internal/data/entities"
	"time"
)

type User struct {
	ID                    int64
	TelegramID            *string
	GoogleID              *string
	GitHubID              *string
	Email                 *string
	Username              *string
	Name                  string
	AvatarURL             *string
	Description           *string
	Role                  string
	SubscriptionType      *string
	SubscriptionExpiresAt *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// ConvertToModel конвертирует entity в бизнес-модель
func ConvertToModel(userEntity *entities.UserEntity) *User {
	user := &User{
		ID:                    userEntity.ID,
		Name:                  userEntity.Name,
		Role:                  string(userEntity.Role),
		TelegramID:            userEntity.TelegramID,
		GoogleID:              userEntity.GoogleID,
		GitHubID:              userEntity.GitHubID,
		Email:                 userEntity.Email,
		Username:              userEntity.Username,
		AvatarURL:             userEntity.AvatarURL,
		Description:           userEntity.Description,
		SubscriptionExpiresAt: userEntity.SubscriptionExpiresAt,
		CreatedAt:             userEntity.CreatedAt,
		UpdatedAt:             userEntity.UpdatedAt,
	}

	if userEntity.SubscriptionType != nil {
		subType := string(*userEntity.SubscriptionType)
		user.SubscriptionType = &subType
	}

	return user
}

// ConvertToEntity конвертирует бизнес-модель в entity
func ConvertToEntity(user *User) *entities.UserEntity {
	userEntity := &entities.UserEntity{
		ID:                    user.ID,
		TelegramID:            user.TelegramID,
		GoogleID:              user.GoogleID,
		GitHubID:              user.GitHubID,
		Email:                 user.Email,
		Username:              user.Username,
		Name:                  user.Name,
		AvatarURL:             user.AvatarURL,
		Description:           user.Description,
		Role:                  entities.UserRole(user.Role),
		SubscriptionExpiresAt: user.SubscriptionExpiresAt,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
	}

	if user.SubscriptionType != nil {
		subType := entities.SubscriptionType(*user.SubscriptionType)
		userEntity.SubscriptionType = &subType
	}

	return userEntity
}
