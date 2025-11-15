# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/supa-manager main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates words wget && \
    ln -s /usr/share/dict/american-english /usr/share/dict/words

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/supa-manager .
COPY --from=builder /app/.env.example .
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Run the application
CMD ["./supa-manager"]
