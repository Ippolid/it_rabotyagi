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

type ModuleQuestion struct {
	QuestionID    int
	QuestionTitle string
	QuestionOrder int
	Difficulty    string
	IsAnswered    bool
	IsCorrect     *bool
}

type ModuleWithProgress struct {
	Module      Module
	Questions   []ModuleQuestion
	IsLocked    bool
	IsCompleted bool
	Progress    ModuleProgressStats
}

type ModuleProgressStats struct {
	TotalQuestions    int
	AnsweredQuestions int
	CorrectAnswers    int
	CorrectnessPct    float32
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
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return time.Time{}, err
	}
	defer tx.Rollback(ctx)

	// Insert or retain existing enrollment
	query := `
INSERT INTO user_course_progress (user_id, course_id, total_modules, completed_modules, module_progress_pct)
VALUES ($1, $2, $3, 0, 0)
ON CONFLICT (user_id, course_id) DO UPDATE SET updated_at = now()
RETURNING started_at`
	var startedAt time.Time
	if err := tx.QueryRow(ctx, query, userID, courseID, totalModules).Scan(&startedAt); err != nil {
		return time.Time{}, err
	}

	// Create user_progress entries for all modules if not already exist
	initProgressQuery := `
INSERT INTO user_progress (user_id, course_id, module_id, is_completed)
SELECT $1, $2, id, false
FROM modules
WHERE course_id = $2
ON CONFLICT (user_id, module_id) DO NOTHING`
	_, err = tx.Exec(ctx, initProgressQuery, userID, courseID)
	if err != nil {
		return time.Time{}, err
	}

	if err := tx.Commit(ctx); err != nil {
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

// GetModuleQuestions возвращает список вопросов модуля с прогрессом пользователя
func (r *CourseRepository) GetModuleQuestions(ctx context.Context, moduleID int, userID *int) ([]ModuleQuestion, error) {
	var query string
	var args []any

	if userID != nil {
		query = `
SELECT
    q.id,
    q.title,
    mq.question_order,
    q.difficulty,
    COALESCE(uqp.question_id IS NOT NULL, false) as is_answered,
    uqp.is_correct
FROM module_questions mq
JOIN questions q ON q.id = mq.question_id
LEFT JOIN user_question_progress uqp ON uqp.question_id = q.id AND uqp.user_id = $2
WHERE mq.module_id = $1
ORDER BY mq.question_order`
		args = []any{moduleID, *userID}
	} else {
		// If no user, just show questions without progress
		query = `
SELECT
    q.id,
    q.title,
    mq.question_order,
    q.difficulty,
    false as is_answered,
    NULL::boolean as is_correct
FROM module_questions mq
JOIN questions q ON q.id = mq.question_id
WHERE mq.module_id = $1
ORDER BY mq.question_order`
		args = []any{moduleID}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []ModuleQuestion
	for rows.Next() {
		var q ModuleQuestion
		if err := rows.Scan(&q.QuestionID, &q.QuestionTitle, &q.QuestionOrder, &q.Difficulty, &q.IsAnswered, &q.IsCorrect); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

// AddQuestionToModule добавляет вопрос в модуль
func (r *CourseRepository) AddQuestionToModule(ctx context.Context, moduleID, questionID, questionOrder int) error {
	query := `INSERT INTO module_questions (module_id, question_id, question_order) VALUES ($1, $2, $3)`
	_, err := r.db.Pool.Exec(ctx, query, moduleID, questionID, questionOrder)
	return err
}

// RemoveQuestionFromModule удаляет вопрос из модуля
func (r *CourseRepository) RemoveQuestionFromModule(ctx context.Context, moduleID, questionID int) error {
	query := `DELETE FROM module_questions WHERE module_id = $1 AND question_id = $2`
	cmd, err := r.db.Pool.Exec(ctx, query, moduleID, questionID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// CanAccessModule проверяет, может ли пользователь получить доступ к модулю
// Правило: module_order == 1 OR previous module is_completed == true
func (r *CourseRepository) CanAccessModule(ctx context.Context, userID, courseID, moduleID int) (bool, error) {
	// Получаем module_order текущего модуля
	var moduleOrder int
	err := r.db.Pool.QueryRow(ctx, `SELECT module_order FROM modules WHERE id = $1 AND course_id = $2`, moduleID, courseID).Scan(&moduleOrder)
	if err != nil {
		return false, err
	}

	// Если это первый модуль, доступ разрешен
	if moduleOrder == 1 {
		return true, nil
	}

	// Проверяем, завершен ли предыдущий модуль
	query := `
SELECT COALESCE(up.is_completed, false)
FROM modules m
LEFT JOIN user_progress up ON up.module_id = m.id AND up.user_id = $1
WHERE m.course_id = $2 AND m.module_order = $3`

	var prevCompleted bool
	err = r.db.Pool.QueryRow(ctx, query, userID, courseID, moduleOrder-1).Scan(&prevCompleted)
	if err != nil {
		return false, err
	}

	return prevCompleted, nil
}

// IsModuleComplete проверяет, завершен ли модуль (все вопросы отвечены)
func (r *CourseRepository) IsModuleComplete(ctx context.Context, userID, moduleID int) (bool, error) {
	query := `
SELECT
    COUNT(mq.question_id) as total_questions,
    COUNT(uqp.question_id) as answered_questions
FROM module_questions mq
LEFT JOIN user_question_progress uqp ON uqp.question_id = mq.question_id AND uqp.user_id = $1
WHERE mq.module_id = $2`

	var total, answered int
	err := r.db.Pool.QueryRow(ctx, query, userID, moduleID).Scan(&total, &answered)
	if err != nil {
		return false, err
	}

	return total > 0 && total == answered, nil
}

// GetModuleWithProgress возвращает модуль с полной информацией о прогрессе
func (r *CourseRepository) GetModuleWithProgress(ctx context.Context, userID, courseID, moduleID int) (*ModuleWithProgress, error) {
	// Получаем модуль
	module, err := r.GetModule(ctx, courseID, moduleID)
	if err != nil {
		return nil, err
	}

	// Проверяем доступ
	canAccess, err := r.CanAccessModule(ctx, userID, courseID, moduleID)
	if err != nil {
		return nil, err
	}

	// Получаем вопросы
	questions, err := r.GetModuleQuestions(ctx, moduleID, &userID)
	if err != nil {
		return nil, err
	}

	// Проверяем завершенность
	var isCompleted bool
	err = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(is_completed, false) FROM user_progress WHERE user_id = $1 AND module_id = $2`, userID, moduleID).Scan(&isCompleted)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Подсчитываем статистику
	totalQuestions := len(questions)
	answeredQuestions := 0
	correctAnswers := 0
	for _, q := range questions {
		if q.IsAnswered {
			answeredQuestions++
			if q.IsCorrect != nil && *q.IsCorrect {
				correctAnswers++
			}
		}
	}

	correctnessPct := float32(0)
	if answeredQuestions > 0 {
		correctnessPct = float32(correctAnswers) * 100.0 / float32(answeredQuestions)
	}

	return &ModuleWithProgress{
		Module:      *module,
		Questions:   questions,
		IsLocked:    !canAccess,
		IsCompleted: isCompleted,
		Progress: ModuleProgressStats{
			TotalQuestions:    totalQuestions,
			AnsweredQuestions: answeredQuestions,
			CorrectAnswers:    correctAnswers,
			CorrectnessPct:    correctnessPct,
		},
	}, nil
}
