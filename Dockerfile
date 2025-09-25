# Начальная сборка с именеи builder, использует официальный образ Go на базе alpine linux
FROM golang:1.25.0-alpine AS builder

# Устанавливает рабочую директорию для всего приложения, будет выполнятся относительной этой директории (COPY, RUN, CMD)
WORKDIR /app

# Копирует файлы в рабочую директорию
COPY go.mod go.sum .env ./

# Выполняем go mod tidy
RUN go mod tidy

# Подгружает зависимости модуля Go, флаг "-x" - показывает подробную информацию о процессе загрузки, что помогает в отладке
RUN go mod download -x

# Копируем весь остальной исходный код
COPY . .

# Компилируем приложение в статический бинарник.
# Параметры: CGO_ENABLED=0      - отключает CGO, чтобы быть не зависимым от системных С-библиотек;
#            GOOS=linux         - указывает целевую ОС, для компиляции;
#            -ldflags="-s -w"   - флаги для уменьшения размера бинарника;
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Финальная сборка (чистая стадия с минимальным образом)
FROM alpine:3.22.1

# Создаёт непривилегированного пользователя и группу, чтобы запускать не от root, уменьшив ущерб при компрометации (от злоумышленников)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Устанавливает рабочую директорию в финальном образе
WORKDIR /app

# Копирует скомпилированный бинарник из стадии builder
COPY --from=builder --chown=appuser:appgroup /app/internal/config ./internal/config
COPY --from=builder --chown=appuser:appgroup /app/main .
COPY --from=builder --chown=appuser:appgroup /app/web ./web


# Копирует статические файлы и шаблоны
#COPY --from=builder /app/tempaltes ./tempaltes
#COPY --from=builder /app/static ./static

# Переключает на непривилегированного пользователя
USER appuser

# Документирует какой порт использет приложение
EXPOSE 8080

# Команда по умолчанию для запуска контейнера
CMD ["./main"]