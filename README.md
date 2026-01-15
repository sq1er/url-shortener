# URL Shortener

Сервис для сокращения ссылок на Go используюзий (GORM + Postgres). Проект содержит HTTP-хендлеры, сервисы и репозитории, JWT-аутентификацию, простой eventbus для асинхронной статистики и набор тестов (unit + интеграционные сценарии).

---

## Основные возможности

- Регистрация / Вход (JWT)
- Создание / Обновление / Удаление коротких ссылок
- Редирект по хэшу (`GET /{hash}`)
- Базовая агрегация статистики (по дням / месяцам)
- Middleware: CORS, логирование, авторизация
- Модели GORM + `AutoMigrate` для разработки
- Тесты: e2e, `sqlmock`, интеграционные и unit тесты

---

## Структура репозитория

```
cmd/               -> точка входа (сервер)
configs/           -> загрузка конфигурации
internal/          -> доменные модули (auth, link, stat, user)
pkg/               -> инфраструктурные/общие пакеты (db, jwt, middleware, req, res, event)
migrations/        -> миграции
Docker-compose.yml -> Postgres
```

`internal/*` разделён по модулям и придерживается паттерна handler -> service -> repository.

---

## Требования

- Go 1.24+
- Docker (опционально, для Postgres)
- Доступ в интернет для скачивания модулей (GOPROXY)

---

## Переменные окружения (пример .env)

Создайте файл `.env` в корне проекта (или в `cmd/` для тестов) с содержимым:

```env
# Пример — не храните реальные секреты в публичных репозиториях
DSN="host=localhost user=postgres password=my_pass dbname=link port=5432 sslmode=disable"
SECRET="your_jwt_secret_here"
```

> Для тестов используется база `link_test` — укажите соответствующий `DSN` при запуске тестов.

---

## Быстрый старт (с Docker Postgres)

1. Запустить Postgres через Docker Compose (в проекте есть `docker-compose.yml`):

```bash
docker exec -it postgres_go psql -U postgres -c "CREATE DATABASE link;"
```

2. Мигрировать схему (AutoMigrate):

```bash
# выполнит миграцию моделей (links, users, stats)
go run migrations/auto.go
```

3. Запустить сервер:

```bash
go run cmd/main.go
# Server is listening on port 8081
```

---

## Тесты

Некоторые тесты используют `sqlmock` и не требуют реальной БД; интеграционные тесты требуют Postgres и корректного `.env`.

---

## API

Все запросы и ответы в формате JSON.

### Auth

- `POST /auth/register` — регистрация

```
- Request:
{
  "email": "user@mail.ru",
  "password":"my_pass",
  "name":"Александр"
}
- Response: 201
{
  "token": "jwt"
}
```

- `POST /auth/login` — вход

```
- Request:
{
  "name":"Александр",
  "password":"my_pass"
}
- Response: 200
{
  "token": "jwt"
}
```

### Link

Защищённые эндпоинты (JWT) — кроме редиректа:

- `POST /link` — создать ссылку (auth)

```
- Request:
{
	"url":"https://ya.ru/"
}
Authorization Bearer jwt

- Response: 201
{
	"ID": 1,
	"CreatedAt": "2026-01-15T19:19:05.5090954+03:00",
	"UpdatedAt": "2026-01-15T19:19:05.5090954+03:00",
	"DeletedAt": null,
	"url": "https://ya.ru/",
	"hash": "wyYvjO",
	"Stats": null
}
```

- `PATCH /link/{id}` — обновить ссылку (auth)

```
- Request:
{
  "ID": 1,
	"url":"https://yandex.ru"
}
Authorization Bearer jwt

- Response: 201
{
	"ID": 1,
	"CreatedAt": "2026-01-15T19:53:33.204199+03:00",
	"UpdatedAt": "2026-01-15T19:56:44.21956+03:00",
	"DeletedAt": null,
	"url": "http://yandex.ru",
	"hash": "cEfbIl",
	"Stats": null
}
```

- `DELETE /link/{id}` — удалить ссылку (auth)
- `GET link?limit=5&offset=0` — список ссылок (auth)
- `GET /{hash}` — переход по короткой ссылке. При переходе редиректится на длинную ссылку, ведется статистика переходов по каждой ссылке, когда и сколько раз были переходы

### Stat

- `GET /stat?from=YYYY-MM-DD&to=YYYY-MM-DD&by=day|month` (auth)
  - Response: `[{ "period":"2026-01-01", "sum": 10 }, ...]`

### Авторизация

Во всех защищённых запросах ожидается header:

```
Authorization: Bearer jwt
```

Токен — JWT, подписанный секретом из `SECRET`.

---

## Миграции и БД

Проект использует `AutoMigrate` GORM для удобства разработки.

Запуск автосинхронизации схемы:

```bash
go run migrations/auto.go
```

---

## Дизайн и реализация

- `internal/*` — доменные модули, каждый содержит handler/service/repository/model/payload
- `pkg/*` — утилитарные и инфраструктурные компоненты (db, jwt, middleware, req, res, event)
- `event.EventBus` — простая внутренняя шина на каналах для асинхронной обработки посещений ссылок
- Пароли хранятся в виде bcrypt-хеша
- JWT для аутентификации

---

## Автор

Если у вас есть идеи или предложения по улучшению — создайте Issue или Pull Request на GitHub.
