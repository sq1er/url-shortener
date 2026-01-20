# URL Shortener

Сервис для сокращения ссылок на **Go** с использованием GORM + PostgreSQL. Проект включает HTTP-хендлеры, службы и репозитории, JWT-аутентификацию, простой event-bus для асинхронной статистики и тесты (unit + интеграционные сценарии).

---

## Содержание <!-- TOC -->

- [Быстрый старт](#быстрый-старт)
- [Возможности](#возможности)
- [Стек технологий](#стек-технологий)
- [Структура проекта](#структура-проекта)
- [Требования](#требования)
- [Переменные окружения](#переменные-окружения)
- [Docker Compose (Postgres)](#docker-compose-postgres)
- [Запуск приложения](#запуск-приложения)
- [Тесты](#тесты)
- [API](#api)
- [БД и миграции](#бд-и-миграции)
- [Реализация и безопасность](#реализация-и-безопасность)
- [Дальнейшие планы](#дальнейшие-планы)
- [Вклад и автор](#вклад-и-автор)

---

Сервис для сокращения ссылок на **Go** с использованием GORM + PostgreSQL.  
Проект включает HTTP-хендлеры, сервисы и репозитории, JWT-аутентификацию, простой event-bus для асинхронной статистики и набор тестов (unit + интеграционные сценарии).

---

## 0. Установка и клонирование репозитория

Клонируйте репозиторий и перейдите в каталог проекта:

```bash
git clone https://github.com/sq1er/url-shortener
cd url-shortener
```

---

## 1. Возможности

- Регистрация и вход пользователей (JWT)
- Создание / обновление / удаление коротких ссылок
- Редирект по хэшу (`GET /{hash}`)
- Сбор и агрегация статистики (по дням / месяцам)
- Middleware: CORS, логирование, авторизация
- GORM модели + `AutoMigrate`
- Тесты: unit, sqlmock, e2e и интеграционные

---

## 2. Стек технологий

- **Язык:** Go 1.24+
- **БД:** PostgreSQL
- **ORM:** GORM
- **Аутентификация:** JWT
- **Асинхронность:** внутренняя шина событий (channels)
- **Инфраструктура:** Docker / Docker Compose

---

## 3. Структура проекта

```
cmd/               -> точка входа (HTTP-сервер)
configs/           -> загрузка конфигурации
internal/          -> доменные модули (auth, link, stat, user)
pkg/               -> общие и инфраструктурные пакеты (db, jwt, middleware, event)
migrations/        -> миграции / AutoMigrate
docker-compose.yml -> Postgres
```

`internal/*` организован по паттерну: **handler → service → repository**.

---

## 4. Требования

- Go 1.24+
- Docker (опционально, для запуска Postgres)
- Интернет-доступ для загрузки Go-модулей

---

## 5. Переменные окружения

Создайте файл `.env` в корне проекта:

```env
DSN="host=localhost user=postgres password=my_pass dbname=link port=5432 sslmode=disable"
SECRET="your_jwt_secret_here"
```

> Для тестов используется база `link_test`.

---

## 6. Docker Compose (Postgres)

```yaml
services:
  postgres:
	container_name: postgres_go
	image: postgres
	environment:
	  POSTGRES_USER: postgres
	  POSTGRES_PASSWORD: my_pass
	  PGDATA: /data/postgres
	volumes:
	  - ./postgres-data:/data/postgres
	ports:
	  - "5432:5432"
```

Запуск БД:

```bash
docker compose up -d
```

Создание базы данных:

```bash
docker exec -it postgres_go psql -U postgres -c "CREATE DATABASE link;"
```

---

## 7. Запуск приложения

Миграция схемы (AutoMigrate):

```bash
go run migrations/auto.go
```

Запуск сервера:

```bash
go run cmd/main.go
```

---

## 8. API

Все запросы и ответы в формате JSON.

### Auth

**POST /auth/register**

```json
{
  "email": "user@mail.ru",
  "password": "my_pass",
  "name": "Александр"
}
```

**POST /auth/login**

```json
{
  "name": "Александр",
  "password": "my_pass"
}
```

Ответ:

```json
{ "token": "jwt" }
```

### Link

- **POST /link** — создать ссылку (JWT)
- **PATCH /link/{id}** — обновить ссылку
- **DELETE /link/{id}** — удалить ссылку
- **GET /link?limit=5&offset=0** — список ссылок
- **GET /{hash}** — редирект на оригинальный URL

Header:

```
Authorization: Bearer <jwt>
```

### Stat

**GET /stat?from=YYYY-MM-DD&to=YYYY-MM-DD&by=day|month**

Ответ:

```json
[{ "period": "2026-01-01", "sum": 10 }]
```

---

## 9. Реализация

- Асинхронный сбор статистики через `event.EventBus`
- Пароли хранятся в виде bcrypt-хеша
- JWT для авторизации
- AutoMigrate для быстрой разработки

---

## 10. Автор

Учебный проект.  
Идеи и улучшения приветствуются — создавайте Issue или Pull Request.
