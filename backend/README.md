# TZshka

Backend проект с авторизацией и регистрацией.

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
| Email | string | Email пользователя |
| Role | string | Роль пользователя (admin/user) |
| CreatedAt | timestamp | Дата создания |

## API

### Регистрация

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

**Ошибка - невалидный role (400):**
```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "invalid request"
  }
}
```

### Вход

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
│   │   ├── domain/          # Доменные модели
│   │   ├── models/          # DTO
│   │   ├── repository/      # Работа с БД
│   │   ├── route/           # HTTP хендлеры
│   │   └── service/         # Бизнес-логика
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
