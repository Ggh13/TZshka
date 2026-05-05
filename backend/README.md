# TZshka

Backend проект с авторизацией, регистрацией, управлением сессиями и LLM интеграцией.

## Технологии

- Go 1.25
- PostgreSQL
- Gin (HTTP фреймворк)
- JWT (аутентификация)
- pgx (работа с PostgreSQL)

## Запуск

### Docker (рекомендуется)

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

#### Обработка файла

```bash
POST /api/llm/file
Content-Type: multipart/form-data

# Поля формы:
# - mode: "Instant" или "Thinking"
# - standard: "ГОСТ-19", "ГОСТ-34" или пусто
# - file: текстовый файл (.txt)
# - sessionId: UUID сессии (опционально)
```

---

### History

#### Получение истории сессии

```bash
GET /api/history/:id
Authorization: Bearer <token>
```

**Ответ (200):**
```json
{
  "corrections": [
    {
      "id": "uuid",
      "sessionId": "uuid",
      "inputContent": "Текст ТЗ",
      "responseData": {...},
      "createdAt": "2026-04-27T12:00:00Z"
    }
  ]
}
```

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
│   ├── session/
│   ├── llm/
│   ├── history/
│   ├── config/
│   └── migrations/
├── pkg/
│   ├── logger/
│   ├── postgres/
│   └── registr/
├── go.mod
└── README.md
```