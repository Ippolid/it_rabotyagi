-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;

-- Создание ENUM типов через функцию
CREATE OR REPLACE FUNCTION create_enum_types() RETURNS void AS $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('user', 'mentor', 'admin');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'subscription_type') THEN
        CREATE TYPE subscription_type AS ENUM ('free', 'pro', 'team', 'enterprise');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'technology_enum') THEN
        CREATE TYPE technology_enum AS ENUM (
            'python', 'javascript', 'typescript', 'java', 'go', 'csharp', 'php', 'ruby',
            'kotlin', 'swift', 'sql', 'devops', 'ml', 'other'
            );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'difficulty_enum') THEN
        CREATE TYPE difficulty_enum AS ENUM ('easy', 'medium', 'hard');
    END IF;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

SELECT create_enum_types();
DROP FUNCTION create_enum_types();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_timestamp_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    IF TG_ARGV[0] = 'edited_at' THEN
        NEW.edited_at := now();
    ELSE
        NEW.updated_at := now();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- USERS
CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     telegram_id TEXT UNIQUE,
                                     google_id TEXT UNIQUE,
                                     github_id TEXT UNIQUE,
                                     email CITEXT UNIQUE CHECK (length(email) <= 250),
                                     username CITEXT UNIQUE CHECK (length(username) <= 50),
                                     name TEXT CHECK (length(name) <= 50),
                                     avatar_url TEXT,
                                     description TEXT CHECK (length(description) <= 150),
                                     role user_role NOT NULL DEFAULT 'user',
                                     subscription_type subscription_type,
                                     subscription_expires_at timestamptz,
                                     created_at timestamptz NOT NULL DEFAULT now(),
                                     updated_at timestamptz NOT NULL DEFAULT now()
);

-- MENTORS
CREATE TABLE IF NOT EXISTS mentors (
                                       id BIGSERIAL PRIMARY KEY,
                                       user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                       specialization TEXT NOT NULL,
                                       grade TEXT,
                                       experience_years INT,
                                       description TEXT,
                                       tags TEXT[],
                                       contacts JSONB,
                                       pricelist JSONB,
                                       created_at timestamptz NOT NULL DEFAULT now(),
                                       updated_at timestamptz NOT NULL DEFAULT now(),
                                       UNIQUE (user_id)
);

-- COURSES
CREATE TABLE IF NOT EXISTS courses (
                                       id BIGSERIAL PRIMARY KEY,
                                       title TEXT NOT NULL,
                                       description TEXT,
                                       technology technology_enum,
                                       is_published BOOLEAN NOT NULL DEFAULT FALSE,
                                       created_at timestamptz NOT NULL DEFAULT now(),
                                       updated_at timestamptz NOT NULL DEFAULT now()
);

-- MODULES
CREATE TABLE IF NOT EXISTS modules (
                                       id BIGSERIAL PRIMARY KEY,
                                       course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
                                       title TEXT NOT NULL,
                                       description TEXT,
                                       module_order INT NOT NULL,
                                       created_at timestamptz NOT NULL DEFAULT now(),
                                       edited_at timestamptz NOT NULL DEFAULT now(),
                                       UNIQUE(course_id, module_order)
);

-- QUESTIONS
CREATE TABLE IF NOT EXISTS questions (
                                         id BIGSERIAL PRIMARY KEY,
                                         title TEXT NOT NULL,
                                         content TEXT NOT NULL,
                                         technology technology_enum,
                                         difficulty difficulty_enum,
                                         options JSONB,
                                         correct_answer TEXT,
                                         explanation TEXT,
                                         company_tag TEXT[],
                                         created_at timestamptz NOT NULL DEFAULT now(),
                                         updated_at timestamptz NOT NULL DEFAULT now()
);

-- MODULE_QUESTIONS
CREATE TABLE IF NOT EXISTS module_questions (
                                                module_id BIGINT NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
                                                question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
                                                question_order INT NOT NULL,
                                                created_at timestamptz NOT NULL DEFAULT now(),
                                                PRIMARY KEY (module_id, question_id),
                                                UNIQUE(module_id, question_order)
);

