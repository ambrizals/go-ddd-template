# Shared Module

Cross-cutting concerns shared across all bounded contexts.

## Overview

This module provides utilities and infrastructure code that are used by multiple modules throughout the application.

## Structure

```
shared/
├── config/       # Configuration management
│   └── config.go
├── logger/       # Structured logging
│   └── logger.go
└── response/     # API response formatting
    ├── response.go
    └── response_dto.go
```

## Components

### Config (`config/`)
Environment-based configuration using Viper.

**Features:**
- Load from `.env` file
- Environment variable overrides
- Type-safe access

**Usage:**
```go
cfg := config.LoadConfig()
dbHost := cfg.GetString("DB_HOST")
dbPort := cfg.GetInt("DB_PORT")
```

### Logger (`logger/`)
Structured logging using Go's slog package.

**Features:**
- JSON output format
- Log levels (debug, info, warn, error)
- Context-aware logging

**Usage:**
```go
logger := logger.GetLogger()
logger.Info("message", "key", "value")
logger.Error("error occurred", "err", err)
```

### Response (`response/`)
Standardized API response formatting.

**Response Structure:**
```json
{
  "success": true,
  "message": "Success",
  "data": { ... }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error title",
  "error": "Error description"
}
```

**Usage:**
```go
// Success response
return c.JSON(response.NewAPIResponse(payload))

// Error response
return c.Status(400).JSON(response.NewAPIErrorResponse[UserResponse]("Bad Request", "Invalid input"))
```

## Adding New Shared Utilities

1. Create directory under `internal/shared/`
2. Implement the utility
3. Export public functions/types
4. Add tests in the same directory

Example:
```go
// internal/shared/newutility/newutility.go
package newutility

func DoSomething() error {
    // implementation
}
```