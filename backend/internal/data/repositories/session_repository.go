package repositories

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"it_rabotyagi/internal/data/database"
	"time"
)

type SessionRepository struct {
	db *database.DB
}

func NewSessionRepository(db *database.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Session представляет сессию пользователя
type Session struct {
	ID               int
	UserID           int
	RefreshTokenHash string
	ExpiresAt        time.Time
	CreatedAt        time.Time
	RevokedAt        *time.Time
}

// HashToken хеширует токен для безопасного хранения
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// CreateSession создает новую сессию
func (r *SessionRepository) CreateSession(ctx context.Context, userID int, refreshToken string, expiresAt time.Time) error {
	tokenHash := HashToken(refreshToken)
	query := `INSERT INTO auth_sessions (user_id, refresh_token_hash, expires_at) 
              VALUES ($1, $2, $3)`

	_, err := r.db.Pool.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}

// GetSessionByToken получает сессию по refresh токену
func (r *SessionRepository) GetSessionByToken(ctx context.Context, refreshToken string) (*Session, error) {
	tokenHash := HashToken(refreshToken)
	query := `SELECT id, user_id, refresh_token_hash, expires_at, created_at, revoked_at 
              FROM auth_sessions 
              WHERE refresh_token_hash = $1 AND revoked_at IS NULL`

	session := &Session{}
	err := r.db.Pool.QueryRow(ctx, query, tokenHash).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshTokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.RevokedAt,
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// RevokeSession отзывает сессию (logout)
func (r *SessionRepository) RevokeSession(ctx context.Context, refreshToken string) error {
	tokenHash := HashToken(refreshToken)
	query := `UPDATE auth_sessions 
              SET revoked_at = now() 
              WHERE refresh_token_hash = $1 AND revoked_at IS NULL`

	_, err := r.db.Pool.Exec(ctx, query, tokenHash)
	return err
}

// RevokeAllUserSessions отзывает все сессии пользователя
func (r *SessionRepository) RevokeAllUserSessions(ctx context.Context, userID int) error {
	query := `UPDATE auth_sessions 
              SET revoked_at = now() 
              WHERE user_id = $1 AND revoked_at IS NULL`

	_, err := r.db.Pool.Exec(ctx, query, userID)
	return err
}

// CleanExpiredSessions удаляет истекшие и отозванные сессии
func (r *SessionRepository) CleanExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM auth_sessions 
              WHERE expires_at < now() OR revoked_at IS NOT NULL`

	_, err := r.db.Pool.Exec(ctx, query)
	return err
}

// GetUserActiveSessions получает список активных сессий пользователя
func (r *SessionRepository) GetUserActiveSessions(ctx context.Context, userID int) ([]*Session, error) {
	query := `SELECT id, user_id, refresh_token_hash, expires_at, created_at, revoked_at 
              FROM auth_sessions 
              WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > now()
              ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		session := &Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.RefreshTokenHash,
			&session.ExpiresAt,
			&session.CreatedAt,
			&session.RevokedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
