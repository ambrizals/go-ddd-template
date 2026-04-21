package update_user

import (
	"context"

	"github.com/ambrizals/go-ddd-template/internal/modules/user/entity"
)

type UpdateUserUseCase struct {
	userRepo entity.UserRepository
}

// NewUpdateUserUseCase creates a new UpdateUserUseCase implementation.
func NewUpdateUserUseCase(repo entity.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{userRepo: repo}
}

// GetUser retrieves a user by ID
func (uc *UpdateUserUseCase) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

// UpdateUserInput represents the input for updating a user
type UpdateUserInput struct {
	ID       uint   `json:"id" validate:"required"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
	FullName string `json:"full_name" validate:"omitempty"`
}

// UpdateUserOutput represents the output for updating a user
type UpdateUserOutput struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

// Execute updates an existing user
func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	// Get existing user
	existingUser, err := uc.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, entity.ErrUserNotFound
	}
	if existingUser == nil {
		return nil, entity.ErrUserNotFound
	}

	// Check if email is being changed and if it's already taken
	if input.Email != "" && input.Email != existingUser.Email {
		emailUser, _ := uc.userRepo.GetByEmail(ctx, input.Email)
		if emailUser != nil {
			return nil, entity.ErrEmailAlreadyExists
		}
	}

	// Update fields if provided
	if input.Email != "" {
		existingUser.Email = input.Email
	}
	if input.Password != "" {
		existingUser.Password = input.Password // Should be hashed in real implementation
	}
	if input.FullName != "" {
		existingUser.FullName = input.FullName
	}

	if err := uc.userRepo.Create(ctx, existingUser); err != nil {
		return nil, err
	}

	return &UpdateUserOutput{
		ID:       existingUser.ID,
		Email:    existingUser.Email,
		FullName: existingUser.FullName,
	}, nil
}
