# Build stage
FROM golang:1.21-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o favouriteApp ./cmd/api

RUN chmod +x /app/favouriteApp

# Production stage
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/favouriteApp /app

CMD ["/app/favouriteApp"]
