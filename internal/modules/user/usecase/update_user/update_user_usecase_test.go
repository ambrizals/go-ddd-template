package update_user

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

func TestUpdateUserUseCase_Execute(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUpdateUserUseCase(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		input := UpdateUserInput{
			ID:       1,
			FullName: "Updated Name",
		}

		existingUser := &entity.User{
			ID:       1,
			Email:    "test@example.com",
			FullName: "Old Name",
		}

		mockRepo.On("GetByID", ctx, input.ID).Return(existingUser, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		result, err := uc.Execute(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Name", result.FullName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		input := UpdateUserInput{
			ID:       999,
			FullName: "Updated Name",
		}

		mockRepo.On("GetByID", ctx, input.ID).Return(nil, errors.New("not found"))

		result, err := uc.Execute(ctx, input)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, entity.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found - nil user", func(t *testing.T) {
		input := UpdateUserInput{
			ID:       999,
			FullName: "Updated Name",
		}

		mockRepo.On("GetByID", ctx, input.ID).Return(nil, nil)

		result, err := uc.Execute(ctx, input)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, entity.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already exists when changing", func(t *testing.T) {
		input := UpdateUserInput{
			ID:    1,
			Email: "taken@example.com",
		}

		existingUser := &entity.User{
			ID:    1,
			Email: "test@example.com",
		}

		takenUser := &entity.User{ID: 2, Email: "taken@example.com"}

		mockRepo.On("GetByID", ctx, input.ID).Return(existingUser, nil)
		mockRepo.On("GetByEmail", ctx, input.Email).Return(takenUser, nil)

		result, err := uc.Execute(ctx, input)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmailAlreadyExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update email successfully", func(t *testing.T) {
		input := UpdateUserInput{
			ID:    1,
			Email: "newemail@example.com",
		}

		existingUser := &entity.User{
			ID:    1,
			Email: "old@example.com",
		}

		mockRepo.On("GetByID", ctx, input.ID).Return(existingUser, nil)
		mockRepo.On("GetByEmail", ctx, input.Email).Return(nil, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		result, err := uc.Execute(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "newemail@example.com", result.Email)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUserUseCase_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUpdateUserUseCase(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		id := uint(1)
		user := &entity.User{ID: id, Email: "test@example.com"}

		mockRepo.On("GetByID", ctx, id).Return(user, nil)

		result, err := uc.GetUser(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		id := uint(999)
		mockRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

		result, err := uc.GetUser(ctx, id)

		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNewUpdateUserUseCase(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUpdateUserUseCase(mockRepo)

	assert.NotNil(t, uc)
	assert.Equal(t, mockRepo, uc.userRepo)
}