# Multi-stage build для frontend

# Stage 1: Build
FROM node:20-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
# Копируем оба файла - package-lock.json должен существовать в build context
COPY package.json package-lock.json ./

# Устанавливаем зависимости
# Используем npm ci для точного воспроизведения зависимостей из lockfile
# Если npm ci не работает, fallback на npm install
RUN npm ci || npm install

# Копируем исходный код
COPY . .

# Собираем приложение для production
RUN npm run build

# Stage 2: Production - используем nginx для раздачи статических файлов
FROM nginx:alpine

# Копируем собранные файлы из builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Копируем конфигурацию nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

