package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ambrizals/go-ddd-template/internal/modules/user/entity"
	"github.com/ambrizals/go-ddd-template/internal/shared/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	cfg := config.LoadTestConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	return db
}

func ensureTestDatabase(t *testing.T) {
	cfg := config.LoadTestConfig()

	defaultDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)
	defaultDB, _ := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{})
	if defaultDB != nil {
		defaultDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.DBName))
		defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	}
}

func setupTestDatabase(t *testing.T, db *gorm.DB) {
	err := db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, email VARCHAR(255) UNIQUE NOT NULL, password TEXT NOT NULL, full_name VARCHAR(255), created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP)").Error
	require.NoError(t, err)
}

func cleanupTestData(t *testing.T, db *gorm.DB) {
	db.Exec("DELETE FROM users")
}

func TestMain(m *testing.M) {
	ensureTestDatabase(nil)
	os.Exit(m.Run())
}

func TestUserRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	setupTestDatabase(t, db)
	cleanupTestData(t, db)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		user := &entity.User{
			Email:    "test@example.com",
			Password: "password123",
			FullName: "Test User",
		}

		err := repo.Create(ctx, user)

		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.NotZero(t, user.CreatedAt)
	})

	t.Run("duplicate email", func(t *testing.T) {
		user := &entity.User{
			Email:    "duplicate@example.com",
			Password: "password123",
			FullName: "First User",
		}

		err := repo.Create(ctx, user)
		assert.NoError(t, err)

		duplicateUser := &entity.User{
			Email:    "duplicate@example.com",
			Password: "password456",
			FullName: "Second User",
		}

		err = repo.Create(ctx, duplicateUser)

		assert.Error(t, err)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	setupTestDatabase(t, db)
	cleanupTestData(t, db)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{
		Email:    "getbyid@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		foundUser, err := repo.GetByID(ctx, user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, user.FullName, foundUser.FullName)
	})

	t.Run("not found", func(t *testing.T) {
		foundUser, err := repo.GetByID(ctx, 99999)

		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	setupTestDatabase(t, db)
	cleanupTestData(t, db)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{
		Email:    "getbyemail@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		foundUser, err := repo.GetByEmail(ctx, "getbyemail@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.FullName, foundUser.FullName)
	})

	t.Run("not found", func(t *testing.T) {
		foundUser, err := repo.GetByEmail(ctx, "nonexistent@example.com")

		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})
}

func TestUserRepository_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	setupTestDatabase(t, db)
	cleanupTestData(t, db)
	repo := NewUserRepository(db)
	ctx := context.Background()

	users := []*entity.User{
		{Email: "user1@example.com", Password: "pass1", FullName: "User One"},
		{Email: "user2@example.com", Password: "pass2", FullName: "User Two"},
		{Email: "user3@example.com", Password: "pass3", FullName: "User Three"},
	}

	for _, u := range users {
		err := repo.Create(ctx, u)
		require.NoError(t, err)
	}

	t.Run("success", func(t *testing.T) {
		foundUsers, err := repo.List(ctx)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(foundUsers), 3)
	})

	t.Run("empty list", func(t *testing.T) {
		cleanupTestData(t, db)

		foundUsers, err := repo.List(ctx)

		assert.NoError(t, err)
		assert.Empty(t, foundUsers)
	})
}