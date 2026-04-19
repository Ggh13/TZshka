# TZshka

Backend проект с авторизацией, регистрацией и управлением сессиями.

## Технологии

- Go 1.25
- PostgreSQL
- Gin (HTTP фреймворк)
- JWT (аутентификация)
- pgx (работа с PostgreSQL)

## Запуск

### Локально

```bash
cd backend
go build -o app.exe ./cmd/main.go
./app.exe
```

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

## API

### Регистрация

```bash
POST /api/register
Content-Type: application/json

{
  "login": "testuser",
  "email": "user@example.com",
  "password": "password123"
}
```

**Успешный ответ (201):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "login": "testuser",
    "email": "user@example.com"
  }
}
```

### Вход

```bash
POST /api/login
Content-Type: application/json

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

### Создание сессии

```bash
POST /api/sessions
Content-Type: application/json
Authorization: Bearer <token>

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

### Получение списка сессий

```bash
GET /api/sessions
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

### Вход в сессию

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

### Удаление сессии

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
{}
```

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
│   ├── auth_test.go         # Тесты auth
│   └── session_test.go      # Тесты session
├── internal/
│   ├── auth/
│   │   ├── domain/          # Доменные модели
│   │   ├── models/          # DTO
│   │   ├── repository/      # Работа с БД
│   │   ├── route/           # HTTP хендлеры
│   │   └── service/         # Бизнес-логика
│   ├── session/
│   │   ├── models/          # DTO
│   │   ├── repository/     # Работа с БД
│   │   ├── route/           # HTTP хендлеры
│   │   └── service/        # Бизнес-логика
│   ├── config/             # Загрузка конфига
│   └── migrations/         # Миграции БД
├── pkg/
│   ├── logger/            # Логгер
│   ├── postgres/          # Подключение к PostgreSQL
│   └── registr/           # JWT токены
├── go.mod
└── README.md
```