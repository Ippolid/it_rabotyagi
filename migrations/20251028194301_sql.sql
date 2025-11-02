-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE CHECK (length(username) <= 50) NOT NULL,
    password TEXT NOT NULL,

    email TEXT UNIQUE CHECK (length(email) <= 250),
    telegram_id TEXT UNIQUE,
    google_id TEXT UNIQUE,
    github_id TEXT UNIQUE,

    name TEXT CHECK (length(name) <= 50),
    avatar_url TEXT,
    description TEXT CHECK (length(description) <= 150),
    role TEXT NOT NULL DEFAULT 'user',
    subscription_type TEXT,
    subscription_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE mentors (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    specialization TEXT NOT NULL,
    grade TEXT,
    experience_years INT,
    description TEXT,
    tags TEXT[],
    contacts JSONB,
    pricelist JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id)
);

CREATE TABLE technologies (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL ,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE course_technologies (
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    technology_id BIGINT NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (course_id, technology_id)
);

CREATE TABLE modules (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    module_order INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    edited_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (course_id, module_order)
);

CREATE TABLE questions (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    difficulty TEXT,
    options JSONB,
    correct_answer TEXT,
    explanation TEXT,
    company_tag TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE question_technologies (
    question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    technology_id BIGINT NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (question_id, technology_id)
);

CREATE TABLE module_questions (
    module_id BIGINT NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    question_order INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (module_id, question_id),
    UNIQUE (module_id, question_order)
);

CREATE TABLE user_course_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    total_modules INT NOT NULL DEFAULT 0,
    completed_modules INT NOT NULL DEFAULT 0,
    module_progress_pct INT,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, course_id)
);

CREATE TABLE user_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    module_id BIGINT NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completion_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, module_id)
);

CREATE TABLE user_question_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    module_id BIGINT NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    attempts INT NOT NULL DEFAULT 1,
    time_spent INTERVAL,
    last_answer TEXT,
    answered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, question_id)
);

CREATE TABLE auth_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at TIMESTAMPTZ
);


-- +goose Down
DROP TABLE IF EXISTS user_question_progress;
DROP TABLE IF EXISTS user_progress;
DROP TABLE IF EXISTS user_course_progress;
DROP TABLE IF EXISTS module_questions;
DROP TABLE IF EXISTS question_technologies;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS course_technologies;
DROP TABLE IF EXISTS technologies;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS mentors;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS auth_sessions;