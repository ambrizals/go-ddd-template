package create_user

import (
	"context"

	"github.com/ambrizals/go-ddd-template/internal/modules/user/entity"
)

type CreateUserUseCase struct {
	userRepo entity.UserRepository
}

// NewCreateUserUseCase creates a new CreateUserUseCase implementation.
func NewCreateUserUseCase(repo entity.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: repo}
}

// CreateUserInput represents the input for creating a user
type CreateUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

// CreateUserOutput represents the output for creating a user
type CreateUserOutput struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

// Execute creates a new user
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	// Business logic: check if email is already taken.
	existingUser, _ := uc.userRepo.GetByEmail(ctx, input.Email)
	if existingUser != nil {
		return nil, entity.ErrEmailAlreadyExists
	}

	user := &entity.User{
		Email:    input.Email,
		Password: input.Password, // Should be hashed in real implementation
		FullName: input.FullName,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
	}, nil
}
