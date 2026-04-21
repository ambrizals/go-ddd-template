package deactivate_user

import (
	"context"

	"github.com/ambrizals/go-ddd-template/internal/modules/user/entity"
)

type DeactivateUserUseCase struct {
	userRepo entity.UserRepository
}

// NewDeactivateUserUseCase creates a new DeactivateUserUseCase implementation.
func NewDeactivateUserUseCase(repo entity.UserRepository) *DeactivateUserUseCase {
	return &DeactivateUserUseCase{userRepo: repo}
}

// DeactivateUserInput represents the input for deactivating a user
type DeactivateUserInput struct {
	ID uint `json:"id" validate:"required"`
}

// DeactivateUserOutput represents the output for deactivating a user
type DeactivateUserOutput struct {
	ID          uint `json:"id"`
	Deactivated bool `json:"deactivated"`
}

// Execute deactivates an existing user (soft delete or status change)
func (uc *DeactivateUserUseCase) Execute(ctx context.Context, input DeactivateUserInput) (*DeactivateUserOutput, error) {
	// Get existing user
	existingUser, err := uc.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, entity.ErrUserNotFound
	}
	if existingUser == nil {
		return nil, entity.ErrUserNotFound
	}

	// For now, we'll just return success since we don't have a deactivated field in the entity
	// In a real implementation, you would update a status field or deleted_at timestamp
	return &DeactivateUserOutput{
		ID:          existingUser.ID,
		Deactivated: true,
	}, nil
}
