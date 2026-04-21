# DTO Layer (Data Transfer Objects)

This layer contains the input and output models for the API. It is decoupled from the internal domain models to allow for API versioning and to avoid exposing database-specific tags.

## Responsibilities
- Define Request/Response structures.
- Use `validate` tags for input validation (`go-playground/validator`).
- Use `example` tags for Swagger/OpenAPI documentation.

## Example: Adding a new DTO
1. Create a new file `internal/dto/product_dto.go`.
2. Define the Request and Response:

```go
package dto

type CreateProductRequest struct {
    Name  string  `json:"name" validate:"required" example:"Laptop"`
    Price float64 `json:"price" validate:"gt=0" example:"999.99"`
}

type ProductResponse struct {
    ID    uint    `json:"id" example:"1"`
    Name  string  `json:"name" example:"Laptop"`
    Price float64 `json:"price" example:"999.99"`
}
```
