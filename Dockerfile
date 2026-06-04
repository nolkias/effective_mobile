FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o subscription-service .

# Финальный этап
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник из этапа сборки
COPY --from=builder /app/subscription-service .

# Копируем миграции и .env
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env .env

EXPOSE 8080

CMD ["./subscription-service"]