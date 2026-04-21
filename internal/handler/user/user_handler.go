package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/go-ddd-template/internal/dto"
	"github.com/user/go-ddd-template/internal/entity"
	"github.com/go-playground/validator/v10"
	"strconv"
)

type UserHandler struct {
	userUseCase entity.UserUseCase
	validate    *validator.Validate
}

func NewUserHandler(useCase entity.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: useCase,
		validate:    validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration details"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /users [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	req := new(dto.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation Failed",
			Message: err.Error(),
		})
	}

	user := &entity.User{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}

	if err := h.userUseCase.Register(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.String(),
	})
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get user details by their unique ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid user ID",
		})
	}

	user, err := h.userUseCase.GetUser(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "Not Found",
			Message: "User not found",
		})
	}

	return c.JSON(dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.String(),
	})
}
