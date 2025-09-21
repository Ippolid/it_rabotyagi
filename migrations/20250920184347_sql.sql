-- +goose Up
-- Создаем перечисляемые типы для ролей и подписок
CREATE TYPE user_role AS ENUM ('user', 'mentor');
CREATE TYPE subscription_type AS ENUM ('trial', 'pro');

-- Таблица пользователей
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,                    -- Внутренний ID (автоинкремент)
    telegram_id BIGINT UNIQUE NOT NULL,          -- Telegram ID (уникальный, но не PK)
    username VARCHAR(255),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    photo_url TEXT,
    role user_role NOT NULL DEFAULT 'user',      -- Роль пользователя ('user' или 'mentor')
    subscription subscription_type NOT NULL DEFAULT 'trial', -- Тип подписки ('trial' или 'pro')
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_users_telegram_id ON users(telegram_id);
CREATE INDEX idx_users_username ON users(username) WHERE username IS NOT NULL;
CREATE INDEX idx_users_role ON users(role); -- Индекс по роли
CREATE INDEX idx_users_subscription ON users(subscription); -- Индекс по типу подписки

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS subscription_type;
