package repositories

import (
	"context"
	"errors"
	"fmt"
	"it_rabotyagi/internal/business/models"
	"it_rabotyagi/internal/data/database"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CheckUser(ctx context.Context, email, username string) (*int, error) {
	query := "SELECT id FROM users WHERE email=$1 OR username=$2"
	var userID int
	err := u.db.Pool.QueryRow(ctx, query, email, username).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Пользователь не найден - это нормально
			return nil, nil
		}
		return nil, err
	}
	return &userID, nil
}

func (u *UserRepository) CreateUser(ctx context.Context, email, username, password string) (int, error) {
	query := "INSERT INTO users (email, username, password, role) VALUES ($1, $2, $3, 'user') RETURNING id"
	var userID int
	err := u.db.Pool.QueryRow(ctx, query, email, username, password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// GetUserByID получает пользователя по ID
func (u *UserRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `SELECT id, username, password, email, telegram_id, google_id, github_id, 
              name, avatar_url, description, role, created_at, updated_at 
              FROM users WHERE id = $1`

	user := &models.User{}
	err := u.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.TelegramID,
		&user.GoogleID,
		&user.GithubID,
		&user.Name,
		&user.AvatarURL,
		&user.Description,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail получает пользователя по email
func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, password, email, telegram_id, google_id, github_id,
              name, avatar_url, description, role, created_at, updated_at
              FROM users WHERE email = $1`

	user := &models.User{}
	err := u.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.TelegramID,
		&user.GoogleID,
		&user.GithubID,
		&user.Name,
		&user.AvatarURL,
		&user.Description,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserProfile обновляет основные данные профиля пользователя
func (u *UserRepository) UpdateUserProfile(ctx context.Context, userID int, username, name, email, description *string) error {
	query := "UPDATE users SET updated_at = NOW()"
	args := []interface{}{userID}
	argPos := 2

	if username != nil {
		query += fmt.Sprintf(", username = $%d", argPos)
		args = append(args, *username)
		argPos++
	}
	if name != nil {
		query += fmt.Sprintf(", name = $%d", argPos)
		args = append(args, *name)
		argPos++
	}
	if email != nil {
		query += fmt.Sprintf(", email = $%d", argPos)
		args = append(args, *email)
		argPos++
	}
	if description != nil {
		query += fmt.Sprintf(", description = $%d", argPos)
		args = append(args, *description)
		argPos++
	}

	query += " WHERE id = $1"

	_, err := u.db.Pool.Exec(ctx, query, args...)
	return err
}

// UpdateUserAvatar обновляет avatar_url пользователя
func (u *UserRepository) UpdateUserAvatar(ctx context.Context, userID int, avatarURL string) error {
	query := "UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2"
	_, err := u.db.Pool.Exec(ctx, query, avatarURL, userID)
	return err
}

// UpdateUserPassword обновляет хеш пароля пользователя
func (u *UserRepository) UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) error {
	query := "UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2"
	_, err := u.db.Pool.Exec(ctx, query, hashedPassword, userID)
	return err
}

// CheckUsernameExists проверяет существование username (исключая текущего пользователя)
func (u *UserRepository) CheckUsernameExists(ctx context.Context, username string, excludeUserID int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id != $2)"
	var exists bool
	err := u.db.Pool.QueryRow(ctx, query, username, excludeUserID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CheckEmailExists проверяет существование email (исключая текущего пользователя)
func (u *UserRepository) CheckEmailExists(ctx context.Context, email string, excludeUserID int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)"
	var exists bool
	err := u.db.Pool.QueryRow(ctx, query, email, excludeUserID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
