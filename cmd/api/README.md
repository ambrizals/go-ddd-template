# API Entry Point

This directory contains the main application entry point and server setup.

## Responsibilities
- Load configuration.
- Initialize infrastructure components (Logger, DB, Redis, OTEL).
- Wire up Dependency Injection.
- Register routes and middlewares.
- Start the Fiber server.

## How to Run
```bash
go run cmd/api/main.go
```
or
```bash
make run
```
