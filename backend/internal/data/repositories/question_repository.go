package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

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
	Tags          []string `json:"tags"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correct_answer"`
	Explanation   *string  `json:"explanation"`
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
		SELECT q.id,
		       q.title,
		       q.content,
		       q.difficulty,
		       q.options,
		       q.correct_answer,
		       q.explanation,
		       q.company_tag,
		       t.name as technology
		FROM questions q
		JOIN question_technologies qt ON q.id = qt.question_id
		JOIN technologies t ON qt.technology_id = t.id
		WHERE q.id = $1
		LIMIT 1
	`

	var q QuestionDetail
	var optionsJSON []byte
	var tags []string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&q.ID,
		&q.Title,
		&q.Content,
		&q.Difficulty,
		&optionsJSON,
		&q.CorrectAnswer,
		&q.Explanation,
		&tags,
		&q.Technology,
	)

	if err != nil {
		return nil, err
	}

	// Парсим JSON с вариантами ответов
	if err := json.Unmarshal(optionsJSON, &q.Options); err != nil {
		return nil, err
	}
	q.Tags = tags

	return &q, nil
}

func (r *QuestionRepository) ensureTechnology(ctx context.Context, technology string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `INSERT INTO technologies (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, technology).Scan(&id)
	return id, err
}

func (r *QuestionRepository) CreateQuestion(ctx context.Context, title, content, difficulty, technology string, tags []string, options []string, correctAnswer string, explanation *string) (int, error) {
	techID, err := r.ensureTechnology(ctx, technology)
	if err != nil {
		return 0, err
	}

	optsJSON, err := json.Marshal(options)
	if err != nil {
		return 0, err
	}

	var questionID int
	err = r.db.QueryRow(ctx, `INSERT INTO questions (title, content, difficulty, options, correct_answer, explanation, company_tag) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		title, content, difficulty, optsJSON, correctAnswer, explanation, tags).Scan(&questionID)
	if err != nil {
		return 0, err
	}

	if _, err := r.db.Exec(ctx, `INSERT INTO question_technologies (question_id, technology_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, questionID, techID); err != nil {
		return 0, err
	}

	return questionID, nil
}

func (r *QuestionRepository) UpdateQuestion(ctx context.Context, id int, title, content, difficulty, technology *string, tags *[]string, options *[]string, correctAnswer, explanation *string) error {
	setParts := []string{}
	args := []any{}
	idx := 1

	if title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", idx))
		args = append(args, *title)
		idx++
	}
	if content != nil {
		setParts = append(setParts, fmt.Sprintf("content = $%d", idx))
		args = append(args, *content)
		idx++
	}
	if difficulty != nil {
		setParts = append(setParts, fmt.Sprintf("difficulty = $%d", idx))
		args = append(args, *difficulty)
		idx++
	}
	if options != nil {
		optsJSON, err := json.Marshal(*options)
		if err != nil {
			return err
		}
		setParts = append(setParts, fmt.Sprintf("options = $%d", idx))
		args = append(args, optsJSON)
		idx++
	}
	if correctAnswer != nil {
		setParts = append(setParts, fmt.Sprintf("correct_answer = $%d", idx))
		args = append(args, *correctAnswer)
		idx++
	}
	if explanation != nil {
		setParts = append(setParts, fmt.Sprintf("explanation = $%d", idx))
		args = append(args, *explanation)
		idx++
	}
	if tags != nil {
		setParts = append(setParts, fmt.Sprintf("company_tag = $%d", idx))
		args = append(args, *tags)
		idx++
	}

	if len(setParts) == 0 && technology == nil {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if len(setParts) > 0 {
		query := fmt.Sprintf("UPDATE questions SET %s, updated_at = now() WHERE id = $%d", strings.Join(setParts, ", "), idx)
		argsWithID := append(args, id)
		cmd, err := tx.Exec(ctx, query, argsWithID...)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return sql.ErrNoRows
		}
	}

	if technology != nil {
		techID, err := r.ensureTechnology(ctx, *technology)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, "DELETE FROM question_technologies WHERE question_id = $1", id); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, "INSERT INTO question_technologies (question_id, technology_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", id, techID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *QuestionRepository) DeleteQuestion(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, "DELETE FROM questions WHERE id = $1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}
