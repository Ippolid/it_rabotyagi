-- +goose Up
BEGIN;

-- Добавляем открытые вопросы (без вариантов ответа, без correct_answer)
-- Используем E-escaping для экранирования символов
-- ID генерируется автоматически через SERIAL
-- Проверяем, что вопросы еще не существуют перед вставкой
INSERT INTO questions (title, content, difficulty, options, correct_answer, explanation, company_tag)
SELECT * FROM (VALUES
    (
        'Что такое замыкания (closures) в JavaScript?',
        'Объясните концепцию замыканий в JavaScript. Приведите примеры использования и объясните, для чего они нужны.',
        'medium',
        '[]'::jsonb,
        '',
        E'Замыкание (closure) — это функция, которая имеет доступ к переменным из внешней (объемлющей) функции, даже после того, как внешняя функция завершила свою работу. Замыкания создаются каждый раз, когда функция создается внутри другой функции.\n\nПример:\nfunction outer() {\n  const message = "Hello";\n  return function inner() {\n    console.log(message);\n  };\n}\nconst closure = outer();\nclosure(); // "Hello"\n\nЗамыкания часто используются для:\n- Создания приватных переменных и методов\n- Реализации паттерна "модуль"\n- Каррирования и частичного применения функций\n- Создания функций-фабрик\n- Сохранения состояния между вызовами функции\n\nВажно: замыкания сохраняют ссылку на переменную, а не ее значение на момент создания.',
        ARRAY['Google', 'Yandex']
    ),
    (
        'Объясните разницу между SQL и NoSQL базами данных',
        'В каких случаях стоит использовать SQL, а в каких NoSQL? Приведите примеры конкретных баз данных и сценариев использования.',
        'medium',
        '[]'::jsonb,
        '',
        E'SQL базы данных (реляционные):\n- Структурированные данные с четкой схемой\n- ACID транзакции для обеспечения консистентности\n- Сложные JOIN операции между таблицами\n- Вертикальное масштабирование (scaling up)\n- Примеры: PostgreSQL, MySQL, Oracle\n- Используйте для: банковских систем, ERP, систем с четкими отношениями между сущностями\n\nNoSQL базы данных (нереляционные):\n- Гибкая или отсутствующая схема данных\n- Горизонтальное масштабирование (scaling out)\n- Высокая производительность для больших объемов данных\n- Различные модели данных: документы, key-value, графы, колоночные\n- Примеры: MongoDB (документы), Redis (key-value), Cassandra (колоночная), Neo4j (графовая)\n- Используйте для: соцсетей, real-time аналитики, кеширования, IoT данных\n\nВыбор зависит от:\n- Требований к консистентности данных (ACID vs BASE)\n- Необходимости масштабирования\n- Структуры и изменчивости данных\n- Паттернов доступа к данным\n- Требований к производительности',
        ARRAY['Meta', 'Amazon', 'Yandex']
    ),
    (
        'Что такое горутины в Go и чем они отличаются от потоков?',
        'Объясните концепцию горутин в Go. В чем их преимущества перед OS потоками? Как работает Go scheduler?',
        'medium',
        '[]'::jsonb,
        '',
        E'Горутины (goroutines) — это легковесные потоки выполнения, управляемые Go runtime, а не операционной системой.\n\nОтличия от OS потоков:\n1. Размер стека: ~2KB для горутины vs ~1-2MB для OS потока\n2. Управление: Go scheduler (в user space) vs OS scheduler\n3. Переключение контекста: намного быстрее у горутин\n4. Стоимость создания: очень низкая для горутин\n5. Количество: можно создать миллионы горутин без проблем\n\nПреимущества:\n- Эффективное использование CPU за счет M:N scheduling\n- Простой синтаксис: go funcName()\n- Встроенная синхронизация через каналы (channels)\n- Автоматическое управление жизненным циклом\n- Минимальные накладные расходы\n\nПример:\nfunc main() {\n    for i := 0; i < 1000; i++ {\n        go func(n int) {\n            fmt.Println(n)\n        }(i)\n    }\n    time.Sleep(time.Second)\n}\n\nGo scheduler:\n- Использует M:N модель: M горутин выполняются на N OS потоках\n- Работа с run queue и локальными очередями для каждого процессора\n- Work stealing для балансировки нагрузки\n- Кооперативное переключение в точках safe points',
        ARRAY['Google', 'Yandex', 'VK']
    )
) AS new_questions(title, content, difficulty, options, correct_answer, explanation, company_tag)
WHERE NOT EXISTS (
    SELECT 1 FROM questions WHERE title = new_questions.title
);

-- Связываем вопросы с технологиями
INSERT INTO question_technologies (question_id, technology_id)
SELECT q.id, t.id
FROM questions q
JOIN technologies t ON t.name IN ('JavaScript', 'PostgreSQL', 'Go')
WHERE q.title IN (
    'Что такое замыкания (closures) в JavaScript?',
    'Объясните разницу между SQL и NoSQL базами данных',
    'Что такое горутины в Go и чем они отличаются от потоков?'
)
AND (
    (q.title = 'Что такое замыкания (closures) в JavaScript?' AND t.name = 'JavaScript') OR
    (q.title = 'Объясните разницу между SQL и NoSQL базами данных' AND t.name = 'PostgreSQL') OR
    (q.title = 'Что такое горутины в Go и чем они отличаются от потоков?' AND t.name = 'Go')
)
ON CONFLICT DO NOTHING;

COMMIT;

-- +goose Down
BEGIN;

-- Удаляем связи с технологиями
DELETE FROM question_technologies
WHERE question_id IN (
    SELECT id FROM questions
    WHERE title IN (
        'Что такое замыкания (closures) в JavaScript?',
        'Объясните разницу между SQL и NoSQL базами данных',
        'Что такое горутины в Go и чем они отличаются от потоков?'
    )
);

-- Удаляем вопросы
DELETE FROM questions
WHERE title IN (
    'Что такое замыкания (closures) в JavaScript?',
    'Объясните разницу между SQL и NoSQL базами данных',
    'Что такое горутины в Go и чем они отличаются от потоков?'
);

COMMIT;
