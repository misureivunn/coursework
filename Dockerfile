# Multi-stage build для Go приложения с Fyne UI

# Stage 1: Builder
FROM golang:1.22-alpine AS builder

# Установка зависимостей для компиляции Fyne
RUN apk add --no-cache \ 
    gcc \ 
    musl-dev \ 
    mesa-dev \ 
    libx11-dev \ 
    libxcursor-dev \ 
    libxrandr-dev \ 
    libxinerama-dev \ 
    libxi-dev \ 
    libgl1-mesa-dev

WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Сборка приложения с поддержкой CGO (требуется для Fyne)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o szi-registry .

# Stage 2: Runtime (минимальный образ)
FROM alpine:latest

# Установка runtime зависимостей для GUI
RUN apk add --no-cache \ 
    ca-certificates \ 
    libx11 \ 
    libxcursor \ 
    libxrandr \ 
    libxinerama \ 
    libxi \ 
    mesa-gl \ 
    font-noto

WORKDIR /app

# Копируем скомпилированный бинарник
COPY --from=builder /app/szi-registry .

# Создаем непривилегированного пользователя
RUN addgroup -g 1000 szi && \
    adduser -D -u 1000 -G szi szi && \
    chown -R szi:szi /app

USER szi

# Переменные окружения для подключения к БД
ENV DB_HOST=postgres-szi \ 
    DB_PORT=5432 \ 
    DB_USER=szi_user \ 
    DB_PASSWORD=MyV3ryS3cur3P@ss2026! \ 
    DB_NAME=szi_registry

# Для GUI-приложений нужен DISPLAY (при запуске с X11 forwarding)
ENV DISPLAY=:0

EXPOSE 8080

CMD ["./szi-registry"]