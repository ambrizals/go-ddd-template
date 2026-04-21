package usecase

import (
	"context"
	"github.com/user/go-ddd-template/internal/entity"
)

type userUseCase struct {
	userRepo entity.UserRepository
}

// NewUserUseCase creates a new UserUseCase implementation.
func NewUserUseCase(repo entity.UserRepository) entity.UserUseCase {
	return &userUseCase{userRepo: repo}
}

func (u *userUseCase) Register(ctx context.Context, user *entity.User) error {
	// Business logic: check if email is already taken.
	existingUser, _ := u.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return entity.ErrEmailAlreadyExists
	}

	return u.userRepo.Create(ctx, user)
}

func (u *userUseCase) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *userUseCase) ListUsers(ctx context.Context) ([]entity.User, error) {
	return u.userRepo.List(ctx)
}
