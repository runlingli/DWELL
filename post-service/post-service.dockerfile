# Build stage
FROM golang:1.21-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o postApp ./cmd/api

RUN chmod +x /app/postApp

# Production stage
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/postApp /app

CMD ["/app/postApp"]
