FROM golang:1.23 as builder

RUN apt-get update && apt-get install -y libssl-dev pkg-config build-essential

WORKDIR /app
COPY bot/go.mod bot/go.sum ./

RUN go mod tidy && go mod download

COPY bot /app
RUN go build -o /app/vote-bot ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/vote-bot .

CMD ["./vote-bot"]
