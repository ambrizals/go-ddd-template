package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ambrizals/go-ddd-template/internal/modules/user/dto"
	"github.com/ambrizals/go-ddd-template/internal/modules/user/entity"
	"github.com/ambrizals/go-ddd-template/internal/modules/user/usecase/create_user"
	"github.com/ambrizals/go-ddd-template/internal/modules/user/usecase/deactivate_user"
	"github.com/ambrizals/go-ddd-template/internal/modules/user/usecase/update_user"
	"github.com/gofiber/fiber/v2"
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

func TestUserHandler_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	createUC := create_user.NewCreateUserUseCase(mockRepo)
	updateUC := update_user.NewUpdateUserUseCase(mockRepo)
	deactivateUC := deactivate_user.NewDeactivateUserUseCase(mockRepo)
	h := NewUserHandler(createUC, updateUC, deactivateUC)
	app := fiber.New()
	app.Post("/users", h.Register)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

		reqBody := `{"email":"test@example.com","password":"password123","full_name":"Test User"}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var apiResp struct {
			Payload dto.UserResponse `json:"payload"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		assert.Equal(t, "test@example.com", apiResp.Payload.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := `{"email":"invalid-email"}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	createUC := create_user.NewCreateUserUseCase(mockRepo)
	updateUC := update_user.NewUpdateUserUseCase(mockRepo)
	deactivateUC := deactivate_user.NewDeactivateUserUseCase(mockRepo)
	h := NewUserHandler(createUC, updateUC, deactivateUC)
	app := fiber.New()
	app.Get("/users/:id", h.GetUser)

	t.Run("success", func(t *testing.T) {
		id := uint(1)
		user := &entity.User{ID: id, Email: "test@example.com", FullName: "Test User", CreatedAt: time.Now()}
		mockRepo.On("GetByID", mock.Anything, id).Return(user, nil)

		req := httptest.NewRequest("GET", "/users/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var apiResp struct {
			Payload dto.UserResponse `json:"payload"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		assert.Equal(t, id, apiResp.Payload.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		id := uint(2)
		mockRepo.On("GetByID", mock.Anything, id).Return(nil, errors.New("user not found"))

		req := httptest.NewRequest("GET", "/users/2", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}
