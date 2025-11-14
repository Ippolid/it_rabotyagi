package repositories

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionRepository struct {
	db *pgxpool.Pool
}

func NewQuestionRepository(db *pgxpool.Pool) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// QuestionListItem представляет краткую информацию о вопросе
type QuestionListItem struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Technology string `json:"technology"`
}

// QuestionDetail представляет полную информацию о вопросе
type QuestionDetail struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Difficulty    string   `json:"difficulty"`
	Technology    string   `json:"technology"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correct_answer"`
	Explanation   string   `json:"explanation"`
}

// GetAllQuestions получает список всех вопросов с их технологиями
func (r *QuestionRepository) GetAllQuestions(ctx context.Context, technology *string, limit, offset int) ([]QuestionListItem, int, error) {
	var questions []QuestionListItem
	var total int

	// Базовый запрос с JOIN для получения технологии
	baseQuery := `
		SELECT q.id, q.title, t.name as technology
		FROM questions q
		JOIN question_technologies qt ON q.id = qt.question_id
		JOIN technologies t ON qt.technology_id = t.id
	`

	countQuery := `
		SELECT COUNT(DISTINCT q.id)
		FROM questions q
		JOIN question_technologies qt ON q.id = qt.question_id
		JOIN technologies t ON qt.technology_id = t.id
	`

	// Добавляем фильтр по технологии, если указан
	if technology != nil && *technology != "" {
		baseQuery += " WHERE t.name = $1"
		countQuery += " WHERE t.name = $1"
	}

	// Добавляем сортировку и пагинацию
	baseQuery += " ORDER BY q.id"

	if limit > 0 {
		if technology != nil && *technology != "" {
			baseQuery += " LIMIT $2 OFFSET $3"
		} else {
			baseQuery += " LIMIT $1 OFFSET $2"
		}
	}

	// Получаем общее количество
	if technology != nil && *technology != "" {
		err := r.db.QueryRow(ctx, countQuery, *technology).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.db.QueryRow(ctx, countQuery).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}

	// Получаем вопросы
	var rows interface{ Close() }
	var err error

	if technology != nil && *technology != "" {
		if limit > 0 {
			rows, err = r.db.Query(ctx, baseQuery, *technology, limit, offset)
		} else {
			rows, err = r.db.Query(ctx, baseQuery, *technology)
		}
	} else {
		if limit > 0 {
			rows, err = r.db.Query(ctx, baseQuery, limit, offset)
		} else {
			rows, err = r.db.Query(ctx, baseQuery)
		}
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	pgxRows := rows.(interface {
		Next() bool
		Scan(...interface{}) error
		Close()
	})

	for pgxRows.Next() {
		var q QuestionListItem
		if err := pgxRows.Scan(&q.ID, &q.Title, &q.Technology); err != nil {
			return nil, 0, err
		}
		questions = append(questions, q)
	}

	return questions, total, nil
}

// GetQuestionByID получает полную информацию о вопросе по ID
func (r *QuestionRepository) GetQuestionByID(ctx context.Context, id int) (*QuestionDetail, error) {
	query := `
		SELECT q.id, q.title, q.content, q.difficulty, q.options, q.correct_answer, q.explanation, t.name as technology
		FROM questions q
		JOIN question_technologies qt ON q.id = qt.question_id
		JOIN technologies t ON qt.technology_id = t.id
		WHERE q.id = $1
		LIMIT 1
	`

	var q QuestionDetail
	var optionsJSON []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&q.ID,
		&q.Title,
		&q.Content,
		&q.Difficulty,
		&optionsJSON,
		&q.CorrectAnswer,
		&q.Explanation,
		&q.Technology,
	)

	if err != nil {
		return nil, err
	}

	// Парсим JSON с вариантами ответов
	if err := json.Unmarshal(optionsJSON, &q.Options); err != nil {
		return nil, err
	}

	return &q, nil
}
