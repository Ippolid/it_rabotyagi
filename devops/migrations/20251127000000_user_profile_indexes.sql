-- +goose Up
-- Индексы для оптимизации статистических запросов

CREATE INDEX IF NOT EXISTS idx_user_course_progress_user_id
ON user_course_progress(user_id);

CREATE INDEX IF NOT EXISTS idx_user_progress_user_id
ON user_progress(user_id);

CREATE INDEX IF NOT EXISTS idx_user_question_progress_user_course
ON user_question_progress(user_id, course_id);

CREATE INDEX IF NOT EXISTS idx_user_question_progress_user_correct
ON user_question_progress(user_id, is_correct);

CREATE INDEX IF NOT EXISTS idx_users_username
ON users(username) WHERE username IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_email
ON users(email) WHERE email IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_user_course_progress_user_id;
DROP INDEX IF EXISTS idx_user_progress_user_id;
DROP INDEX IF EXISTS idx_user_question_progress_user_course;
DROP INDEX IF EXISTS idx_user_question_progress_user_correct;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;