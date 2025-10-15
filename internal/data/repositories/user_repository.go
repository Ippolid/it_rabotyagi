package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"itpath/internal/data/database"
	"itpath/internal/data/entities"
	"itpath/internal/logger"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindUserByTelegramID находит пользователя по Telegram ID
func (r *UserRepository) FindUserByTelegramID(telegramID string) (*entities.UserEntity, error) {
	query := `SELECT id, telegram_id, google_id, github_id, email, username, name, avatar_url, description, role, subscription_type, subscription_expires_at, created_at, updated_at 
	          FROM users WHERE telegram_id = $1`

	user := &entities.UserEntity{}
	err := r.db.Pool.QueryRow(context.Background(), query, telegramID).Scan(
		&user.ID, &user.TelegramID, &user.GoogleID, &user.GitHubID, &user.Email, &user.Username,
		&user.Name, &user.AvatarURL, &user.Description, &user.Role, &user.SubscriptionType,
		&user.SubscriptionExpiresAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by telegram_id: %w", err)
	}

	logger.Debug("DB: FindByITelegramD",
		zap.String("telegramid", telegramID),
		zap.Any("user", user),
		zap.Error(err),
	)

	return user, nil
}

// FindUserByGoogleID находит пользователя по Google ID
func (r *UserRepository) FindUserByGoogleID(googleID string) (*entities.UserEntity, error) {
	query := `SELECT id, telegram_id, google_id, github_id, email, username, name, avatar_url, description, role, subscription_type, subscription_expires_at, created_at, updated_at 
	          FROM users WHERE google_id = $1`

	user := &entities.UserEntity{}
	err := r.db.Pool.QueryRow(context.Background(), query, googleID).Scan(
		&user.ID, &user.TelegramID, &user.GoogleID, &user.GitHubID, &user.Email, &user.Username,
		&user.Name, &user.AvatarURL, &user.Description, &user.Role, &user.SubscriptionType,
		&user.SubscriptionExpiresAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by google_id: %w", err)
	}

	logger.Debug("DB: FindByGoogleID",
		zap.String("GoogleId", googleID),
		zap.Any("user", user),
		zap.Error(err),
	)

	return user, nil
}

// FindUserByGitHubID находит пользователя по GitHub ID
func (r *UserRepository) FindUserByGitHubID(githubID string) (*entities.UserEntity, error) {
	query := `SELECT id, telegram_id, google_id, github_id, email, username, name, avatar_url, description, role, subscription_type, subscription_expires_at, created_at, updated_at 
	          FROM users WHERE github_id = $1`

	user := &entities.UserEntity{}
	err := r.db.Pool.QueryRow(context.Background(), query, githubID).Scan(
		&user.ID, &user.TelegramID, &user.GoogleID, &user.GitHubID, &user.Email, &user.Username,
		&user.Name, &user.AvatarURL, &user.Description, &user.Role, &user.SubscriptionType,
		&user.SubscriptionExpiresAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by github_id: %w", err)
	}

	logger.Debug("DB: FindByGitHubID",
		zap.String("githubId", githubID),
		zap.Any("user", user),
		zap.Error(err),
	)

	return user, nil
}

// FindUserByEmail находит пользователя по email
func (r *UserRepository) FindUserByEmail(email string) (*entities.UserEntity, error) {
	query := `SELECT id, telegram_id, google_id, github_id, email, username, name, avatar_url, description, role, subscription_type, subscription_expires_at, created_at, updated_at 
	          FROM users WHERE email = $1`

	user := &entities.UserEntity{}
	err := r.db.Pool.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.TelegramID, &user.GoogleID, &user.GitHubID, &user.Email, &user.Username,
		&user.Name, &user.AvatarURL, &user.Description, &user.Role, &user.SubscriptionType,
		&user.SubscriptionExpiresAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	logger.Debug("DB: FindByID",
		zap.String("email", email),
		zap.Any("user", user),
		zap.Error(err),
	)

	return user, nil
}

// FindUserByID находит пользователя по ID
func (r *UserRepository) FindUserByID(id int64) (*entities.UserEntity, error) {
	query := `SELECT id, telegram_id, google_id, github_id, email, username, name, avatar_url, description, role, subscription_type, subscription_expires_at, created_at, updated_at 
	          FROM users WHERE id = $1`

	user := &entities.UserEntity{}
	err := r.db.Pool.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.TelegramID, &user.GoogleID, &user.GitHubID, &user.Email, &user.Username,
		&user.Name, &user.AvatarURL, &user.Description, &user.Role, &user.SubscriptionType,
		&user.SubscriptionExpiresAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	logger.Debug("DB: FindByID",
		zap.Int64("id", id),
		zap.Any("user", user),
		zap.Error(err),
	)

	return user, nil
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(user *entities.UserEntity) error {
	query := `INSERT INTO users (telegram_id, google_id, github_id, email, username, name, avatar_url, description, role) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	          RETURNING id, created_at, updated_at`

	err := r.db.Pool.QueryRow(
		context.Background(),
		query,
		user.TelegramID, user.GoogleID, user.GitHubID, user.Email, user.Username,
		user.Name, user.AvatarURL, user.Description, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	logger.Debug("DB: Create user", zap.Int64("user", user.ID))

	return nil
}

// UpdateUser обновляет данные пользователя
func (r *UserRepository) UpdateUser(user *entities.UserEntity) error {
	query := `UPDATE users 
	          SET telegram_id = $1, google_id = $2, github_id = $3, email = $4, username = $5, 
	              name = $6, avatar_url = $7, description = $8, role = $9
	          WHERE id = $10`

	_, err := r.db.Pool.Exec(
		context.Background(),
		query,
		user.TelegramID, user.GoogleID, user.GitHubID, user.Email, user.Username,
		user.Name, user.AvatarURL, user.Description, user.Role, user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	logger.Debug("DB: Update user", zap.Int64("user", user.ID))

	return nil
}
