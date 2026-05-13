# TZshka

Backend проект с авторизацией, управлением сессиями, LLM интеграцией и историей правок.

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
| Password | string | Хэш пароля |
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
| UserID | UUID | ID пользователя |
| InputContent | text | Входной текст |
| ResponseData | JSONB | Ответ от LLM |
| CreatedAt | timestamp | Дата создания |

---

## API

> **Авторизация:** эндпоинты с пометкой 🔒 требуют заголовок `Authorization: Bearer <token>`.
> Токен выдаётся при входе (`POST /api/login`).

### Auth

#### Регистрация

```
POST /api/register
Content-Type: application/json
```

**Body:**

| Поле | Тип | Обязательное | Описание |
|------|-----|:----------:|----------|
| login | string | ✅ | Логин пользователя (мин. 3 символа) |
| email | string | ✅ | Email пользователя |
| password | string | ✅ | Пароль (мин. 6 символов) |

```json
{
  "login": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**Успешный ответ (201):**

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "login": "testuser",
    "email": "test@example.com",
    "createdAt": "2026-04-27T12:00:00Z"
  }
}
```

**Ошибка — пользователь уже существует (400):**

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "login already exists"
  }
}
```

#### Вход 🔒

```
POST /api/login
Content-Type: application/json
```

**Body:**

| Поле | Тип | Обязательное | Описание |
|------|-----|:----------:|----------|
| login | string | ✅ | Логин пользователя |
| password | string | ✅ | Пароль |

```json
{
  "login": "testuser",
  "password": "password123"
}
```

**Успешный ответ (200):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ошибка — неверные credentials (401):**

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "invalid credentials"
  }
}
```

---

### Session 🔒

Все эндпоинты сессий требуют заголовок `Authorization: Bearer <token>`.

#### Создание сессии

```
POST /api/sessions/create
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**

| Поле | Тип | Обязательное | Описание |
|------|-----|:----------:|----------|
| name | string | ✅ | Название сессии (мин. 1 символ) |

```json
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
    "creatorId": "00000000-0000-0000-0000-000000000002",
    "createdAt": "2026-04-27T12:00:00Z"
  }
}
```

#### Получение сессий пользователя

```
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
      "creatorId": "00000000-0000-0000-0000-000000000002",
      "createdAt": "2026-04-27T12:00:00Z"
    }
  ]
}
```

#### Вход в сессию

```
POST /api/sessions/join
Content-Type: application/json
```

**Body:**

| Поле | Тип | Обязательное | Описание |
|------|-----|:----------:|----------|
| sessionId | string | ✅ | UUID сессии |

```json
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

**Ошибка — сессия не найдена (404):**

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "session not found"
  }
}
```

#### Удаление сессии

```
DELETE /api/sessions
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**

| Поле | Тип | Обязательное | Описание |
|------|-----|:----------:|----------|
| sessionId | string | ✅ | UUID сессии |

```json
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

### LLM 🔒

Эндпоинты LLM требуют `Authorization: Bearer <token>` для сохранения результатов в историю. Без токена запрос будет обработан, но в историю сохранён не будет.

#### Обработка текста

```
POST /api/llm/text
Content-Type: application/json
```

**Body:**

| Поле | Тип | Обязательное | Описание |
|------|-----|:----------:|----------|
| mode | string | ✅ | Режим: `"Instant"` или `"Thinking"` |
| content | string | ✅ | Текст технического задания |
| standard | string | ❌ | Стандарт: `"ГОСТ-19"`, `"ГОСТ-34"` или пустая строка |
| sessionId | string | ❌ | UUID сессии для сохранения в историю |

```json
{
  "mode": "Instant",
  "standard": "ГОСТ-19",
  "content": "Текст технического задания",
  "sessionId": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Успешный ответ (200):**

```json
{
  "success": true,
  "code": 200,
  "data": {
    "rules_checker": {
      "status": "issues_found",
      "feedback": "Нашел несколько мест, которые могут трактоваться неоднозначно.",
      "issues": [
        {
          "rule_id": "R1",
          "problem": "Фраза ... слишком общая и не задает четкий измеримый результат.",
          "explanation": "Система или исполнитель не смогут однозначно понять..."
        }
      ]
    },
    "standard_checker": {
      "status": "issues_found",
      "feedback": "Проверка по ГОСТ-19 выявила неточности оформления и структуры.",
      "issues": [
        {
          "rule_id": "ГОСТ-19",
          "problem": "Структура требования выглядит неполной для выбранного стандарта.",
          "explanation": "Желательно разделить цель, функциональные требования..."
        }
      ]
    }
  }
}
```

**Ошибка (500):**

```json
{
  "success": false,
  "code": 500,
  "data": null,
  "error": {
    "type": "internal_error",
    "message": "LLM service unavailable"
  }
}
```

#### Обработка файла

```
POST /api/llm/file
Content-Type: multipart/form-data
```

**Поля формы:**

| Поле | Тип | Обязательное | Описание |
|------|-----|:----------:|----------|
| mode | string | ✅ | Режим: `"Instant"` или `"Thinking"` |
| file | file | ✅ | Текстовый файл (.txt) |
| standard | string | ❌ | Стандарт: `"ГОСТ-19"`, `"ГОСТ-34"` или пустая строка |
| sessionId | string | ❌ | UUID сессии для сохранения в историю |

**Успешный ответ (200):** аналогичен ответу `/api/llm/text`.

---

### History 🔒

#### Получение истории сессии

```
GET /api/history/:id
Authorization: Bearer <token>
```

Где `:id` — UUID сессии.

**Успешный ответ (200):**

```json
{
  "corrections": [
    {
      "id": "uuid",
      "sessionId": "uuid",
      "userId": "uuid",
      "inputContent": "Текст ТЗ",
      "responseData": {
        "data": {
          "rules_checker": { ... },
          "standard_checker": { ... }
        },
        "success": true,
        "code": 200
      },
      "createdAt": "2026-04-27T12:00:00Z"
    }
  ]
}
```

**Ошибка — невалидный id (400):**

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
│   └── main.go               # Точка входа
├── config/
│   └── config.yaml           # Конфигурация
├── docker/
│   └── Dockerfile            # Docker образ
├── tests/
│   ├── auth_test.go
│   └── history_test.go
├── internal/
│   ├── auth/
│   │   ├── models/           # DTO
│   │   ├── repository/       # Работа с БД
│   │   ├── route/            # HTTP хендлеры
│   │   └── service/          # Бизнес-логика
│   ├── session/
│   │   ├── models/           # DTO
│   │   ├── repository/       # Работа с БД
│   │   ├── route/            # HTTP хендлеры
│   │   └── service/          # Бизнес-логика
│   ├── history/
│   │   ├── models/           # DTO
│   │   ├── repository/       # Работа с БД
│   │   ├── route/            # HTTP хендлеры
│   │   └── service/          # Бизнес-логика
│   ├── llm/
│   │   ├── models/           # DTO
│   │   ├── service/          # Прокси к llm_api
│   │   └── route/            # HTTP хендлеры
│   ├── config/               # Загрузка конфига
│   └── migrations/           # Миграции БД
├── pkg/
│   ├── logger/               # Логгер (zap)
│   ├── postgres/             # Подключение к PostgreSQL (pgxpool)
│   └── registr/              # JWT токены
├── docker-compose.yml        # Docker Compose
├── go.mod
└── README.md
```