package repositories

import (
	"context"
	"errors"
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

func (u *UserRepository) CheckUser(ctx context.Context, email, username string) (*int64, error) {
	query := "SELECT id FROM users WHERE email=$1 OR username=$2"
	var userID int64
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
