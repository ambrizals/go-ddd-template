package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "github.com/user/go-ddd-template/docs" // Import generated docs
	"github.com/user/go-ddd-template/internal/config"
	"github.com/user/go-ddd-template/internal/handler/user"
	"github.com/user/go-ddd-template/internal/infrastructure/database"
	"github.com/user/go-ddd-template/internal/infrastructure/logger"
	"github.com/user/go-ddd-template/internal/infrastructure/otel"
	"github.com/user/go-ddd-template/internal/infrastructure/redis"
)

// @title Go DDD API Template
// @version 1.0
// @description This is a Go backend template using DDD architecture.
// @contact.name API Support
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Initialize Logger
	logger.InitLogger()

	// 3. Initialize OpenTelemetry
	tp, err := otel.InitOTEL()
	if err != nil {
		log.Printf("Failed to init OTEL: %v", err)
	}
	defer otel.ShutdownOTEL(tp)

	// 4. Initialize Database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 5. Run Migrations
	if err := database.RunMigrations(cfg); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	// 6. Initialize Redis
	if err := redis.InitRedis(cfg); err != nil {
		log.Printf("Redis warning: %v", err)
	}

	// 7. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Go DDD API",
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(fiberLogger.New())

	// Routes
	api := app.Group("/api/v1")
	user.SetupUserRoutes(api, db)

	// Swagger UI
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
