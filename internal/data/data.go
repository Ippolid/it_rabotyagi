package data

import "context"

// UserRepository определяет контракт работы с данными пользователей
type UserRepository interface {
	CheckUser(ctx context.Context, email, nickname string) (*int, error)
	CreateUser(ctx context.Context, email, nickname, passwordHash string) (int, error)
}
