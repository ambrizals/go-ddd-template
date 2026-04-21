# Handler Layer (Delivery)

This layer implements the HTTP handlers using the Fiber framework. It is responsible for parsing requests, validating inputs (via DTOs), calling Use Cases, and returning HTTP responses.

## Responsibilities
- Parse HTTP requests (Body, Query, Params).
- Validate inputs using `go-playground/validator` and DTOs.
- Call Use Case methods.
- Handle HTTP status codes and response formatting.
- Include Swaggo annotations for OpenAPI documentation.

## Example: Adding a new Handler
1. Create a new subdirectory `internal/handler/product`.
2. Create a new file `internal/handler/product/product_handler.go`.
2. Implement the handler:

```go
package product

import (
    "github.com/gofiber/fiber/v2"
    "github.com/user/go-ddd-template/internal/dto"
    "github.com/user/go-ddd-template/internal/entity"
)

type ProductHandler struct {
    productUseCase entity.ProductUseCase
}

func NewProductHandler(useCase entity.ProductUseCase) *ProductHandler {
    return &ProductHandler{productUseCase: useCase}
}

// CreateProduct godoc
// @Summary Create a product
// @Tags products
// @Param request body dto.CreateProductRequest true "Product info"
// @Success 201 {object} dto.ProductResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
    // ... logic
}
```
