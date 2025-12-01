# Interview Trainer Service

Сервис для тренировки прохождения собеседований с использованием AI.

## Возможности

- **Транскрипция аудио**: Whisper API от Groq для преобразования речи в текст
- **Оценка ответов**: LLM (Llama 3.3 70B) для анализа и оценки ответов
- **Фидбек**: Детальная обратная связь с оценкой, сильными сторонами и рекомендациями

## API Endpoints

### 1. Health Check
```
GET /
```

### 2. Транскрипция аудио
```
POST /api/v1/transcribe
Content-Type: multipart/form-data

Parameters:
- audio: Audio file (mp3, wav, ogg, webm, m4a)
- language: Optional language code (ru, en, etc.)

Response:
{
  "text": "Транскрибированный текст",
  "language": "ru",
  "duration": 15.5
}
```

### 3. Оценка ответа
```
POST /api/v1/evaluate
Content-Type: application/json

Body:
{
  "question": "Вопрос собеседования",
  "ideal_answer": "Идеальный ответ",
  "user_answer": "Ответ пользователя",
  "language": "ru"
}

Response:
{
  "score": 85,
  "feedback": "Полный текст фидбека",
  "strengths": [
    "Хорошо объяснил концепцию",
    "Привел примеры из практики"
  ],
  "improvements": [
    "Можно добавить больше деталей о...",
    "Стоит упомянуть..."
  ],
  "overall_comment": "Общий комментарий"
}
```

### 4. Полный процесс (транскрипция + оценка)
```
POST /api/v1/interview/complete
Content-Type: multipart/form-data

Parameters:
- audio: Audio file
- question: Interview question
- ideal_answer: Expected answer
- language: Optional language code

Response:
{
  "transcription": { ... },
  "evaluation": { ... }
}
```

## Запуск локально

### Требования
- Python 3.11+
- Groq API key

### Установка
```bash
cd ml
pip install -r requirements.txt
```

### Запуск
```bash
export GROQ_API_KEY="your_groq_api_key"
python interview_trainer.py
```

Сервис будет доступен на `http://localhost:8001`

## Запуск в Docker

```bash
cd devops/docker
docker compose up ml-service
```

## Переменные окружения

- `GROQ_API_KEY` (обязательно) - API ключ от Groq

## Интеграция с Backend

Backend должен:
1. Получить аудио файл от пользователя
2. Отправить на `/api/v1/interview/complete` с вопросом и идеальным ответом из БД
3. Сохранить результат в таблицу `interview_attempts`

Пример запроса из Go:
```go
// Формируем multipart request
body := &bytes.Buffer{}
writer := multipart.NewWriter(body)

// Добавляем аудио файл
part, _ := writer.CreateFormFile("audio", "answer.webm")
io.Copy(part, audioFile)

// Добавляем текстовые поля
writer.WriteField("question", question)
writer.WriteField("ideal_answer", idealAnswer)
writer.WriteField("language", "ru")
writer.Close()

// Отправляем запрос
resp, err := http.Post(
    "http://ml-service:8001/api/v1/interview/complete",
    writer.FormDataContentType(),
    body,
)
```

## Costs

При использовании Groq:
- Whisper: **БЕСПЛАТНО** (лимит в бесплатном tier)
- LLM (Llama 3.3 70B): **БЕСПЛАТНО** (лимит ~14,400 запросов/день)

Примерные затраты на одно интервью: **$0** в пределах лимитов

## Модели

- **Whisper**: `whisper-large-v3-turbo` - быстрая и точная транскрипция
- **LLM**: `llama-3.3-70b-versatile` - качественная оценка и фидбек

## Swagger документация

После запуска доступна по адресу: `http://localhost:8001/docs`
