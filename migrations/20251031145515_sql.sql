-- +goose Up
-- Обертываем все запросы в транзакцию для атомарности
BEGIN;

-- Вставляем технологии, указывая ID вручную.
-- ON CONFLICT (id) DO NOTHING игнорирует вставку, если запись с таким ID уже существует.
INSERT INTO technologies (id, name) VALUES (1, 'Go') ON CONFLICT (id) DO NOTHING;
INSERT INTO technologies (id, name) VALUES (2, 'Python') ON CONFLICT (id) DO NOTHING;

-- Вставляем курс
INSERT INTO courses (id, title, description, is_published)
VALUES (1, 'Программирование на Go', 'Основы и продвинутые концепции языка Go', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Устанавливаем связь между курсом и технологией
INSERT INTO course_technologies (course_id, technology_id) VALUES (1, 1);

-- Вставляем модули для курса (course_id = 1)
INSERT INTO modules (id, course_id, title, description, module_order, edited_at)
VALUES
    (1, 1, 'Введение в Go', 'Знакомство с языком Go и настройка рабочего окружения.', 1, now()),
    (2, 1, 'Основы синтаксиса', 'Типы данных, переменные, функции и управление потоком.', 2, now()),
    (3, 1, 'Конкурентность в Go', 'Горутины, каналы и продвинутые паттерны конкурентного программирования.', 3, now())
ON CONFLICT (id) DO NOTHING;

-- Вставляем вопросы
INSERT INTO questions (id, title, content, difficulty, options, correct_answer, explanation)
VALUES
    (1, 'Что такое слайс (slice) в Go?', 'Выберите наиболее точное определение для слайса.', 'easy',
     '["Вариант A: Массив фиксированной длины", "Вариант B: Динамический массив, который ссылается на базовый массив", "Вариант C: Структура для хранения ключ-значение"]',
     'Вариант B: Динамический массив, который ссылается на базовый массив',
     'Слайс представляет собой гибкую обертку над массивом, позволяя работать с последовательностями данных динамически.'),
    (2, 'С помощью какого ключевого слова запускается горутина?', 'Как запустить функцию myFunc() в отдельной горутине?', 'easy',
     '["Вариант A: goroutine myFunc()", "Вариант B: start myFunc()", "Вариант C: go myFunc()"]',
     'Вариант C: go myFunc()',
     'Ключевое слово `go` используется перед вызовом функции для ее асинхронного выполнения в новой горутине.')
ON CONFLICT (id) DO NOTHING;

-- Связываем вопросы с технологией Go
INSERT INTO question_technologies (question_id, technology_id) VALUES (1, 1), (2, 1);

-- Добавляем вопросы в модуль (в данном случае, в модуль с id=2 "Основы синтаксиса")
INSERT INTO module_questions (module_id, question_id, question_order) VALUES (2, 1, 1), (2, 2, 2);

-- Создаем пользователя и ментора
INSERT INTO users (id, username, password, email, name, role)
VALUES (1, 'gopher123', 'hashed_password', 'gopher123@example.com', 'Alex', 'user')
ON CONFLICT (id) DO NOTHING;

INSERT INTO mentors (id, user_id, specialization, grade, experience_years, tags)
VALUES (1, 1, 'Backend Development', 'Senior', 5, '{"Go", "PostgreSQL", "Microservices"}')
ON CONFLICT (id) DO NOTHING;

-- Добавляем пример прогресса для пользователя
INSERT INTO user_course_progress (id, user_id, course_id, total_modules, completed_modules, module_progress_pct)
VALUES (1, 1, 1, 3, 1, 33)
ON CONFLICT (id) DO NOTHING;

COMMIT;

-- +goose Down
-- Эта команда полностью очистит таблицы от всех данных и сбросит счетчики.
-- Идеально для отката тестового наполнения.
BEGIN;

TRUNCATE TABLE
    user_question_progress,
    user_progress,
    user_course_progress,
    module_questions,
    question_technologies,
    questions,
    modules,
    course_technologies,
    courses,
    mentors,
    users,
    auth_sessions,
    technologies
    RESTART IDENTITY CASCADE;

COMMIT;
