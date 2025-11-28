package repositories

import (
	"context"
	"it_rabotyagi/internal/data/database"
	"time"

	"github.com/jackc/pgx/v5"
)

type StatisticsRepository struct {
	db *database.DB
}

func NewStatisticsRepository(db *database.DB) *StatisticsRepository {
	return &StatisticsRepository{db: db}
}

// Типы данных
type OverallStatistics struct {
	CoursesEnrolled    int
	CoursesCompleted   int
	OverallProgressPct float32
	QuestionStatistics QuestionStatistics
}

type CourseStatistics struct {
	CourseID         int
	CourseTitle      string
	TotalModules     int
	CompletedModules int
	ProgressPct      float32
	StartedAt        time.Time
	TimeSpent        float32
}

type QuestionStatistics struct {
	TotalSolved       int
	CorrectAnswersPct float32
	TotalAttempts     int
	AvgTimeSpent      float32
}

type QuestionStatisticsByCourse struct {
	CourseID    int
	CourseTitle string
	Statistics  QuestionStatistics
}

// GetOverallStatistics получает общую статистику пользователя
func (r *StatisticsRepository) GetOverallStatistics(ctx context.Context, userID int) (*OverallStatistics, error) {
	query := `
		WITH course_stats AS (
			SELECT
				COUNT(*) as enrolled,
				COUNT(CASE WHEN completed_modules = total_modules AND total_modules > 0 THEN 1 END) as completed,
				AVG(module_progress_pct) as avg_progress
			FROM user_course_progress
			WHERE user_id = $1
		),
		question_stats AS (
			SELECT
				COUNT(DISTINCT question_id) as total_solved,
				COUNT(CASE WHEN is_correct THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0) as correct_pct,
				SUM(attempts) as total_attempts,
				AVG(EXTRACT(EPOCH FROM time_spent)) as avg_time
			FROM user_question_progress
			WHERE user_id = $1
		)
		SELECT
			COALESCE(cs.enrolled, 0),
			COALESCE(cs.completed, 0),
			COALESCE(cs.avg_progress, 0),
			COALESCE(qs.total_solved, 0),
			COALESCE(qs.correct_pct, 0),
			COALESCE(qs.total_attempts, 0),
			COALESCE(qs.avg_time, 0)
		FROM course_stats cs
		CROSS JOIN question_stats qs`

	stats := &OverallStatistics{}
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&stats.CoursesEnrolled,
		&stats.CoursesCompleted,
		&stats.OverallProgressPct,
		&stats.QuestionStatistics.TotalSolved,
		&stats.QuestionStatistics.CorrectAnswersPct,
		&stats.QuestionStatistics.TotalAttempts,
		&stats.QuestionStatistics.AvgTimeSpent,
	)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetCourseStatistics получает детальную статистику по курсам
func (r *StatisticsRepository) GetCourseStatistics(ctx context.Context, userID int) ([]CourseStatistics, error) {
	query := `
		SELECT
			ucp.course_id,
			c.title,
			ucp.total_modules,
			ucp.completed_modules,
			ucp.module_progress_pct,
			ucp.started_at,
			EXTRACT(EPOCH FROM (ucp.updated_at - ucp.started_at)) / 3600.0
		FROM user_course_progress ucp
		JOIN courses c ON c.id = ucp.course_id
		WHERE ucp.user_id = $1
		ORDER BY ucp.started_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CourseStatistics
	for rows.Next() {
		var stat CourseStatistics
		err := rows.Scan(
			&stat.CourseID,
			&stat.CourseTitle,
			&stat.TotalModules,
			&stat.CompletedModules,
			&stat.ProgressPct,
			&stat.StartedAt,
			&stat.TimeSpent,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, stat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetQuestionStatisticsOverall получает общую статистику по вопросам
func (r *StatisticsRepository) GetQuestionStatisticsOverall(ctx context.Context, userID int) (*QuestionStatistics, error) {
	query := `
		SELECT
			COALESCE(COUNT(DISTINCT question_id), 0),
			COALESCE(COUNT(CASE WHEN is_correct THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0), 0),
			COALESCE(SUM(attempts), 0),
			COALESCE(AVG(EXTRACT(EPOCH FROM time_spent)), 0)
		FROM user_question_progress
		WHERE user_id = $1`

	stats := &QuestionStatistics{}
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&stats.TotalSolved,
		&stats.CorrectAnswersPct,
		&stats.TotalAttempts,
		&stats.AvgTimeSpent,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Если нет данных, возвращаем нулевую статистику
			return stats, nil
		}
		return nil, err
	}

	return stats, nil
}

// GetQuestionStatisticsByCourse получает статистику по вопросам с разбивкой по курсам
func (r *StatisticsRepository) GetQuestionStatisticsByCourse(ctx context.Context, userID int, courseID *int) ([]QuestionStatisticsByCourse, error) {
	query := `
		SELECT
			uqp.course_id,
			c.title,
			COUNT(DISTINCT uqp.question_id),
			COUNT(CASE WHEN uqp.is_correct THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0),
			SUM(uqp.attempts),
			AVG(EXTRACT(EPOCH FROM uqp.time_spent))
		FROM user_question_progress uqp
		JOIN courses c ON c.id = uqp.course_id
		WHERE uqp.user_id = $1
		  AND ($2::int IS NULL OR uqp.course_id = $2)
		GROUP BY uqp.course_id, c.title
		ORDER BY uqp.course_id`

	rows, err := r.db.Pool.Query(ctx, query, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []QuestionStatisticsByCourse
	for rows.Next() {
		var stat QuestionStatisticsByCourse
		err := rows.Scan(
			&stat.CourseID,
			&stat.CourseTitle,
			&stat.Statistics.TotalSolved,
			&stat.Statistics.CorrectAnswersPct,
			&stat.Statistics.TotalAttempts,
			&stat.Statistics.AvgTimeSpent,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, stat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}