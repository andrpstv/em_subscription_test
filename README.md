# Subscription Service

## Архитектура

Проект реализован с использованием многослойной архитектуры для обеспечения масштабируемости:

- **Repository слой**: Работа с базой данных через sqlx
- **Service слой**: Бизнес-логика и валидация
- **Handler слой**: HTTP обработчики
- **Config**: Управление конфигурацией
- **Migrations**: Миграции базы данных через Goose

## Технологии

- **Go 1.23**
- **Gin** - HTTP фреймворк
- **sqlx** - Расширение database/sql для удобной работы с PostgreSQL
- **Goose** - Миграции базы данных
- **Logrus** - Логирование
- **Swagger** - Документация API
- **Docker & Docker Compose** - Контейнеризация

## Запуск

### Запуск через Docker Compose
```bash
docker-compose up --build
```

## API Документация

Swagger документация доступна по адресу: `http://localhost:8080/swagger/index.html`

### Основные эндпоинты

#### Подписки
- `POST /api/v1/subscriptions` - Создание подписки
- `GET /api/v1/subscriptions` - Список подписок (с фильтрами)
- `GET /api/v1/subscriptions/{id}` - Получение подписки по ID
- `PUT /api/v1/subscriptions/{id}` - Обновление подписки
- `DELETE /api/v1/subscriptions/{id}` - Удаление подписки

#### Расчет стоимости
- `POST /api/v1/subscriptions/total-cost` - Расчет суммарной стоимости за период

### Пример запроса на создание подписки
```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
```

## Конфигурация

Конфигурационные переменные задаются через переменные окружения:

- `DB_HOST` - Хост PostgreSQL
- `DB_PORT` - Порт PostgreSQL
- `DB_USER` - Пользователь БД
- `DB_PASSWORD` - Пароль БД
- `DB_NAME` - Имя базы данных
- `DB_SSLMODE` - Режим SSL
- `SERVER_PORT` - Порт сервера
- `LOG_LEVEL` - Уровень логирования
