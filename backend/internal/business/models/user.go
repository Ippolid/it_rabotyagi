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
