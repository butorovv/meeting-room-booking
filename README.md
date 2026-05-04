# Meeting Room Booking Service

Backend-сервис бронирования переговорных комнат.

## Стек

Go, PostgreSQL, Redis, Docker, Docker Compose, Swagger, JWT, gomock, testify.

## Архитектура

Проект построен по принципам Clean Architecture:

```text
delivery -> usecase -> repository -> domain
```

- `delivery` — HTTP-обработчики и middleware;
- `usecase` — бизнес-логика;
- `repository` — работа с PostgreSQL;
- `domain` — бизнес-сущности и ошибки;
- `pkg` — переиспользуемые компоненты.

## Что реализовано

- генерация слотов на лету из расписания комнаты;
- стабильный ID слота на основе `room_id + start_time + end_time`;
- защита от двойного бронирования через уникальный индекс PostgreSQL;
- Redis-кеширование списка доступных слотов;
- инвалидация кеша при создании и отмене брони;
- асинхронная отправка уведомлений через in-memory worker pool;
- graceful shutdown;
- Swagger-документация;
- unit- и e2e-тесты.

## Запуск

```bash
git clone https://github.com/butorovv/meeting-room-booking.git
cd meeting-room-booking
docker-compose up -d --build
make seed
```

Swagger:

```text
http://localhost:8080/swagger/index.html
```

## Тестирование

```bash
make test
make e2e
```

## Структура

```text
cmd/            - точка входа
internal/       - основной код
pkg/            - переиспользуемые компоненты
db/migrations/  - SQL-миграции
scripts/        - e2e и seed
docs/           - Swagger
```

## Контакты

GitHub: github.com/butorovv  
Telegram: @butorovvv  
Email: george.butorov@mail.ru