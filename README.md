# TZshka

Полноценное веб-приложение с авторизацией.

## Структура

```
TZshka/
├── backend/              # Go backend
├── frontend/             # Frontend (скоро)
├── docker-compose.yml    # Запуск всех сервисов
├── .env.example          # Переменные окружения
└── README.md            # Этот файл
```

## Быстрый старт

### Docker

```bash
docker-compose up --build
```

### Локально (только backend)

```bash
cd backend
go build -o app.exe ./cmd/main.go
./app.exe
```

## API Endpoints

- `POST /register` - Регистрация
- `POST /login` - Вход

## Тесты

```bash
cd backend
go test -v ./tests/...
```
