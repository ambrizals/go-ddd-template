# User Module

User management bounded context following DDD principles.

## Overview

This module handles all user-related functionality including registration, profile management, and account deactivation.

## Structure

```
user/
├── entity/           # Domain model and interfaces
│   └── user.go       # User entity, repository & use case interfaces
├── repository/       # Data access layer
│   └── user_repository.go
├── usecase/          # Business logic (workflow-specific)
│   ├── create_user/      # User registration workflow
│   ├── update_user/      # User update workflow
│   └── deactivate_user/  # User deactivation workflow
├── handler/          # HTTP layer
│   └── http/
│       ├── user_handler.go  # HTTP handlers
│       └── user_routes.go   # Route definitions
└── dto/              # Data transfer objects
    ├── user_dto.go           # Request DTOs
    └── user_response_dto.go  # Response DTOs
```

## Domain Model

### User Entity
- `ID` - Primary key
- `Email` - Unique user email
- `Password` - Hashed password
- `FullName` - User's full name
- `CreatedAt`, `UpdatedAt` - Timestamps

### Repository Interface
- `Create(ctx, user)` - Create new user
- `GetByID(ctx, id)` - Get user by ID
- `GetByEmail(ctx, email)` - Get user by email
- `List(ctx)` - List all users

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/users` | Register new user |
| GET | `/api/v1/users/:id` | Get user by ID |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Deactivate user |

## Use Cases

### CreateUserUseCase
Handles user registration with email uniqueness validation.

**Input:**
- Email (required, valid email format)
- Password (required, min 6 characters)
- FullName (required)

**Output:**
- User ID
- Email
- FullName

### UpdateUserUseCase
Updates user profile information.

**Input:**
- ID (user ID)
- Email (optional)
- FullName (optional)

**Output:**
- User ID
- Email
- FullName

### DeactivateUserUseCase
Deactivates a user account (soft delete).

**Input:**
- ID (user ID)

**Output:**
- ID
- Deactivated (boolean)

## Adding a New Use Case

1. Create directory `internal/modules/user/usecase/new_feature/`
2. Implement `NewNewFeatureUseCase` constructor
3. Add `Execute` method with input struct
4. Add DTOs in `internal/modules/user/dto/`
5. Register in handler

```go
// internal/modules/user/usecase/new_feature/new_feature_usecase.go
type NewFeatureInput struct {
    // fields
}

type NewFeatureOutput struct {
    // fields
}

type NewFeatureUseCase struct {
    userRepo entity.UserRepository
}

func NewNewFeatureUseCase(userRepo entity.UserRepository) *NewFeatureUseCase {
    return &NewFeatureUseCase{userRepo: userRepo}
}

func (uc *NewFeatureUseCase) Execute(ctx context.Context, input NewFeatureInput) (*NewFeatureOutput, error) {
    // implementation
}
```