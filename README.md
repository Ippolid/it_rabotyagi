# IT-RABOTYAGI - Платформа для обучения и менторинга в IT

## 📋 Описание проекта

IT-RABOTYAGI - это платформа для обучения программированию с возможностью менторинга. Проект включает в себя систему курсов, базу вопросов для практики и поиск менторов.

### Основные функции:
- 🔐 Аутентификация и авторизация с JWT
- 👤 Управление профилями пользователей
- 📚 База вопросов для практики
- 🎓 Система курсов и модулей
- 👨‍🏫 Поиск менторов
- 💬 API документация через Swagger

## 🛠 Технологии

- **Backend**: Go 1.24+
- **Framework**: Echo
- **Database**: PostgreSQL 17
- **Миграции**: Goose
- **API**: OpenAPI 3.0 / Swagger
- **Контейнеризация**: Docker, Docker Compose
- **JWT**: golang-jwt/jwt/v5

## 🚀 Быстрый старт

### Требования
- Docker и Docker Compose
- Make (опционально)
- Go 1.24+ (для локальной разработки)

### 1. Клонирование проекта

```bash
git clone <repository-url>
cd it_rabotyagi
```

### 2. Настройка окружения

```bash
# Переход в папку devops
cd devops/docker

# Создание .env файла из примера
cp .env.example .env

# Отредактируйте .env файл, изменив значения по умолчанию
# на свои собственные (пароли, секретные ключи и т.д.)
```

### 3. Запуск через Docker Compose

```bash
# Запуск всех сервисов (PostgreSQL, миграции, приложение)
docker compose up --build

# Запуск в фоновом режиме
docker compose up -d --build
```

Приложение будет доступно по адресу:
- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/api/openapi/index.html

### 4. Локальная разработка

```bash
# Переход в папку backend
cd backend

# Установка зависимостей
make install-deps

# Генерация кода из OpenAPI
make openapi

# Запуск миграций
make local-migration-up

# Сборка и запуск
go build -o app ./cmd/app/main.go
./app
```

### 5. Остановка сервисов

```bash
cd devops/docker

# Остановка
docker compose down

# Остановка с удалением volumes (очистка БД)
docker compose down -v
```

## 📁 Структура проекта

```
it_rabotyagi/
├── backend/                  # Backend приложение на Go
│   ├── api/openapi/          # OpenAPI спецификация
│   ├── cmd/app/              # Точка входа
│   ├── internal/             # Внутренний код
│   │   ├── app/              # Инициализация
│   │   ├── business/         # Бизнес-логика
│   │   ├── config/           # Конфигурация
│   │   ├── data/             # Слой данных
│   │   ├── logger/           # Логирование
│   │   └── server/           # HTTP сервер
│   ├── go.mod
│   ├── Makefile
│   └── README.md
├── devops/                   # DevOps конфигурация
│   ├── docker/               # Docker файлы
│   │   ├── Dockerfile
│   │   ├── migration.Dockerfile
│   │   ├── docker-compose.yaml
│   │   ├── migration.sh
│   │   └── .env.example
│   ├── migrations/           # SQL миграции
│   └── README.md
├── frontend/                 # Frontend (будущая разработка)
│   └── README.md
└── README.md
```

## 🔧 Конфигурация

Приложение настраивается через переменные окружения в `devops/docker/.env`.

**Важно**: Создайте `.env` файл из `.env.example` шаблона:
```bash
cd devops/docker
cp .env.example .env
```

Затем отредактируйте `.env`, установив свои значения для:
- `PG_USER`, `PG_PASSWORD` - учетные данные PostgreSQL
- `JWT_SECRET` - секретный ключ для JWT токенов
- Другие параметры по необходимости

Пример переменных окружения:

```yaml
environment:
  DB_URL: "postgres://user:password@postgres:5432/it_rabotyagi?sslmode=disable"
  SERVER_HOST: "0.0.0.0"
  SERVER_PORT: "8080"
  AUTH_SECRET: "your-secret-key"
  AUTH_TOKEN_DURATION: "60"        # минуты
  AUTH_REFRESH_DURATION: "10080"   # минуты (7 дней)
  LOG_LEVEL: "info"
```

## 📡 API Endpoints

### Аутентификация

#### POST `/api/v1/auth/register`
Регистрация нового пользователя

**Request:**
```json
{
  "email": "user@example.com",
  "nickname": "username",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 3600
}
```

#### POST `/api/v1/auth/login`
Вход пользователя

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 3600
}
```

#### POST `/api/v1/auth/refresh`
Обновление access токена

**Request:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 3600
}
```

### Защищенные endpoints (требуют Authorization header)

