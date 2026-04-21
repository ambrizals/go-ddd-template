# Use Case Layer (Business Logic)

This layer implements the `entity.UserUseCase` interface. It orchestrates the domain models and coordinates with the repository to execute business rules.

## Responsibilities
- Implement Business Rules.
- Orchestrate one or more Repositories.
- Map between Entities and DTOs if necessary (or leave it to Handlers).
- Handle Domain Errors.

## Example: Adding a new Use Case
1. Create a new file `internal/usecase/product_usecase.go`.
2. Implement the interface:

```go
package usecase

import (
    "context"
    "github.com/user/go-ddd-template/internal/entity"
)

type productUseCase struct {
    productRepo entity.ProductRepository
}

func NewProductUseCase(repo entity.ProductRepository) entity.ProductUseCase {
    return &productUseCase{productRepo: repo}
}

func (u *productUseCase) CreateProduct(ctx context.Context, p *entity.Product) error {
    // Add business validation
    if p.Price <= 0 {
        return errors.New("invalid price")
    }
    return u.productRepo.Create(ctx, p)
}
```
