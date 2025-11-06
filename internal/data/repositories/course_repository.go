package repositories

import (
    "context"
    "it_rabotyagi/internal/business/models"
    "it_rabotyagi/internal/data/database"
)

type CourseRepository struct {
    db *database.DB
}

func NewCourseRepository(db *database.DB) *CourseRepository {
    return &CourseRepository{db: db}
}

func (r *CourseRepository) ListPublished(ctx context.Context, limit, offset int) ([]*models.Course, error) {
    if limit <= 0 { limit = 12 }
    if offset < 0 { offset = 0 }
    rows, err := r.db.Pool.Query(ctx, `
        SELECT id, title, description, is_published, created_at, updated_at
        FROM courses
        WHERE is_published = TRUE
        ORDER BY id ASC
        LIMIT $1 OFFSET $2
    `, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []*models.Course
    for rows.Next() {
        c := &models.Course{}
        if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.IsPublished, &c.CreatedAt, &c.UpdatedAt); err != nil {
            return nil, err
        }
        out = append(out, c)
    }
    return out, nil
}


func (r *CourseRepository) GetByID(ctx context.Context, id int64) (*models.Course, error) {
    row := r.db.Pool.QueryRow(ctx, `
        SELECT id, title, description, is_published, created_at, updated_at
        FROM courses
        WHERE id = $1
    `, id)

    c := &models.Course{}
    if err := row.Scan(&c.ID, &c.Title, &c.Description, &c.IsPublished, &c.CreatedAt, &c.UpdatedAt); err != nil {
        return nil, err
    }
    return c, nil
}