#### GET `/api/v1/users/me`
Получение профиля текущего пользователя

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "fullName": "Username",
  "role": "user"
}
```

#### GET `/api/v1/mentors`
Список менторов (опциональная авторизация)

**Response (200):**
```json
{
  "items": [
    {
      "id": 1,
      "fullName": "Alice Smith",
      "title": "Senior Go Developer",
      "skills": ["Go", "Docker", "Kubernetes"],
      "yearsOfExperience": 5
    }
  ],
  "total": 1
}
```

## 🧪 Тестирование API

### Через Swagger UI
1. Откройте http://localhost:8080/api-docs/index.html
2. Используйте интерфейс Swagger для тестирования endpoints

### Через curl

**Регистрация:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "nickname": "testuser",
    "password": "password123"
  }'
```

**Логин:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Получение профиля:**
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <YOUR_ACCESS_TOKEN>"
```

## 🗃 База данных

### Схема основных таблиц

**users** - Пользователи
- id, username, email, password
- avatar_url, role, description
- created_at, updated_at

**auth_sessions** - Сессии пользователей
- id, user_id, refresh_token_hash
- expires_at, created_at, revoked_at

**courses** - Курсы
- id, title, description, is_published

**questions** - Вопросы для практики
- id, title, content, difficulty
- options (JSONB), correct_answer, explanation

**mentors** - Менторы
- id, user_id, specialization, grade
- experience_years, description, tags

### Применение миграций вручную

```bash
# Войти в контейнер с БД
docker exec -it it_rabotyagi_postgres psql -U user -d it_rabotyagi

# Проверить текущую версию миграций
\dt goose_db_version

# Выполнить SQL команды
# Например, обновить роль для существующих пользователей:
UPDATE users SET role = 'user' WHERE role = '' OR role IS NULL;
```

## 🔐 Безопасность

### Хранение паролей
- Пароли хешируются с использованием SHA256
- Рекомендация: мигрировать на bcrypt для большей безопасности

### JWT токены
- **Access Token**: срок жизни 60 минут
- **Refresh Token**: срок жизни 7 дней
- Refresh токены хранятся в БД в виде SHA256 хеша

### Middleware
- `AuthMiddleware` - обязательная проверка токена
- `OptionalAuthMiddleware` - опциональная проверка токена
- `RoleMiddleware` - проверка роли пользователя

## 🛠 Разработка

### Установка зависимостей

```bash
go mod download
```

### Генерация кода из OpenAPI

```bash
make generate
# или
oapi-codegen -package openapi -generate types,server,spec \
  api/openapi/openapi.yaml > api/openapi/server.gen.go
```

### Локальный запуск (без Docker)

```bash
# Запустить только PostgreSQL
docker-compose up postgres -d

# Применить миграции
export DB_URL="postgres://user:password@localhost:5432/it_rabotyagi?sslmode=disable"
./bin/goose -dir migrations up

# Запустить приложение
go run cmd/app/main.go
```

### Создание новой миграции

```bash
./bin/goose -dir migrations create migration_name sql
```

### Сборка бинарника

```bash
make build
# или
go build -o bin/app ./cmd/app
```

## 📊 Логирование

Приложение использует структурированное логирование с zap:
- **info** - основные операции
- **error** - ошибки
- **debug** - детальная отладочная информация

Логи выводятся в stdout и могут быть просмотрены через:
```bash
docker-compose logs -f app
```

## 🐛 Troubleshooting

### Проблема: Контейнер БД не запускается
```bash
docker-compose down -v
docker-compose up --build
```

### Проблема: Миграции не применяются
```bash
docker-compose down
docker volume rm it_rabotyagi_postgres-data
docker-compose up --build
```

### Проблема: Порт 8080 занят
```bash
# Изменить порт в docker-compose.yaml
ports:
  - "8081:8080"  # внешний:внутренний
```

### Проблема: Роль пользователя пустая
```bash
# Обновить роли в БД
docker exec it_rabotyagi_postgres psql -U user -d it_rabotyagi \
  -c "UPDATE users SET role = 'user' WHERE role = '' OR role IS NULL;"
```

## 📚 Дополнительная документация

- [AUTH_GUIDE.md](AUTH_GUIDE.md) - Подробная документация по системе аутентификации
- [SWAGGER_GUIDE.md](SWAGGER_GUIDE.md) - Работа со Swagger
- [api/openapi/openapi.yaml](api/openapi/openapi.yaml) - OpenAPI спецификация

## 👥 Команда

Проект IT-RABOTYAGI

## 📝 Лицензия

MIT License

## 🤝 Контрибьютинг

1. Fork проекта
2. Создайте feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit изменения (`git commit -m 'Add some AmazingFeature'`)
4. Push в branch (`git push origin feature/AmazingFeature`)
5. Откройте Pull Request

---

**Для проверяющего**:
- Все endpoints доступны через Swagger UI по адресу http://localhost:8080/api-docs/index.html
- Тестовые данные автоматически загружаются через миграции
- Для быстрого теста: `docker-compose up --build` и откройте Swagger UI

