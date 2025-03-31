# 1. Промежуточный образ для сборки
FROM golang:1.23 as builder

RUN apt-get update && apt-get install -y libssl-dev pkg-config build-essential

# Рабочая директория сборки
WORKDIR /app

# 2. Скопируем только go.mod/go.sum из bot/
COPY bot/go.mod bot/go.sum ./

# 3. Установим зависимости
RUN go mod tidy && go mod download

# 4. Скопируем всё содержимое папки bot в /app
COPY bot /app

# В итоге в /app лежат:
#   go.mod, go.sum
#   cmd/main.go
#   internal/
#   ... и т.д.

# 5. Собираем бинарник
RUN go build -o /app/vote-bot ./cmd/main.go

# 6. Финальный контейнер
FROM alpine:latest
WORKDIR /app

# Скопируем только бинарник
COPY --from=builder /app/vote-bot .

CMD ["./vote-bot"]
