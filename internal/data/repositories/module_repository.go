package repositories

import (
    "context"
    "it_rabotyagi/internal/business/models"
    "it_rabotyagi/internal/data/database"
)

type ModuleRepository struct {
    db *database.DB
}

func NewModuleRepository(db *database.DB) *ModuleRepository {
    return &ModuleRepository{db: db}
}

func (r *ModuleRepository) ListByCourse(ctx context.Context, courseID int64) ([]*models.Module, error) {
    rows, err := r.db.Pool.Query(ctx, `
        SELECT id, course_id, title, description, module_order, created_at, edited_at
        FROM modules
        WHERE course_id=$1
        ORDER BY module_order ASC, id ASC
    `, courseID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []*models.Module
    for rows.Next() {
        m := &models.Module{}
        if err := rows.Scan(&m.ID, &m.CourseID, &m.Title, &m.Description, &m.ModuleOrder, &m.CreatedAt, &m.EditedAt); err != nil {
            return nil, err
        }
        out = append(out, m)
    }
    return out, nil
}

func (r *ModuleRepository) GetByID(ctx context.Context, id int64) (*models.Module, error) {
    row := r.db.Pool.QueryRow(ctx, `
        SELECT id, course_id, title, description, module_order, created_at, edited_at
        FROM modules
        WHERE id = $1
    `, id)
    m := &models.Module{}
    if err := row.Scan(&m.ID, &m.CourseID, &m.Title, &m.Description, &m.ModuleOrder, &m.CreatedAt, &m.EditedAt); err != nil {
        return nil, err
    }
    return m, nil
}
