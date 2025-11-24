-- +goose Up
CREATE TABLE IF NOT EXISTS allowed_editors (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_allowed_editors_user_id ON allowed_editors (user_id);

-- +goose Down
DROP TABLE IF EXISTS allowed_editors;
