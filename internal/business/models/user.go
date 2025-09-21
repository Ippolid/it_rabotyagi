package models

import (
	"itpath/internal/data/entities"
	"time"
)

// User - бизнес-модель пользователя
type User struct {
	ID           int64                     `json:"id"`
	TelegramID   int64                     `json:"telegram_id"`
	Username     *string                   `json:"username,omitempty"`
	FirstName    string                    `json:"first_name"`
	LastName     *string                   `json:"last_name,omitempty"`
	PhotoURL     *string                   `json:"photo_url,omitempty"`
	Email        *string                   `json:"email,omitempty"`
	Role         entities.UserRole         `json:"role"`
	Subscription entities.SubscriptionType `json:"subscription"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// PublicUser - публичная модель пользователя (без чувствительных данных)
type PublicUser struct {
	ID           int64                     `json:"id"`
	Username     *string                   `json:"username,omitempty"`
	FirstName    string                    `json:"first_name"`
	LastName     *string                   `json:"last_name,omitempty"`
	PhotoURL     *string                   `json:"photo_url,omitempty"`
	Role         entities.UserRole         `json:"role"`
	Subscription entities.SubscriptionType `json:"subscription"`
	CreatedAt    time.Time                 `json:"created_at"`
}

// AuthResult - результат авторизации
type AuthResult struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *PublicUser `json:"user"`
	ExpiresIn    int64       `json:"expires_in"`
}

// TelegramAuthData - данные авторизации от Telegram
// Поля соответствуют данным, которые возвращает виджет Telegram Login.
type TelegramAuthData struct {
	ID        int64  `json:"id"`
	AuthDate  int64  `json:"auth_date"`
	FirstName string `json:"first_name"`
	Hash      string `json:"hash"`
	// Необязательные поля
	LastName string `json:"last_name,omitempty"`
	Username string `json:"username,omitempty"`
	PhotoURL string `json:"photo_url,omitempty"`
}

// UpdateProfileRequest - запрос на обновление профиля
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Username  *string `json:"username,omitempty"`
	PhotoURL  *string `json:"photo_url,omitempty"`
}

// Конвертация из Entity в Business Model
func FromEntity(entity *entities.UserEntity) *User {
	if entity == nil {
		return nil
	}

	return &User{
		ID:           entity.ID,
		TelegramID:   entity.TelegramID,
		Username:     entity.Username,
		FirstName:    entity.FirstName,
		LastName:     entity.LastName,
		PhotoURL:     entity.PhotoURL,
		Role:         entity.Role,
		Subscription: entity.Subscription,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

// ToPublic - конвертация в публичную модель
func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:           u.ID,
		Username:     u.Username,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		PhotoURL:     u.PhotoURL,
		Role:         u.Role,
		Subscription: u.Subscription,
		CreatedAt:    u.CreatedAt,
	}
}
