# Go DDD Backend Template

This is a production-ready Go backend template following Domain-Driven Design (DDD) principles.

## Tech Stack
- **Framework**: [Fiber](https://gofiber.io/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **OpenAPI/Swagger**: [swaggo/swag](https://github.com/swaggo/swag)
- **Caching**: Redis
- **Observability**: OpenTelemetry
- **Config**: Viper

## Project Structure
```text
.
├── cmd/api             # Entry point
├── docs                # Swagger docs (auto-generated)
├── internal/
│   ├── entity          # Domain models & Interfaces
│   ├── usecase         # Business Logic
│   ├── repository      # Data Source Implementation
│   ├── handler         # HTTP Handlers (Delivery)
│   ├── dto             # Data Transfer Objects & Validation
│   ├── config          # Configuration
│   ├── infrastructure  # DB, Redis, OTEL, Logger
│   └── middleware      # Fiber Middlewares
├── pkg                 # Shared utilities
└── test                # Integration tests
```

## Getting Started

### Prerequisites
- Go 1.22+
- Docker & Docker Compose

### Setup
1. Clone the repository.
2. Run the environment:
   ```bash
   docker-compose up -d
   ```
3. Run the application:
   ```bash
   make run
   ```

### Swagger Documentation
The API documentation is available at `http://localhost:8080/swagger`.

To regenerate swagger docs:
```bash
make swag-gen
```

### Database Migrations
Migrations are automatically applied on startup. To add a new migration:
1. Create `.up.sql` and `.down.sql` files in `internal/infrastructure/database/migrations`.

### Running Tests
```bash
make test
```
