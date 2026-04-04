# TZshka

## Быстрый старт

### Docker

```bash
# Первый запуск или сброс данных
docker-compose down -v
docker-compose up --build

# Последующие запуски
docker-compose up
```

## Структура

- `backend/` - Go backend
- `docker-compose.yml` - Запуск всех сервисов
