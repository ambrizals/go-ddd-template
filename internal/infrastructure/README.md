# Infrastructure Layer

This layer handles external concerns like database connections, caching, observability, and logging.

## Responsibilities
- Initialize database connections (GORM, Redis).
- Run database migrations (`golang-migrate`).
- Set up logging (`slog`).
- Set up observability (OpenTelemetry).

## Components
- **database**: GORM & Migration logic.
- **redis**: Redis client.
- **logger**: Structured logging.
- **otel**: OpenTelemetry tracing.

## Example: Adding a new Infrastructure Client
1. Create a new directory `internal/infrastructure/s3`.
2. Implement the client initialization.
3. Call it in `cmd/api/main.go`.

```go
package s3

func InitS3() {
    // ...
}
```
