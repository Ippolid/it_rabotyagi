package repositories

import (
    "context"
    "encoding/json"
    "it_rabotyagi/internal/business/models"
    "it_rabotyagi/internal/data/database"
)

type QuestionRepository struct {
    db *database.DB
}

func NewQuestionRepository(db *database.DB) *QuestionRepository {
    return &QuestionRepository{db: db}
}

func (r *QuestionRepository) ListByModule(ctx context.Context, moduleID int64) ([]*models.Question, error) {
    rows, err := r.db.Pool.Query(ctx, `
        SELECT q.id, q.title, q.content, q.difficulty, q.options, q.correct_answer, q.explanation, q.created_at, q.updated_at
        FROM questions q
        JOIN module_questions mq ON mq.question_id = q.id
        WHERE mq.module_id = $1
        ORDER BY mq.question_order ASC, q.id ASC
    `, moduleID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []*models.Question
    for rows.Next() {
        q := &models.Question{}
        var rawOptions []byte
        if err := rows.Scan(&q.ID, &q.Title, &q.Content, &q.Difficulty, &rawOptions, &q.CorrectAnswer, &q.Explanation, &q.CreatedAt, &q.UpdatedAt); err != nil {
            return nil, err
        }
        if len(rawOptions) > 0 {
            var opts []string
            if err := json.Unmarshal(rawOptions, &opts); err == nil {
                q.Options = opts
            }
        }
        out = append(out, q)
    }
    return out, nil
}
