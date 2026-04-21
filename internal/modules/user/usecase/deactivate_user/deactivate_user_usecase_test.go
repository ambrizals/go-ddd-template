package deactivate_user

import (
	"context"
	"errors"
	"testing"

	"github.com/ambrizals/go-ddd-template/internal/modules/user/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.User), args.Error(1)
}

func TestDeactivateUserUseCase_Execute(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewDeactivateUserUseCase(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		input := DeactivateUserInput{
			ID: 1,
		}

		existingUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Test User",
		}

		mockRepo.On("GetByID", ctx, input.ID).Return(existingUser, nil)

		result, err := uc.Execute(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.ID, result.ID)
		assert.True(t, result.Deactivated)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		input := DeactivateUserInput{
			ID: 999,
		}

		mockRepo.On("GetByID", ctx, input.ID).Return(nil, errors.New("not found"))

		result, err := uc.Execute(ctx, input)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, entity.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found - nil user", func(t *testing.T) {
		input := DeactivateUserInput{
			ID: 999,
		}

		mockRepo.On("GetByID", ctx, input.ID).Return(nil, nil)

		result, err := uc.Execute(ctx, input)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, entity.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNewDeactivateUserUseCase(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewDeactivateUserUseCase(mockRepo)

	assert.NotNil(t, uc)
	assert.Equal(t, mockRepo, uc.userRepo)
}