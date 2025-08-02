FROM golang:1.24.2 AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go build -o prometheus-consumer ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/prometheus-consumer .
COPY ./internal/config ./config

EXPOSE 8090

RUN apk add libc6-compat

CMD ["./prometheus-consumer"]