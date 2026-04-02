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
docker-compose up --build
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
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "role": "user"
}
```

Ответ:
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "role": "user"
  }
}
```

### Вход

```bash
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Ответ:
```json
{
  "token": "jwt-token"
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
├── docker-compose.yml       # Docker Compose
├── go.mod
└── README.md
```
