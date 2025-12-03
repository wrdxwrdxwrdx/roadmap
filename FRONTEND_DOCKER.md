# Frontend Docker Setup

Документация по запуску frontend через Docker.

## Быстрый старт

### Production режим

```bash
# Запустить все сервисы (postgres, api, frontend)
make up

# Frontend будет доступен на http://localhost:3000
```

### Development режим (с hot reload)

```bash
# Запустить все сервисы с frontend в dev режиме
make up-dev

# Frontend будет доступен на http://localhost:3000
# Изменения в коде будут автоматически применяться (hot reload)
```

## Доступные команды

### Основные команды

- `make up` - Запустить все сервисы в production режиме
- `make up-dev` - Запустить все сервисы с frontend в dev режиме
- `make down` - Остановить все сервисы
- `make build` - Собрать все Docker образы
- `make rebuild` - Пересобрать все Docker образы без кэша

### Frontend команды

- `make frontend-build-docker` - Собрать frontend Docker образ для production
- `make frontend-dev-docker` - Запустить frontend в dev режиме (отдельно)
- `make frontend-restart` - Перезапустить frontend сервис
- `make frontend-logs` - Просмотреть логи frontend
- `make frontend-shell` - Открыть shell в frontend контейнере
- `make frontend-clean-docker` - Очистить frontend Docker образы и контейнеры

### Просмотр логов

- `make logs` - Все логи
- `make logs-frontend` - Только логи frontend
- `make logs-api` - Только логи API
- `make logs-db` - Только логи PostgreSQL

## Структура Docker файлов

### Production

- `frontend/Dockerfile` - Multi-stage build:
  - Stage 1: Сборка приложения (Node.js)
  - Stage 2: Раздача статики (Nginx)

- `frontend/nginx.conf` - Конфигурация Nginx для production:
  - SPA fallback (все маршруты → index.html)
  - Прокси для API запросов (`/api` → `http://api:8080`)
  - Gzip compression
  - Кэширование статических файлов
  - Security headers

### Development

- `frontend/Dockerfile.dev` - Dev контейнер с Vite dev server
- `docker-compose.dev.yml` - Override конфигурация для dev режима:
  - Volume mounts для hot reload
  - Polling для файловой системы

## Переменные окружения

Можно настроить через `.env` файл или переменные окружения:

- `FRONTEND_PORT` - Порт для frontend (по умолчанию: 3000)
- `API_PORT` - Порт для API (по умолчанию: 8080)
- `FRONTEND_DOCKERFILE` - Какой Dockerfile использовать (по умолчанию: Dockerfile)

## Особенности

### Production

- Frontend собирается в статические файлы
- Nginx раздает статику и проксирует API запросы
- Оптимизирован для production (минификация, кэширование)

### Development

- Используется Vite dev server с hot module replacement (HMR)
- Исходный код монтируется как volume для мгновенных изменений
- Polling включен для работы в Docker

## Решение проблем

### Frontend не обновляется в dev режиме

1. Проверьте, что volumes правильно смонтированы:
   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.dev.yml ps
   ```

2. Проверьте логи:
   ```bash
   make frontend-logs
   ```

### API запросы не работают

1. Проверьте, что API контейнер запущен:
   ```bash
   make ps
   ```

2. В dev режиме proxy настроен на `http://api:8080`
3. В production nginx проксирует `/api` на `http://api:8080`

### Очистка

Если нужно полностью пересобрать frontend:

```bash
make frontend-clean-docker
make frontend-build-docker
```

