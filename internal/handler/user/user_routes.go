package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/go-ddd-template/internal/repository"
	"github.com/user/go-ddd-template/internal/usecase"
	"gorm.io/gorm"
)

func SetupUserRoutes(router fiber.Router, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	h := NewUserHandler(userUseCase)

	users := router.Group("/users")
	users.Post("/", h.Register)
	users.Get("/:id", h.GetUser)
}
