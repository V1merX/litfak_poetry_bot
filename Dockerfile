# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/bot ./cmd/bot

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/bot /app/bot
COPY --from=builder /app/configs /app/configs

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Run as non-root user
RUN adduser -D -g '' appuser
USER appuser

CMD ["/app/bot"]