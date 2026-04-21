# Entity Layer

This layer contains the core domain models and interfaces. It is the center of the DDD architecture and should have no dependencies on other layers.

## Responsibilities
- Define Domain Models (Structs).
- Define Repository Interfaces.
- Define Use Case (Service) Interfaces.

## Example: Adding a new Entity
1. Create a new file `internal/entity/product.go`.
2. Define the struct and interfaces:

```go
package entity

type Product struct {
    ID    uint    `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

type ProductRepository interface {
    Create(ctx context.Context, p *Product) error
    // ...
}

type ProductUseCase interface {
    CreateProduct(ctx context.Context, p *Product) error
    // ...
}
```
