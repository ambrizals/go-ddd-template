# Testing Layer

This directory contains integration and End-to-End (E2E) tests. Unit tests for each layer are located within the respective package directories (e.g., `internal/usecase/user_usecase_test.go`).

## Responsibilities
- Integration Tests (testing interactions between multiple layers).
- E2E Tests (testing the full HTTP flow with a real database).

## Example: Adding an Integration Test
1. Create a new file `test/user_integration_test.go`.
2. Use Fiber's Test method:

```go
func TestUserIntegration(t *testing.T) {
    app := fiber.New()
    // ... setup and register handlers ...
    
    resp, err := app.Test(httptest.NewRequest("GET", "/api/v1/users/1", nil))
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```
