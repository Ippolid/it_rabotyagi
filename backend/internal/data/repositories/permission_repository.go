package repositories

import (
	"context"

	"it_rabotyagi/internal/data/database"
)

// PermissionRepository хранит информацию о пользователях с правами на изменение сущностей.
type PermissionRepository struct {
	db *database.DB
}

func NewPermissionRepository(db *database.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// IsEditor проверяет, есть ли у пользователя право редактировать курсы/модули/вопросы.
func (r *PermissionRepository) IsEditor(ctx context.Context, userID int) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM allowed_editors WHERE user_id = $1)`, userID).Scan(&exists)
	return exists, err
}
