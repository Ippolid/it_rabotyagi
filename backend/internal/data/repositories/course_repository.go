package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"it_rabotyagi/internal/data/database"
)

// CourseRepository инкапсулирует работу с курсами, модулями и прогрессом.
type CourseRepository struct {
	db *database.DB
}

func NewCourseRepository(db *database.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

type Course struct {
	ID          int
	Title       string
	Description string
	IsPublished bool
}

type Module struct {
	ID          int
	CourseID    int
	Title       string
	Description string
	Order       int
	Content     *string
}

type ModuleProgress struct {
	Module      Module
	Completed   bool
	CompletedAt *time.Time
}

func (r *CourseRepository) ListCourses(ctx context.Context, limit, offset int) ([]Course, int, error) {
	var total int
	if err := r.db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM courses").Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, title, description, is_published FROM courses ORDER BY id`
	args := []any{}
	if limit > 0 {
		query += " LIMIT $1 OFFSET $2"
		args = append(args, limit, offset)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.IsPublished); err != nil {
			return nil, 0, err
		}
		courses = append(courses, c)
	}

	return courses, total, nil
}

func (r *CourseRepository) GetCourse(ctx context.Context, id int) (*Course, error) {
	query := `SELECT id, title, description, is_published FROM courses WHERE id = $1`
	var c Course
	if err := r.db.Pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.Title, &c.Description, &c.IsPublished); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CourseRepository) CreateCourse(ctx context.Context, title, description string, isPublished bool) (int, error) {
	query := `INSERT INTO courses (title, description, is_published) VALUES ($1, $2, $3) RETURNING id`
	var id int
	if err := r.db.Pool.QueryRow(ctx, query, title, description, isPublished).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CourseRepository) UpdateCourse(ctx context.Context, id int, title, description *string, isPublished *bool) error {
	setParts := []string{}
	args := []any{}
	idx := 1

	if title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", idx))
		args = append(args, *title)
		idx++
	}
	if description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", idx))
		args = append(args, *description)
		idx++
	}
	if isPublished != nil {
		setParts = append(setParts, fmt.Sprintf("is_published = $%d", idx))
		args = append(args, *isPublished)
		idx++
	}

	if len(setParts) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE courses SET %s, updated_at = now() WHERE id = $%d", strings.Join(setParts, ", "), idx)
	args = append(args, id)

	cmd, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CourseRepository) DeleteCourse(ctx context.Context, id int) error {
	cmd, err := r.db.Pool.Exec(ctx, "DELETE FROM courses WHERE id = $1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CourseRepository) ListModules(ctx context.Context, courseID int) ([]Module, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT id, course_id, title, description, module_order FROM modules WHERE course_id = $1 ORDER BY module_order`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []Module
	for rows.Next() {
		var m Module
		if err := rows.Scan(&m.ID, &m.CourseID, &m.Title, &m.Description, &m.Order); err != nil {
			return nil, err
		}
		m.Content = nil
		modules = append(modules, m)
	}
	return modules, nil
}

func (r *CourseRepository) GetModule(ctx context.Context, courseID, moduleID int) (*Module, error) {
	var m Module
	err := r.db.Pool.QueryRow(ctx, `SELECT id, course_id, title, description, module_order FROM modules WHERE course_id = $1 AND id = $2`, courseID, moduleID).
		Scan(&m.ID, &m.CourseID, &m.Title, &m.Description, &m.Order)
	if err != nil {
		return nil, err
	}
	m.Content = nil
	return &m, nil
}

func (r *CourseRepository) CreateModule(ctx context.Context, courseID int, title, description string, order int, content *string) (int, error) {
	var id int
	err := r.db.Pool.QueryRow(ctx, `INSERT INTO modules (course_id, title, description, module_order) VALUES ($1, $2, $3, $4) RETURNING id`, courseID, title, description, order).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CourseRepository) UpdateModule(ctx context.Context, courseID, moduleID int, title, description, content *string, order *int) error {
	setParts := []string{}
	args := []any{}
	idx := 1
	if title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", idx))
		args = append(args, *title)
		idx++
	}
	if description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", idx))
		args = append(args, *description)
		idx++
	}
	if order != nil {
		setParts = append(setParts, fmt.Sprintf("module_order = $%d", idx))
		args = append(args, *order)
		idx++
	}

	if len(setParts) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE modules SET %s, edited_at = now() WHERE course_id = $%d AND id = $%d", strings.Join(setParts, ", "), idx, idx+1)
	args = append(args, courseID, moduleID)

	cmd, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CourseRepository) DeleteModule(ctx context.Context, courseID, moduleID int) error {
	cmd, err := r.db.Pool.Exec(ctx, "DELETE FROM modules WHERE course_id = $1 AND id = $2", courseID, moduleID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CourseRepository) EnrollUser(ctx context.Context, userID, courseID, totalModules int) (time.Time, error) {
	// Insert or retain existing enrollment
	query := `
INSERT INTO user_course_progress (user_id, course_id, total_modules, completed_modules, module_progress_pct)
VALUES ($1, $2, $3, 0, 0)
ON CONFLICT (user_id, course_id) DO UPDATE SET updated_at = now()
RETURNING started_at`
	var startedAt time.Time
	if err := r.db.Pool.QueryRow(ctx, query, userID, courseID, totalModules).Scan(&startedAt); err != nil {
		return time.Time{}, err
	}
	return startedAt, nil
}

func (r *CourseRepository) IsUserEnrolled(ctx context.Context, userID, courseID int) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM user_course_progress WHERE user_id = $1 AND course_id = $2)`, userID, courseID).Scan(&exists)
	return exists, err
}

func (r *CourseRepository) GetCourseProgress(ctx context.Context, userID, courseID int) (int, int, []ModuleProgress, error) {
	modules, err := r.ListModules(ctx, courseID)
	if err != nil {
		return 0, 0, nil, err
	}

	progressRows, err := r.db.Pool.Query(ctx, `SELECT module_id, is_completed, completion_date FROM user_progress WHERE user_id = $1 AND course_id = $2`, userID, courseID)
	if err != nil {
		return 0, 0, nil, err
	}
	defer progressRows.Close()

	progressMap := make(map[int]ModuleProgress)
	for progressRows.Next() {
		var moduleID int
		var completed bool
		var completedAt *time.Time
		if err := progressRows.Scan(&moduleID, &completed, &completedAt); err != nil {
			return 0, 0, nil, err
		}
		progressMap[moduleID] = ModuleProgress{
			Module:      Module{ID: moduleID},
			Completed:   completed,
			CompletedAt: completedAt,
		}
	}

	var completedCount int
	moduleProgresses := make([]ModuleProgress, 0, len(modules))
	for _, m := range modules {
		mp := ModuleProgress{Module: m}
		if existing, ok := progressMap[m.ID]; ok {
			mp.Completed = existing.Completed
			mp.CompletedAt = existing.CompletedAt
			if mp.Completed {
				completedCount++
			}
		}
		moduleProgresses = append(moduleProgresses, mp)
	}

	return len(modules), completedCount, moduleProgresses, nil
}

func (r *CourseRepository) CompleteModule(ctx context.Context, userID, courseID, moduleID int) (*time.Time, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var belongsToCourse bool
	if err := tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM modules WHERE id = $1 AND course_id = $2)", moduleID, courseID).Scan(&belongsToCourse); err != nil {
		return nil, err
	}
	if !belongsToCourse {
		return nil, sql.ErrNoRows
	}

	completedAt := time.Now()
	_, err = tx.Exec(ctx, `
