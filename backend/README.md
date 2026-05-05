# TZshka

<<<<<<< HEAD
Backend проект с авторизацией, управлением сессиями и историей правок.
=======
Backend проект с авторизацией, регистрацией, управлением сессиями и LLM интеграцией.
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e

## Технологии

- Go 1.25
- PostgreSQL
- Gin (HTTP фреймворк)
- JWT (аутентификация)
- pgx (работа с PostgreSQL)

## Запуск

<<<<<<< HEAD
### Docker
=======
### Docker (рекомендуется)
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e

```bash
# Первый запуск или сброс данных
docker-compose down -v
docker-compose up --build

# Последующие запуски
docker-compose up
```

### Остановка

```bash
docker-compose down
```

## Тесты

```bash
go test -v ./tests/...
```

## Сущности

### User

| Поле | Тип | Описание |
|------|-----|----------|
| ID | UUID | Уникальный идентификатор |
| Login | string | Логин пользователя |
| Email | string | Email пользователя |
| CreatedAt | timestamp | Дата создания |

### Session

| Поле | Тип | Описание |
|------|-----|----------|
| ID | UUID | Уникальный идентификатор |
| Name | string | Название сессии |
| CreatorID | UUID | ID создателя |
| CreatedAt | timestamp | Дата создания |

### CorrectionsHistory

| Поле | Тип | Описание |
|------|-----|----------|
| ID | UUID | Уникальный идентификатор |
| SessionID | UUID | ID сессии |
| InputContent | text | Входной текст |
| ResponseData | JSONB | Ответ от LLM |
| CreatedAt | timestamp | Дата создания |

### Session

| Поле | Тип | Описание |
|------|-----|----------|
| ID | UUID | Уникальный идентификатор |
| Name | string | Название сессии |
| CreatorID | UUID | ID создателя |
| CreatedAt | timestamp | Дата создания |

### CorrectionsHistory

| Поле | Тип | Описание |
|------|-----|----------|
| ID | UUID | Уникальный идентификатор |
| SessionID | UUID | ID сессии |
| InputContent | text | Входной текст |
| ResponseData | JSONB | Ответ от LLM |
| CreatedAt | timestamp | Дата создания |

---

## API

### Auth

#### Регистрация

```bash
POST /api/register
Content-Type: application/json

{
  "login": "user1",
  "email": "user1@test.com",
  "password": "password123"
}
```

**Ответ (201):**
```json
{"user": {"id": "uuid", "login": "user1", "email": "user1@test.com"}}
```

<<<<<<< HEAD
**Ошибка - пользователь уже существует (400):**
```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "email already exists"
  }
}
```

=======
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
#### Вход

```bash
POST /api/login
Content-Type: application/json

{
  "login": "user1",
  "password": "password123"
}
```

**Ответ (200):**
```json
{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
```

---

### Session

#### Создание сессии

```bash
POST /api/sessions/create
Authorization: Bearer <token>
Content-Type: application/json

{"name": "My Session"}
```

**Ответ (201):**
```json
{"session": {"id": "uuid", "name": "My Session", "creatorId": "uuid"}}
```

#### Получение сессий пользователя

```bash
GET /api/sessions
Authorization: Bearer <token>
```

**Ответ (200):**
```json
{"sessions": [{"id": "uuid", "name": "My Session", "creatorId": "uuid"}]}
```

#### Вход в сессию

```bash
POST /api/sessions/join
Content-Type: application/json

{"sessionId": "uuid"}
```

#### Удаление сессии

```bash
DELETE /api/sessions
Authorization: Bearer <token>
Content-Type: application/json

{"sessionId": "uuid"}
```

---

### LLM

#### Обработка текста (JSON)

```bash
POST /api/llm/text
Content-Type: application/json

{
  "mode": "Instant",       # или "Thinking"
  "standard": "ГОСТ-19",  # или "ГОСТ-34", или пусто
  "content": "Текст технического задания",
  "sessionId": "uuid"     # опционально, для сохранения в историю
}
```

**Ответ (200):**
```json
{
  "success": true,
  "code": 200,
  "data": {
    "rules_checker": "AI вывод",
    "standard_checker": "AI вывод по стандарту"
  }
}
```

<<<<<<< HEAD
---

### Session

#### Создание сессии

```bash
POST /api/sessions/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "My Session"
}
```

**Успешный ответ (201):**
```json
{
  "session": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "My Session",
    "creatorId": "00000000-0000-0000-0000-000000000002"
  }
}
```

**Ошибка - без токена (401):**
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "invalid token"
  }
}
```

#### Получение сессий пользователя

```bash
GET /api/sessions
Authorization: Bearer <token>
```

**Успешный ответ (200):**
```json
{
  "sessions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "My Session",
      "creatorId": "00000000-0000-0000-0000-000000000002"
    }
  ]
}
```

#### Вход в сессию

```bash
POST /api/sessions/join
Content-Type: application/json

{
  "sessionId": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Успешный ответ (200):**
```json
{
  "session": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "My Session",
    "creatorId": "00000000-0000-0000-0000-000000000002"
  }
}
```

**Ошибка - сессия не найдена (404):**
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "session not found"
  }
}
```

#### Удаление сессии

```bash
DELETE /api/sessions
Authorization: Bearer <token>
Content-Type: application/json

{
  "sessionId": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Успешный ответ (200):**
```json
{
  "message": "session deleted"
}
=======
#### Обработка файла

```bash
POST /api/llm/file
Content-Type: multipart/form-data

# Поля формы:
# - mode: "Instant" или "Thinking"
# - standard: "ГОСТ-19", "ГОСТ-34" или пусто
# - file: текстовый файл (.txt)
# - sessionId: UUID сессии (опционально)
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
```

---

### History

#### Получение истории сессии

```bash
GET /api/history/:id
<<<<<<< HEAD
```

**Успешный ответ (200):**
=======
Authorization: Bearer <token>
```

**Ответ (200):**
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
```json
{
  "corrections": [
    {
      "id": "uuid",
      "sessionId": "uuid",
      "inputContent": "Текст ТЗ",
<<<<<<< HEAD
      "responseData": {},
=======
      "responseData": {...},
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
      "createdAt": "2026-04-27T12:00:00Z"
    }
  ]
}
```

<<<<<<< HEAD
**Ошибка - невалидный id (400):**
```json
{
  "error": "invalid session id"
}
```

---

=======
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
## Структура проекта

```
backend/
├── cmd/main.go
├── config/config.yaml
├── docker/Dockerfile
├── tests/
│   ├── auth_test.go
│   └── history_test.go
├── internal/
│   ├── auth/
<<<<<<< HEAD
│   │   ├── models/          # DTO
│   │   ├── repository/      # Работа с БД
│   │   ├── route/           # HTTP хендлеры
│   │   └── service/         # Бизнес-логика
│   ├── session/
│   │   ├── models/          # DTO
│   │   ├── repository/      # Работа с БД
│   │   ├── route/           # HTTP хендлеры
│   │   └── service/        # Бизнес-логика
│   ├── history/
│   │   ├── models/         # DTO
│   │   ├── repository/     # Работа с БД
│   │   ├── route/          # HTTP хендлеры
│   │   └── service/        # Бизнес-логика
│   ├── config/              # Загрузка конфига
│   └── migrations/          # Миграции БД
=======
│   ├── session/
│   ├── llm/
│   ├── history/
│   ├── config/
│   └── migrations/
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
├── pkg/
│   ├── logger/
│   ├── postgres/
│   └── registr/
├── go.mod
└── README.md
```