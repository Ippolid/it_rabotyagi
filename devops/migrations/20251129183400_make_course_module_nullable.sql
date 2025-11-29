-- +goose Up
-- Make course_id and module_id nullable in user_question_progress
-- to support standalone questions (not part of a course/module)
ALTER TABLE user_question_progress
    ALTER COLUMN course_id DROP NOT NULL,
    ALTER COLUMN module_id DROP NOT NULL;

-- +goose Down
-- Revert back to NOT NULL constraints
-- Note: This will fail if there are NULL values in the columns
ALTER TABLE user_question_progress
    ALTER COLUMN course_id SET NOT NULL,
    ALTER COLUMN module_id SET NOT NULL;