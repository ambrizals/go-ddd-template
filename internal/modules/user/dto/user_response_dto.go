package dto

import (
	"github.com/ambrizals/go-ddd-template/internal/shared/response"
)

type UserResponseDTO = response.APIResponse[UserResponse]
type UserListResponseDTO = response.APIListResponse[UserResponse]
type DeactivateUserResponseDTO = response.APIResponse[DeactivateUserResponseData]

type DeactivateUserResponseData struct {
	ID          uint `json:"id"`
	Deactivated bool `json:"deactivated"`
}