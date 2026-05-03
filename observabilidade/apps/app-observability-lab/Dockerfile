FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY main.go ./
COPY controllers ./controllers

RUN go build -o app-observability-lab .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app-observability-lab .

EXPOSE 8080

CMD ["./app-observability-lab"]