INSERT INTO user_progress (user_id, course_id, module_id, is_completed, completion_date)
VALUES ($1, $2, $3, TRUE, $4)
ON CONFLICT (user_id, module_id)
DO UPDATE SET is_completed = EXCLUDED.is_completed, completion_date = EXCLUDED.completion_date, course_id = EXCLUDED.course_id`, userID, courseID, moduleID, completedAt)
	if err != nil {
		return nil, err
	}

	var totalModules int
	if err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM modules WHERE course_id = $1", courseID).Scan(&totalModules); err != nil {
		return nil, err
	}
	var completedModules int
	if err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM user_progress WHERE user_id = $1 AND course_id = $2 AND is_completed = TRUE", userID, courseID).Scan(&completedModules); err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO user_course_progress (user_id, course_id, total_modules, completed_modules, module_progress_pct)
VALUES ($1, $2, $3, $4, CASE WHEN $3 > 0 THEN ($4 * 100) / $3 ELSE 0 END)
ON CONFLICT (user_id, course_id)
DO UPDATE SET total_modules = EXCLUDED.total_modules,
              completed_modules = EXCLUDED.completed_modules,
              module_progress_pct = EXCLUDED.module_progress_pct,
              updated_at = now()`, userID, courseID, totalModules, completedModules)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &completedAt, nil
}
