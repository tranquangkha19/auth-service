FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/server

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/auth-service .
COPY configs/config.yaml .

EXPOSE 8080

CMD ["./auth-service"]
