# TZshka

Система проверки технического задания (ТЗ) с использованием LLM.

## Возможности

- **Авторизация** — регистрация и вход по логину и паролю
- **Управление сессиями** — создание, просмотр, присоединение и удаление сессий
- **Проверка ТЗ** — отправка текста или файла на проверку по ГОСТ-19/ГОСТ-34
- **История** — сохранение и просмотр истории проверок

## Быстрый старт

### Требования

- Docker и Docker Compose
- OpenRouter API ключ (для LLM)

### Настройка

1. Создайте файл `.env` (скопируйте из `.env.example` если есть):
   ```
   OPENROUTER_API_KEY=ваш-ключ
   JWT_SECRET=ваш-секретный-ключ
   ```

### Запуск

```bash
# Первый запуск или сброс данных
docker-compose down -v
docker-compose up --build

# Последующие запуски
docker-compose up
```

### Сервисы

| Сервис | Порт | Описание |
|--------|------|-----------|
| app | 8080 | Go backend API |
| llm_api | 8001 | Python LLM сервис |
| frontend | 5173 | React frontend |
| postgres | 5432 | База данных |

## API Эндпоинты

Все эндпоинты Go backend начинаются с `/api`. Токен авторизации передаётся в заголовке `Authorization: Bearer <token>`.

### Авторизация

| Метод | Эндпоинт | Описание | Авторизация | Body |
|-------|----------|----------|-------------|------|
| POST | `/register` | Регистрация нового пользователя | Нет | `{ "login": string, "email": string, "password": string }` |
| POST | `/login` | Вход в систему | Нет | `{ "login": string, "password": string }` |

**Ответ (успех):**
- Register: `201` — `{ "user": { "id", "login", "email", "createdAt" } }`
- Login: `200` — `{ "token": string }`

**Ответ (ошибка):**
- `400` — `{ "error": { "code": string, "message": string } }`
- `500` — `{ "error": { "code": "INTERNAL_ERROR", "message": "internal server error" } }`

---

### Сессии

| Метод | Эндпоинт | Описание | Body / Params |
|-------|----------|----------|---------------|
| POST | `/sessions/create` | Создать сессию | `{ "name": string }` |
| GET | `/sessions` | Получить список своих сессий | — |
| DELETE | `/sessions` | Удалить сессию | `{ "sessionId": string }` |
| POST | `/sessions/join` | Присоединиться к сессии | `{ "sessionId": string }` |

**Ответ (успех):**
- Create: `201` — `{ "session": { "id", "name", "creatorId", "createdAt" } }`
- Get: `200` — `{ "sessions": [...] }`
- Join: `200` — `{ "session": { ... } }`
- Delete: `200` — `{ "message": "session deleted" }`

**Ответ (ошибка):**
- `401` — `{ "error": { "code": "UNAUTHORIZED", "message": "invalid token" } }`
- `400` — `{ "error": { "code": "INVALID_REQUEST", "message": "..." } }`
- `404` — `{ "error": { "code": "NOT_FOUND", "message": "session not found" } }`

---

### LLM — Проверка ТЗ

| Метод | Эндпоинт | Описание | Body / Params |
|-------|----------|----------|---------------|
| POST | `/llm/text` | Проверка текста | `mode` (string, form), `standard` (string, form, опц.), `content` (string, form), `sessionId` (string, form, опц.) |
| POST | `/llm/file` | Проверка файла (.txt) | `mode` (form), `standard` (form, опц.), `file` (file, form), `sessionId` (form, опц.) |

**Параметры:**
- `mode` — режим проверки: `"Instant"` или `"Thinking"`
- `standard` — стандарт: `"ГОСТ-19"`, `"ГОСТ-34"` или пустое значение (только проверка правил)
- `content` — текст ТЗ (для `/text`)
- `file` — `.txt` файл (для `/file`)
- `sessionId` — ID сессии для сохранения в историю (опционально)

**Ответ (успех):**
```json
{
  "success": true,
  "code": 200,
  "data": {
    "rules_checker": {
      "status": "valid | issues_found",
      "feedback": "string",
      "issues": [
        { "rule_id": "string", "problem": "string", "explanation": "string" }
      ]
    },
    "standard_checker": {
      "status": "valid | issues_found",
      "feedback": "string",
      "issues": [
        { "rule_id": "string", "problem": "string", "explanation": "string" }
      ]
    }
  }
}
```

**Ответ (ошибка):**
- `400` — `{ "error": { "type": "validation_error", "message": "..." } }`
- `500` — `{ "error": { "type": "internal_error", "message": "..." } }`

---

### История проверок

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| GET | `/history/:id` | Получить историю проверок сессии (`:id` — sessionId) |

**Ответ (успех):**
```json
{
  "corrections": [
    {
      "id": "uuid",
      "sessionId": "uuid",
      "userId": "uuid",
      "inputContent": "string",
      "responseData": { "data": {...}, "success": true, "code": 200 },
      "createdAt": "timestamp"
    }
  ]
}
```

---

### Внутренний LLM API (Python)

Эндпоинты Python-сервиса (`localhost:8001`):

| Метод | Эндпоинт | Описание | Body / Params |
|-------|----------|----------|---------------|
| GET | `/` | Health check | — |
| POST | `/llm/text` | Получить анализ от LLM | form: `mode`, `standard`, `content` |

---

### База данных (PostgreSQL)

| Таблица | Описание |
|---------|----------|
| `users` | Пользователи (`id`, `login`, `email`, `password_hash`, `created_at`) |
| `sessions` | Сессии (`id`, `name`, `creator_id`, `created_at`) |
| `corrections_history` | История проверок (`id`, `session_id`, `user_id`, `input_content`, `response_data` JSONB, `created_at`) |

---

## Структура

- `backend/` — Go backend (Gin, pgxpool, JWT)
- `frontend/` — React frontend (Vite)
- `llm_api/` — Python LLM сервис (FastAPI, LangChain, OpenRouter)
- `docker-compose.yml` — Запуск всех сервисов