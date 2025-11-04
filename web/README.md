# IT-RABOTYAGI Web (Frontend)

Отдельный фронтенд-проект на React + TypeScript + Vite + Tailwind. Бэкенд (Go + Echo) остаётся нетронутым.

## Локальный запуск

1. Установите зависимости:
```bash
cd web
npm install
```

2. Запустите dev-сервер:
```bash
npm run dev
```

- Приложение откроется на `http://localhost:5173`
- Все запросы на `/api` проксируются на `http://localhost:8080` (см. `vite.config.ts`).

## Сборка
```bash
npm run build
npm run preview
```

## Структура
- `src/App.tsx` — главная страница (Hero, Features, Mentors, Courses, CTA)
- `src/styles.css` — Tailwind стили и утилиты

## Примечание
- Чтобы раздавать собранную статическую версию через nginx, можно добавить отдельный сервис в docker-compose и монтировать `web/dist`.





