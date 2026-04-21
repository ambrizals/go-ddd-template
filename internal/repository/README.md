# Repository Layer (Infrastructure Implementation)

This layer implements the `entity.UserRepository` interface using GORM and PostgreSQL. It is responsible for all database interactions.

## Responsibilities
- Implement data access logic.
- Manage DB transactions if needed.
- Hide implementation details (GORM/SQL) from the domain.

## Example: Adding a new Repository
1. Create a new file `internal/repository/product_repository.go`.
2. Implement the interface:

```go
package repository

import (
    "context"
    "github.com/user/go-ddd-template/internal/entity"
    "gorm.io/gorm"
)

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) entity.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *entity.Product) error {
    return r.db.WithContext(ctx).Create(p).Error
}
```
