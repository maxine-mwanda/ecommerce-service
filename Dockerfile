# Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN apk add --no-cache git
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ecommerce-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/ecommerce-service .

EXPOSE 8080
CMD ["./ecommerce-service"]

