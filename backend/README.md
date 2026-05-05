# TZshka

Backend проект с авторизацией, управлением сессиями и историей правок.

## Технологии

- Go 1.25
- PostgreSQL
- Gin (HTTP фреймворк)
- JWT (аутентификация)
- pgx (работа с PostgreSQL)

## Запуск

### Docker

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
| Email | string | Email пользователя |
| Role | string | Роль пользователя (admin/user) |
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
  "email": "user@example.com",
  "password": "password123",
  "role": "user"
}
```

**Успешный ответ (201):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "role": "user"
  }
}
```

**Ошибка - пользователь уже существует (400):**
```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "email already exists"
  }
}
```

#### Вход

```bash
POST /api/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Успешный ответ (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ошибка - неверные credentials (401):**
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "invalid credentials"
  }
}
```

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
```

---

### History

#### Получение истории сессии

```bash
GET /api/history/:id
```

**Успешный ответ (200):**
```json
{
  "corrections": [
    {
      "id": "uuid",
      "sessionId": "uuid",
      "inputContent": "Текст ТЗ",
      "responseData": {},
      "createdAt": "2026-04-27T12:00:00Z"
    }
  ]
}
```

**Ошибка - невалидный id (400):**
```json
{
  "error": "invalid session id"
}
```

---

## Структура проекта

```
backend/
├── cmd/
│   └── main.go              # Точка входа
├── config/
│   └── config.yaml          # Конфигурация
├── docker/
│   └── Dockerfile           # Docker образ
├── tests/
│   └── auth_test.go         # Тесты
├── internal/
│   ├── auth/
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
├── pkg/
│   ├── logger/             # Логгер
│   ├── postgres/            # Подключение к PostgreSQL
│   └── registr/             # JWT токены
├── docker-compose.yml       # Docker Compose
├── go.mod
└── README.md
```