-- USER_COURSE_PROGRESS
CREATE TABLE IF NOT EXISTS user_course_progress (
                                                    id BIGSERIAL PRIMARY KEY,
                                                    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                                    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
                                                    total_modules INT NOT NULL DEFAULT 0,
                                                    completed_modules INT NOT NULL DEFAULT 0,
                                                    module_progress_pct INT GENERATED ALWAYS AS (
                                                        CASE WHEN total_modules > 0
                                                                 THEN (completed_modules * 100) / total_modules
                                                             ELSE 0 END
                                                        ) STORED,
                                                    started_at timestamptz NOT NULL DEFAULT now(),
                                                    updated_at timestamptz NOT NULL DEFAULT now(),
                                                    UNIQUE (user_id, course_id)
);

-- USER_PROGRESS
CREATE TABLE IF NOT EXISTS user_progress (
                                             id BIGSERIAL PRIMARY KEY,
                                             user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                             course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
                                             module_id BIGINT NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
                                             is_completed BOOLEAN NOT NULL DEFAULT FALSE,
                                             completion_date timestamptz,
                                             created_at timestamptz NOT NULL DEFAULT now(),
                                             UNIQUE (user_id, module_id)
);

-- USER_QUESTION_PROGRESS
CREATE TABLE IF NOT EXISTS user_question_progress (
                                                      id BIGSERIAL PRIMARY KEY,
                                                      user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                                      course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
                                                      module_id BIGINT NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
                                                      question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
                                                      is_correct BOOLEAN NOT NULL DEFAULT FALSE,
                                                      attempts INT NOT NULL DEFAULT 1,
                                                      time_spent INTERVAL,
                                                      last_answer TEXT,
                                                      answered_at timestamptz NOT NULL DEFAULT now(),
                                                      updated_at timestamptz NOT NULL DEFAULT now(),
                                                      UNIQUE (user_id, question_id)
);

-- Триггеры
CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_timestamp_updated_at();

CREATE TRIGGER trg_mentors_updated_at
    BEFORE UPDATE ON mentors
    FOR EACH ROW EXECUTE FUNCTION set_timestamp_updated_at();

CREATE TRIGGER trg_courses_updated_at
    BEFORE UPDATE ON courses
    FOR EACH ROW EXECUTE FUNCTION set_timestamp_updated_at();

CREATE TRIGGER trg_modules_edited_at
    BEFORE UPDATE ON modules
    FOR EACH ROW EXECUTE FUNCTION set_timestamp_updated_at('edited_at');

CREATE TRIGGER trg_questions_updated_at
    BEFORE UPDATE ON questions
    FOR EACH ROW EXECUTE FUNCTION set_timestamp_updated_at();

CREATE TRIGGER trg_user_course_progress_updated_at
    BEFORE UPDATE ON user_course_progress
    FOR EACH ROW EXECUTE FUNCTION set_timestamp_updated_at();

CREATE TRIGGER trg_user_question_progress_updated_at
    BEFORE UPDATE ON user_question_progress
    FOR EACH ROW EXECUTE FUNCTION set_timestamp_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_user_question_progress_updated_at ON user_question_progress;
DROP TRIGGER IF EXISTS trg_user_course_progress_updated_at ON user_course_progress;
DROP TRIGGER IF EXISTS trg_questions_updated_at ON questions;
DROP TRIGGER IF EXISTS trg_modules_edited_at ON modules;
DROP TRIGGER IF EXISTS trg_courses_updated_at ON courses;
DROP TRIGGER IF EXISTS trg_mentors_updated_at ON mentors;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;

DROP FUNCTION IF EXISTS set_timestamp_updated_at;

DROP TABLE IF EXISTS user_question_progress;
DROP TABLE IF EXISTS user_progress;
DROP TABLE IF EXISTS user_course_progress;
DROP TABLE IF EXISTS module_questions;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS mentors;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS difficulty_enum;
DROP TYPE IF EXISTS technology_enum;
DROP TYPE IF EXISTS subscription_type;
DROP TYPE IF EXISTS user_role;

DROP EXTENSION IF EXISTS citext;
