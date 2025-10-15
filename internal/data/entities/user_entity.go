package entities

import "time"

type UserRole string

const (
	RoleUser   UserRole = "user"
	RoleMentor UserRole = "mentor"
	RoleAdmin  UserRole = "admin"
)

type SubscriptionType string

const (
	SubscriptionFree       SubscriptionType = "free"
	SubscriptionPro        SubscriptionType = "pro"
	SubscriptionTeam       SubscriptionType = "team"
	SubscriptionEnterprise SubscriptionType = "enterprise"
)

type UserEntity struct {
	ID                    int64             `db:"id"`
	TelegramID            *string           `db:"telegram_id"`
	GoogleID              *string           `db:"google_id"`
	GitHubID              *string           `db:"github_id"`
	Email                 *string           `db:"email"`
	Username              *string           `db:"username"`
	Name                  string            `db:"name"`
	AvatarURL             *string           `db:"avatar_url"`
	Description           *string           `db:"description"`
	Role                  UserRole          `db:"role"`
	SubscriptionType      *SubscriptionType `db:"subscription_type"`
	SubscriptionExpiresAt *time.Time        `db:"subscription_expires_at"`
	CreatedAt             time.Time         `db:"created_at"`
	UpdatedAt             time.Time         `db:"updated_at"`
}
