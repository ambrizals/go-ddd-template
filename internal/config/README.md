# Config Layer

This layer handles application configuration using environment variables and `viper`.

## Responsibilities
- Define the configuration structure.
- Set default values.
- Load environment variables.

## Usage
Add new environment variables to the `Config` struct in `internal/config/config.go` and set defaults in `LoadConfig()`.

```go
type Config struct {
    // ...
    NewServiceKey string `mapstructure:"NEW_SERVICE_KEY"`
}
```
