# Go DDD Backend Template

A production-ready Go backend template following Domain-Driven Design (DDD) principles with Bounded Context organization.

## Tech Stack

| Category | Technology |
|----------|------------|
| Framework | [Fiber](https://gofiber.io/) |
| ORM | [GORM](https://gorm.io/) |
| Database | PostgreSQL |
| Caching | Redis |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| OpenAPI/Swagger | [swaggo/swag](https://github.com/swaggo/swag) |
| Observability | OpenTelemetry |
| Config | Viper |

## Project Structure

```
.
├── cmd/api/                  # Application entry point
├── internal/
│   ├── modules/              # Bounded contexts (DDD)
│   │   └── user/             # User module
│   │       ├── entity/       # Domain models & interfaces
│   │       ├── repository/   # Data access layer
│   │       ├── usecase/      # Business logic (workflow-specific)
│   │       ├── handler/      # HTTP handlers
│   │       └── dto/          # Data transfer objects
│   ├── shared/               # Shared utilities
│   │   ├── config/           # Configuration management
│   │   ├── logger/           # Structured logging
│   │   └── response/         # API response formatting
│   └── infrastructure/       # External services
│       ├── database/         # PostgreSQL & migrations
│       ├── redis/            # Redis client
│       └── otel/             # OpenTelemetry tracing
├── migrations/               # SQL migrations
├── test/                     # Integration tests
├── pkg/                      # Public utilities
├── docs/                     # Swagger documentation
└── docker-compose.yaml       # Local development environment
```

## Modules

### User Module (`internal/modules/user/`)
User management bounded context with:
- **Entity**: User domain model with repository interface
- **Use Cases**: Create, Update, Deactivate user workflows
- **Handler**: HTTP REST endpoints
- **DTOs**: Request/Response data transfer objects

### Shared (`internal/shared/`)
Cross-cutting concerns used across all modules:
- **config**: Environment-based configuration via Viper
- **logger**: Structured logging with slog
- **response**: Standardized API response formatting

### Infrastructure (`internal/infrastructure/`)
External service integrations:
- **database**: PostgreSQL connection, GORM, migrations
- **redis**: Caching layer
- **otel**: Distributed tracing

## Getting Started

### Prerequisites
- Go 1.23+
- Docker & Docker Compose

### Quick Start

```bash
# Start dependencies (PostgreSQL, Redis)
docker-compose up -d

# Run the application
make run
# or: go run cmd/api/main.go
```

### API Documentation
- Swagger UI: http://localhost:8080/swagger
- Scalar API: http://localhost:8080/scalar

### Available Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the application |
| `make test` | Run tests |
| `make swag-gen` | Generate Swagger docs |
| `make build` | Build binary |

### Environment Variables
Copy `.env.example` to `.env` and configure:
- `APP_PORT` - Server port (default: 8080)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL
- `REDIS_HOST`, `REDIS_PORT` - Redis