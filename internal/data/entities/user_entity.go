package entities

import (
	"database/sql/driver"
	"time"
)

// Enum типы для базы данных
type UserRole string
type SubscriptionType string

const (
	RoleUser   UserRole = "user"
	RoleMentor UserRole = "mentor"
)

const (
	SubscriptionTrial SubscriptionType = "trial"
	SubscriptionPro   SubscriptionType = "pro"
)

// Реализация driver.Valuer для ENUM'ов
func (ur UserRole) Value() (driver.Value, error) {
	return string(ur), nil
}

func (st SubscriptionType) Value() (driver.Value, error) {
	return string(st), nil
}

// UserEntity - сущность для базы данных
type UserEntity struct {
	ID           int64            `db:"id"`
	TelegramID   int64            `db:"telegram_id"`
	Username     *string          `db:"username"`
	FirstName    string           `db:"first_name"`
	LastName     *string          `db:"last_name"`
	PhotoURL     *string          `db:"photo_url"`
	Email        *string          `db:"email"`
	Role         UserRole         `db:"role"`
	Subscription SubscriptionType `db:"subscription"`
	CreatedAt    time.Time        `db:"created_at"`
	UpdatedAt    time.Time        `db:"updated_at"`
}

// CreateUserRequest - данные для создания пользователя
type CreateUserRequest struct {
	TelegramID   int64
	Username     string
	FirstName    string
	LastName     string
	PhotoURL     string
	Role         UserRole
	Subscription SubscriptionType
}

// UpdateUserRequest - данные для обновления пользователя
type UpdateUserRequest struct {
	Username  *string
	FirstName *string
	LastName  *string
	PhotoURL  *string
	Email     *string
}
