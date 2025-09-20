package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"itpath/internal/data"
	"itpath/internal/data/entities"
)

type userRepository struct {
	db *pgxpool.Pool
}

func (r *userRepository) Update(ctx context.Context, id int64, req entities.UpdateUserRequest) (*entities.UserEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) UpdateRole(ctx context.Context, userID int64, role entities.UserRole) (*entities.UserEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) UpdateSubscription(ctx context.Context, userID int64, subscription entities.SubscriptionType) (*entities.UserEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.UserEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entities.UserEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) UsernameExists(ctx context.Context, username string, excludeUserID ...int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) CountTotal(ctx context.Context) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository(db *pgxpool.Pool) data.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*entities.UserEntity, error) {
	query := `
        SELECT id, telegram_id, username, first_name, last_name, photo_url, 
               email, github_id, google_id, role, subscription, created_at, updated_at
        FROM users 
        WHERE id = $1
    `

	var user entities.UserEntity
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName,
		&user.LastName, &user.PhotoURL, &user.Email, &user.GithubID,
		&user.GoogleID, &user.Role, &user.Subscription, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*entities.UserEntity, error) {
	query := `
        SELECT id, telegram_id, username, first_name, last_name, photo_url, 
               email, github_id, google_id, role, subscription, created_at, updated_at
        FROM users 
        WHERE telegram_id = $1
    `

	var user entities.UserEntity
	err := r.db.QueryRow(ctx, query, telegramID).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName,
		&user.LastName, &user.PhotoURL, &user.Email, &user.GithubID,
		&user.GoogleID, &user.Role, &user.Subscription, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by telegram_id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, req entities.CreateUserRequest) (*entities.UserEntity, error) {
	query := `
        INSERT INTO users (telegram_id, username, first_name, last_name, photo_url, role, subscription)
        VALUES ($1, NULLIF($2, ''), $3, NULLIF($4, ''), NULLIF($5, ''), $6, $7)
        RETURNING id, telegram_id, username, first_name, last_name, photo_url, 
                  email, github_id, google_id, role, subscription, created_at, updated_at
    `

	var user entities.UserEntity
	err := r.db.QueryRow(ctx, query,
		req.TelegramID,
		req.Username,
		req.FirstName,
		req.LastName,
		req.PhotoURL,
		req.Role,
		req.Subscription,
	).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName,
		&user.LastName, &user.PhotoURL, &user.Email, &user.GithubID,
		&user.GoogleID, &user.Role, &user.Subscription, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) UpdateTelegramData(ctx context.Context, telegramID int64, username, firstName, lastName, photoURL string) (*entities.UserEntity, error) {
	query := `
        UPDATE users 
        SET username = NULLIF($2, ''), 
            first_name = $3, 
            last_name = NULLIF($4, ''), 
            photo_url = NULLIF($5, ''), 
            updated_at = NOW()
        WHERE telegram_id = $1
        RETURNING id, telegram_id, username, first_name, last_name, photo_url, 
                  email, github_id, google_id, role, subscription, created_at, updated_at
    `

	var user entities.UserEntity
	err := r.db.QueryRow(ctx, query, telegramID, username, firstName, lastName, photoURL).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName,
		&user.LastName, &user.PhotoURL, &user.Email, &user.GithubID,
		&user.GoogleID, &user.Role, &user.Subscription, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update telegram data: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByRole(ctx context.Context, role entities.UserRole, limit, offset int) ([]*entities.UserEntity, error) {
	query := `
        SELECT id, telegram_id, username, first_name, last_name, photo_url, 
               email, github_id, google_id, role, subscription, created_at, updated_at
        FROM users 
        WHERE role = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(ctx, query, role, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	defer rows.Close()

	var users []*entities.UserEntity
	for rows.Next() {
		var user entities.UserEntity
		err := rows.Scan(
			&user.ID, &user.TelegramID, &user.Username, &user.FirstName,
			&user.LastName, &user.PhotoURL, &user.Email, &user.GithubID,
			&user.GoogleID, &user.Role, &user.Subscription, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string, excludeUserID ...int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1`
	args := []interface{}{email}

	if len(excludeUserID) > 0 {
		query += ` AND id != $2`
		args = append(args, excludeUserID[0])
	}
	query += `)`

	var exists bool
	err := r.db.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

func (r *userRepository) CountByRole(ctx context.Context, role entities.UserRole) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE role = $1`

	var count int64
	err := r.db.QueryRow(ctx, query, role).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}

	return count, nil
}

// Реализация остальных методов...
