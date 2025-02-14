# Используем официальный образ Go как базовый
FROM golang:1.24 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код проекта
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/service

# Используем минимальный образ для запуска
FROM alpine:latest

# Устанавливаем необходимые зависимости
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Устанавливаем переменную окружения для пути к конфигу
ENV CONFIG_PATH=.env

# Открываем порт, который использует приложение
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]