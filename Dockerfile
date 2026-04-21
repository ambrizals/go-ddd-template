# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum (if it exists)
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /api cmd/api/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /api .
COPY --from=builder /app/internal/infrastructure/database/migrations ./internal/infrastructure/database/migrations

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api"]
