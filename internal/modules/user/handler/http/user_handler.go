package user

import (
	"strconv"

	"github.com/ambrizals/go-ddd-template/internal/modules/user/dto"
	"github.com/ambrizals/go-ddd-template/internal/modules/user/usecase/create_user"
	"github.com/ambrizals/go-ddd-template/internal/modules/user/usecase/deactivate_user"
	"github.com/ambrizals/go-ddd-template/internal/modules/user/usecase/update_user"
	"github.com/ambrizals/go-ddd-template/internal/shared/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	createUserUseCase     *create_user.CreateUserUseCase
	updateUserUseCase     *update_user.UpdateUserUseCase
	deactivateUserUseCase *deactivate_user.DeactivateUserUseCase
	validate              *validator.Validate
}

func NewUserHandler(
	createUC *create_user.CreateUserUseCase,
	updateUC *update_user.UpdateUserUseCase,
	deactivateUC *deactivate_user.DeactivateUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase:     createUC,
		updateUserUseCase:     updateUC,
		deactivateUserUseCase: deactivateUC,
		validate:              validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration details"
// @Success 201 {object} dto.UserResponseDTO
// @Failure 400 {object} dto.UserResponseDTO
// @Router /users [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	req := new(dto.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewAPIErrorResponse[dto.UserResponse]("Bad Request", err.Error()))
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewAPIErrorResponse[dto.UserResponse]("Validation Failed", err.Error()))
	}

	createInput := create_user.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}

	output, err := h.createUserUseCase.Execute(c.Context(), createInput)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewAPIErrorResponse[dto.UserResponse]("Internal Server Error", err.Error()))
	}

	payload := &dto.UserResponse{
		ID:       output.ID,
		Email:    output.Email,
		FullName: output.FullName,
	}
	return c.Status(fiber.StatusCreated).JSON(response.NewAPIResponse(payload))
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get user details by their unique ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserResponseDTO
// @Failure 404 {object} dto.UserResponseDTO
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewAPIErrorResponse[dto.UserResponse]("Bad Request", "Invalid user ID"))
	}

	user, err := h.updateUserUseCase.GetUser(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.NewAPIErrorResponse[dto.UserResponse]("Not Found", "User not found"))
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(response.NewAPIErrorResponse[dto.UserResponse]("Not Found", "User not found"))
	}

	payload := &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.String(),
	}
	return c.JSON(response.NewAPIResponse(payload))
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update user details by their unique ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body update_user.UpdateUserInput true "User update details"
// @Success 200 {object} dto.UserResponseDTO
// @Failure 400 {object} dto.UserResponseDTO
// @Failure 404 {object} dto.UserResponseDTO
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewAPIErrorResponse[dto.UserResponse]("Bad Request", "Invalid user ID"))
	}

	req := new(update_user.UpdateUserInput)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewAPIErrorResponse[dto.UserResponse]("Bad Request", err.Error()))
	}

	req.ID = uint(id)

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewAPIErrorResponse[dto.UserResponse]("Validation Failed", err.Error()))
	}

	output, err := h.updateUserUseCase.Execute(c.Context(), *req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewAPIErrorResponse[dto.UserResponse]("Internal Server Error", err.Error()))
	}

	payload := &dto.UserResponse{
		ID:       output.ID,
		Email:    output.Email,
		FullName: output.FullName,
	}
	return c.JSON(response.NewAPIResponse(payload))
}

// DeactivateUser godoc
// @Summary Deactivate an existing user
// @Description Deactivate user by their unique ID
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} dto.DeactivateUserResponseDTO
// @Failure 400 {object} dto.DeactivateUserResponseDTO
// @Failure 404 {object} dto.DeactivateUserResponseDTO
// @Router /users/{id} [delete]
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewAPIErrorResponse[dto.DeactivateUserResponseData]("Bad Request", "Invalid user ID"))
	}

	input := deactivate_user.DeactivateUserInput{
		ID: uint(id),
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewAPIErrorResponse[dto.DeactivateUserResponseData]("Validation Failed", err.Error()))
	}

	output, err := h.deactivateUserUseCase.Execute(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewAPIErrorResponse[dto.DeactivateUserResponseData]("Internal Server Error", err.Error()))
	}

	payload := &dto.DeactivateUserResponseData{
		ID:          output.ID,
		Deactivated: output.Deactivated,
	}
	return c.JSON(response.NewAPIResponse(payload))
}