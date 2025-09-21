package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"itpath/internal/data"
	"itpath/internal/data/database"
	"itpath/internal/data/entities"
	"strings"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *database.DB) data.UserRepository {
	return &userRepository{
		db: db.Pool,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, req entities.CreateUserRequest) (*entities.UserEntity, error) {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, photo_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, telegram_id, username, first_name, last_name, photo_url, role, subscription, created_at, updated_at
	`

	user := &entities.UserEntity{}

	err := r.db.QueryRow(ctx, query,
		req.TelegramID,
		req.Username,
		req.FirstName,
		req.LastName,
		req.PhotoURL,
	).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PhotoURL,
		&user.Role,
		&user.Subscription,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// ИЗМЕНЕНО: теперь обновляет по telegram_id, а не по id
func (r *userRepository) Update(ctx context.Context, telegramId int64, req entities.UpdateUserRequest) (*entities.UserEntity, error) {
	var setClauses []string
	var args []interface{}
	argCount := 1

	if req.Username != nil {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", argCount))
		args = append(args, *req.Username)
		argCount++
	}
	if req.FirstName != nil {
		setClauses = append(setClauses, fmt.Sprintf("first_name = $%d", argCount))
		args = append(args, *req.FirstName)
		argCount++
	}
	if req.LastName != nil {
		setClauses = append(setClauses, fmt.Sprintf("last_name = $%d", argCount))
		args = append(args, *req.LastName)
		argCount++
	}
	if req.PhotoURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("photo_url = $%d", argCount))
		args = append(args, *req.PhotoURL)
		argCount++
	}

	if len(setClauses) == 0 {
		setClauses = append(setClauses, "updated_at = NOW()")
	} else {
		setClauses = append(setClauses, "updated_at = NOW()")
	}

	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE telegram_id = $%d
		RETURNING id, telegram_id, username, first_name, last_name, photo_url, role, subscription, created_at, updated_at
	`, strings.Join(setClauses, ", "), argCount)

	args = append(args, telegramId)

	user := &entities.UserEntity{}

	err := r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PhotoURL,
		&user.Role,
		&user.Subscription,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*entities.UserEntity, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, photo_url, role, subscription, created_at, updated_at
		FROM users
		WHERE telegram_id = $1
	`

	user := &entities.UserEntity{}

	err := r.db.QueryRow(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PhotoURL,
		&user.Role,
		&user.Subscription,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by telegram id from db: %w", err)
	}

	return user, nil
}
