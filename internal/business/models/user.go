package models

import "time"

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
