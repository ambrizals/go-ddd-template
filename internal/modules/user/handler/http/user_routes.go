package user

import (
	"github.com/ambrizals/go-ddd-template/internal/modules/user/repository"
	"github.com/ambrizals/go-ddd-template/internal/modules/user/usecase/create_user"
	"github.com/ambrizals/go-ddd-template/internal/modules/user/usecase/deactivate_user"
	"github.com/ambrizals/go-ddd-template/internal/modules/user/usecase/update_user"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupUserRoutes sets up all user-related HTTP routes
func SetupUserRoutes(router fiber.Router, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)

	// Initialize workflow-specific use cases
	createUserUC := create_user.NewCreateUserUseCase(userRepo)
	updateUserUC := update_user.NewUpdateUserUseCase(userRepo)
	deactivateUserUC := deactivate_user.NewDeactivateUserUseCase(userRepo)

	// Initialize handler with all use cases
	h := NewUserHandler(createUserUC, updateUserUC, deactivateUserUC)

	users := router.Group("/users")
	users.Post("/", h.Register)
	users.Get("/:id", h.GetUser)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeactivateUser)
}
