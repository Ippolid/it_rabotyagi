package models

import "time"

// User представляет модель пользователя
type User struct {
	ID          int
	Username    string
	Password    string
	Email       *string
	TelegramID  *string
	GoogleID    *string
	GithubID    *string
	Name        *string
	AvatarURL   *string
	Description *string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserProfileUpdate представляет DTO для обновления профиля пользователя
type UserProfileUpdate struct {
	Username    *string
	Name        *string
	Email       *string
	Description *string
}

// PasswordChange представляет DTO для смены пароля
type PasswordChange struct {
	OldPassword string
	NewPassword string
}
