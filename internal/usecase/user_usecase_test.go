package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-ddd-template/internal/entity"
	"testing"
)

// MockUserRepository is a mock of entity.UserRepository
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

func TestUserUseCase_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		user := &entity.User{Email: "test@example.com", FullName: "Test User"}
		
		mockRepo.On("GetByEmail", mock.Anything, user.Email).Return(nil, errors.New("not found"))
		mockRepo.On("Create", mock.Anything, user).Return(nil)

		err := useCase.Register(context.Background(), user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email exists", func(t *testing.T) {
		user := &entity.User{Email: "existing@example.com"}
		existingUser := &entity.User{ID: 1, Email: "existing@example.com"}

		mockRepo.On("GetByEmail", mock.Anything, user.Email).Return(existingUser, nil)

		err := useCase.Register(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmailAlreadyExists, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		id := uint(1)
		expectedUser := &entity.User{ID: id, Email: "test@example.com"}
		mockRepo.On("GetByID", mock.Anything, id).Return(expectedUser, nil)

		user, err := useCase.GetUser(context.Background(), id)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		id := uint(2)
		mockRepo.On("GetByID", mock.Anything, id).Return(nil, errors.New("not found"))

		user, err := useCase.GetUser(context.Background(), id)

		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_ListUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedUsers := []entity.User{
			{ID: 1, Email: "user1@example.com"},
			{ID: 2, Email: "user2@example.com"},
		}
		mockRepo.On("List", mock.Anything).Return(expectedUsers, nil)

		users, err := useCase.ListUsers(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		mockRepo.AssertExpectations(t)
	})
}
