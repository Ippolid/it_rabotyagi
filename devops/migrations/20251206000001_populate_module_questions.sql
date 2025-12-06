-- +goose Up
BEGIN;

-- Migration: Populate modules with questions
-- Description: Add existing questions to course modules
-- Date: 2025-12-06

-- Add questions to Module #1: "Введение в Go"
-- This module covers Go basics, so we add fundamental questions
INSERT INTO module_questions (module_id, question_id, question_order) VALUES
  (1, 5, 1),  -- Какое значение по умолчанию имеет bool?
  (1, 1, 2),  -- Что такое слайс (slice) в Go?
  (1, 3, 3)   -- Когда выполняется функция, объявленная с defer?
ON CONFLICT (module_id, question_id) DO NOTHING;

-- Add questions to Module #3: "Конкурентность в Go"
-- This module covers concurrency, so we add goroutines and channels questions
INSERT INTO module_questions (module_id, question_id, question_order) VALUES
  (3, 2, 1),  -- С помощью какого ключевого слова запускается горутина?
  (3, 4, 2)   -- Что произойдет при отправке в nil-канал?
ON CONFLICT (module_id, question_id) DO NOTHING;

-- Note: Module #2 "Основы синтаксиса" already has questions from previous setup

COMMIT;

-- +goose Down
BEGIN;

-- Remove questions from Module #1
DELETE FROM module_questions
WHERE module_id = 1 AND question_id IN (5, 1, 3);

-- Remove questions from Module #3
DELETE FROM module_questions
WHERE module_id = 3 AND question_id IN (2, 4);

COMMIT;
