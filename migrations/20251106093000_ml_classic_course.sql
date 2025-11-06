-- +goose Up
BEGIN;

-- Technology: Machine Learning
INSERT INTO technologies (id, name)
VALUES (3, 'Machine Learning')
ON CONFLICT (id) DO NOTHING;

-- Course: Classic Machine Learning
INSERT INTO courses (id, title, description, is_published, created_at, updated_at)
VALUES (
  2,
  'Classic Machine Learning',
  'Полный практический курс по классическому ML: регрессия, классификация, ансамбли, SVM, кластеризация, понижение размерности и оценка качества.',
  TRUE,
  NOW(), NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Link course to technology
INSERT INTO course_technologies (course_id, technology_id)
VALUES (2, 3)
ON CONFLICT (course_id, technology_id) DO NOTHING;

-- Modules (10)
INSERT INTO modules (id, course_id, title, description, module_order, created_at, edited_at) VALUES
  (200, 2, 'Линейная регрессия', 'Модель, предположения, метрики, полиномиальная регрессия.', 1, NOW(), NOW()),
  (201, 2, 'Градиентный спуск', 'Оптимизация MSE: batch/SGD/mini-batch, выбор шага.', 2, NOW(), NOW()),
  (202, 2, 'Обобщающая способность и кросс-валидация', 'Bias–variance, train/val/test, схемы CV, data leakage.', 3, NOW(), NOW()),
  (203, 2, 'Регуляризация и масштабирование', 'L1/L2/ElasticNet, нормализация признаков.', 4, NOW(), NOW()),
  (204, 2, 'Логистическая регрессия и базовые метрики', 'Бинарная классификация, Confusion matrix, Accuracy/Precision/Recall/F1.', 5, NOW(), NOW()),
  (205, 2, 'ROC/PR и калибровка', 'ROC‑AUC vs PR‑AUC, калибровка вероятностей и выбор порога.', 6, NOW(), NOW()),
  (206, 2, 'SVM', 'Линейный/ядерный SVM, параметры C и γ, масштабирование.', 7, NOW(), NOW()),
  (207, 2, 'Деревья и Random Forest', 'Критерии, глубина, переобучение, бэггинг.', 8, NOW(), NOW()),
  (208, 2, 'Градиентный бустинг', 'Идея бустинга, XGBoost/LightGBM/CatBoost, важность признаков.', 9, NOW(), NOW()),
  (209, 2, 'Кластеризация и понижение размерности', 'k‑Means, DBSCAN, PCA, t‑SNE, EDA‑практики.', 10, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Questions (20 sample)
-- Note: difficulty: easy|medium|hard; options as JSON array where уместно
INSERT INTO questions (id, title, content, difficulty, options, correct_answer, explanation, company_tag, created_at, updated_at) VALUES
  (2000, 'Что минимизирует линейная регрессия?', 'Выберите верный вариант.', 'easy',
    '["Среднеквадратичную ошибку (MSE)", "Среднюю абсолютную ошибку (MAE)", "Логистическую функцию потерь"]',
    'Среднеквадратичную ошибку (MSE)', 'В классической постановке OLS минимизирует сумму квадратов ошибок (MSE).', NULL, NOW(), NOW()),
  (2001, 'Полиномиальная регрессия', 'Что фактически делает полиномиальная регрессия?', 'easy',
    '["Меняет модель", "Меняет пространство признаков", "Меняет функцию потерь"]',
    'Меняет пространство признаков', 'Мы добавляем полиномиальные признаки и обучаем линейную модель в расширенном пространстве.', NULL, NOW(), NOW()),
  (2002, 'Градиентный спуск', 'Какая величина определяет размер шага в градиентном спуске?', 'easy',
    '["Momentum", "Learning rate", "Batch size"]',
    'Learning rate', 'Скорость обучения (learning rate) задаёт величину шага по направлению антиградиента.', NULL, NOW(), NOW()),
  (2003, 'Cross‑validation', 'Зачем нужна кросс‑валидация?', 'medium',
    '["Оценить обобщающую способность", "Ускорить обучение", "Увеличить выборку"]',
    'Оценить обобщающую способность', 'Кросс‑валидация стабилизирует оценку качества за счёт усреднения по фолдам.', NULL, NOW(), NOW()),
  (2004, 'Data leakage', 'Пример утечки данных (data leakage) — это…', 'medium',
    '["Стандартизация после train/test split", "Стандартизация до train/test split", "Случайное разбиение"]',
    'Стандартизация до train/test split', 'Нормализацию нужно вычислять только по train и применять к val/test.', NULL, NOW(), NOW()),
  (2005, 'L2‑регуляризация', 'Что делает L2‑регуляризация с весами модели?', 'easy',
    '["Обнуляет многие веса", "Сжимает веса к нулю", "Усиливает большие веса"]',
    'Сжимает веса к нулю', 'L2 штрафует за большие значения весов, делая их меньше (но не обнуляя массово).', NULL, NOW(), NOW()),
  (2006, 'L1‑регуляризация', 'Чем известна L1‑регуляризация?', 'easy',
    '["Стимулирует разреженность", "Ускоряет инференс", "Увеличивает вес смещения"]',
    'Стимулирует разреженность', 'L1 приводит к нулевым коэффициентам у части признаков, отбирая важные.', NULL, NOW(), NOW()),
  (2007, 'Логистическая регрессия', 'Какую величину оценивает логистическая регрессия?', 'easy',
    '["Условную вероятность класса", "Среднеквадратичную ошибку", "AUC"]',
    'Условную вероятность класса', 'Модель оценивает P(y=1|x) через сигмоиду линейной комбинации признаков.', NULL, NOW(), NOW()),
  (2008, 'Precision vs Recall', 'Когда стоит предпочесть Recall над Precision?', 'medium',
    '["Когда важны все найденные положительные", "Когда критичны ложноположительные", "Когда классы сбалансированы"]',
    'Когда важны все найденные положительные', 'Например, в медицинском скрининге важнее не пропустить заболевших.', NULL, NOW(), NOW()),
  (2009, 'ROC‑AUC vs PR‑AUC', 'Почему PR‑AUC часто предпочтительнее при сильном дизбалансе?', 'medium',
    NULL,
    'PR‑AUC лучше отражает качество на положительном классе', 'При сильном дисбалансе ROC может быть завышен, PR‑AUC концентрируется на Precision‑Recall.', NULL, NOW(), NOW()),
  (2010, 'Калибровка вероятностей', 'Зачем нужна калибровка вероятностей?', 'medium',
    NULL,
    'Чтобы вероятности соответствовали реальности', 'Например, у 70% прогнозов с p≈0.7 примерно 70% должны быть верными.', NULL, NOW(), NOW()),
  (2011, 'SVM', 'Что контролирует параметр C в SVM?', 'easy',
    '["Ширину зазора и штраф за ошибки", "Тип ядра", "Число деревьев"]',
    'Ширину зазора и штраф за ошибки', 'Большой C стремится минимизировать ошибки на обучении, малый C — шире зазор.', NULL, NOW(), NOW()),
  (2012, 'Ядерный SVM', 'За счёт чего ядро помогает SVM?', 'medium',
    NULL,
    'Неявно отображает данные в более высокое пространство признаков', 'Kernel trick позволяет обучать линейный разделитель в скрытом пространстве.', NULL, NOW(), NOW()),
  (2013, 'Деревья решений', 'Главный риск одиночного дерева?', 'easy',
    '["Переобучение", "Недообучение", "Избыточная скорость"]',
    'Переобучение', 'Глубокие деревья легко переобучаются — поэтому их ограничивают/принят pruning.', NULL, NOW(), NOW()),
  (2014, 'Random Forest', 'Как бэггинг помогает снизить дисперсию?', 'medium',
    NULL,
    'Усредняет предсказания разных деревьев, уменьшая разброс', 'Случайные подвыборки и подмножества признаков повышают разнообразие деревьев.', NULL, NOW(), NOW()),
  (2015, 'Gradient Boosting', 'Идея градиентного бустинга — это…', 'medium',
    NULL,
    'Итеративное исправление ошибок предшествующих моделей', 'Каждый следующий слабый ученик обучается на остатках предыдущего.', NULL, NOW(), NOW()),
  (2016, 'Feature importance', 'SHAP/Permutation важности — это…', 'medium',
    NULL,
    'Методы интерпретации вклада признаков', 'Permutation важности измеряют падение качества при перемешивании признака; SHAP — аддитивные вклады.', NULL, NOW(), NOW()),
  (2017, 'k‑Means', 'Как выбрать k в k‑Means?', 'easy',
    '["Правилом локтя/силуэт", "По числу классов", "Случайно"]',
    'Правилом локтя/силуэт', 'Метрики помогают выбрать число кластеров по структуре данных.', NULL, NOW(), NOW()),
  (2018, 'DBSCAN', 'В чём преимущество DBSCAN над k‑Means?', 'medium',
    NULL,
    'Находит кластеры произвольной формы и “шум” без задания k', 'Использует eps и minPts, может выделять выбросы.', NULL, NOW(), NOW()),
  (2019, 'PCA', 'Что делает PCA?', 'easy',
    '["Сжимает данные по главным компонентам", "Увеличивает размерность", "Удаляет выбросы"]',
    'Сжимает данные по главным компонентам', 'Линейное понижение размерности вдоль направлений максимальной дисперсии.', NULL, NOW(), NOW());

-- Map questions to technology Machine Learning
INSERT INTO question_technologies (question_id, technology_id)
SELECT q.id, 3 FROM questions q WHERE q.id BETWEEN 2000 AND 2019
ON CONFLICT (question_id, technology_id) DO NOTHING;

-- Attach questions to modules (minimal, 2 per module)
INSERT INTO module_questions (module_id, question_id, question_order) VALUES
  (200, 2000, 1), (200, 2001, 2),
  (201, 2002, 1), (202, 2003, 1),
  (202, 2004, 2), (203, 2005, 1),
  (203, 2006, 2), (204, 2007, 1),
  (204, 2008, 2), (205, 2009, 1),
  (205, 2010, 2), (206, 2011, 1),
  (206, 2012, 2), (207, 2013, 1),
  (207, 2014, 2), (208, 2015, 1),
  (208, 2016, 2), (209, 2017, 1),
  (209, 2018, 2), (209, 2019, 3)
ON CONFLICT (module_id, question_id) DO NOTHING;

COMMIT;

-- +goose Down
BEGIN;
  -- detach module_questions
  DELETE FROM module_questions WHERE module_id BETWEEN 200 AND 209;
  -- delete questions
  DELETE FROM question_technologies WHERE question_id BETWEEN 2000 AND 2019;
  DELETE FROM questions WHERE id BETWEEN 2000 AND 2019;
  -- delete modules and course
  DELETE FROM course_technologies WHERE course_id = 2;
  DELETE FROM modules WHERE course_id = 2;
  DELETE FROM courses WHERE id = 2;
  -- delete technology if unused
  DELETE FROM technologies WHERE id = 3 AND NOT EXISTS (
    SELECT 1 FROM course_technologies ct WHERE ct.technology_id = 3
  ) AND NOT EXISTS (
    SELECT 1 FROM question_technologies qt WHERE qt.technology_id = 3
  );
COMMIT;
