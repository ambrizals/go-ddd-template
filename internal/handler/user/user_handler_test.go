package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-ddd-template/internal/dto"
	"github.com/user/go-ddd-template/internal/entity"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockUserUseCase is a mock of entity.UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Register(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUseCase) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) ListUsers(ctx context.Context) ([]entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.User), args.Error(1)
}

func TestUserHandler_Register(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase)
	app := fiber.New()
	app.Post("/users", h.Register)

	t.Run("success", func(t *testing.T) {
		reqBody := `{"email":"test@example.com","password":"password123","full_name":"Test User"}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		mockUseCase.On("Register", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		
		var userResp dto.UserResponse
		json.NewDecoder(resp.Body).Decode(&userResp)
		assert.Equal(t, "test@example.com", userResp.Email)
		mockUseCase.AssertExpectations(t)
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
	mockUseCase := new(MockUserUseCase)
	h := NewUserHandler(mockUseCase)
	app := fiber.New()
	app.Get("/users/:id", h.GetUser)

	t.Run("success", func(t *testing.T) {
		id := uint(1)
		user := &entity.User{ID: id, Email: "test@example.com"}
		mockUseCase.On("GetUser", mock.Anything, id).Return(user, nil)

		req := httptest.NewRequest("GET", "/users/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		
		var userResp dto.UserResponse
		json.NewDecoder(resp.Body).Decode(&userResp)
		assert.Equal(t, id, userResp.ID)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		id := uint(2)
		mockUseCase.On("GetUser", mock.Anything, id).Return(nil, errors.New("not found"))

		req := httptest.NewRequest("GET", "/users/2", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}
